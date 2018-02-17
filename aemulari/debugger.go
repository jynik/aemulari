package aemulari

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"

	cs "github.com/bnagy/gapstone"
	"github.com/op/go-logging"
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

// Top-level debugger object
type Debugger struct {
	arch   Architecture
	mu     uc.Unicorn     // Unicorn emulator handle
	cs     cs.Engine      // Capstone disassembly engine handle
	cfg    DebuggerConfig // Configuration settings
	mapped MemRegions     // Mapped memory regions
	step   codeStep       // Code stepping metadata
	bps    breakpoints    // Breakpoint settings
	exInfo exceptionInfo  // CPU Exception handling
}

// Configuration of Debugger's initial state
type DebuggerConfig struct {
	RegDefs []RegisterValue // Default register values
	Mem     MemRegions      // Memory region configuration
}

// A single disassembled instruction
type Disassembly struct {
	AddressU64 uint64		// Address of the instruction, as a uint64
	Address    string		// Address of the instruction, as a string
	Opcode     string		// String representation of the binary opcode
	Mnemonic   string		// String representation of the instruction mnemonic
	Operands   string		// String representation of the instruction operands
}

// Returns true if two instructions are equals, and false otherwise
func (d Disassembly) Equals(other Disassembly) bool {
	return d.AddressU64 == other.AddressU64 &&
		d.Address == other.Address &&
		d.Opcode == other.Opcode &&
		d.Mnemonic == other.Mnemonic &&
		d.Operands == other.Operands
}

// Exception callback handling
type exceptionInfo struct {
	dbg  *Debugger
	hook uc.Hook
	last Exception // Most recently occurring exception
}

// Data used to implement stepping and breakpoints
type codeStep struct {
	dbg     *Debugger
	count   int64
	hook    uc.Hook
	options uc.UcOptions

	// Need to backup state prior to stopping emulator and restore it
	// after we return from our execution. Unclear if this is necessitated
	// due to a Unicorn defect, or our own misuse of the framework
	regs []RegisterValue
}

var log = logging.MustGetLogger("")

// Instantiate and configure a new Debugger
func NewDebugger(a Architecture, c DebuggerConfig) (*Debugger, error) {
	var d Debugger

	if err := d.init(a, c, false); err != nil {
		return nil, err
	}

	return &d, nil
}

func (d *Debugger) Reset() error {
	d.mu.Close()
	d.cs.Close()

	/* Reset to original configuration, but keep memory mappings
	 * This is intended to allow us to map regions as we discover they're
	 * necessary, and not have to re-do that with each reset.
	 *
	 * TODO: Make this an opt-out configuration item?
	 */
	newConfig := d.cfg
	newConfig.Mem = d.mapped

	return d.init(d.arch, newConfig, true)
}

func (d *Debugger) init(arch Architecture, cfg DebuggerConfig, reset bool) error {
	var err error

	d.cfg = cfg
	d.arch = arch

	// Keep existing breakpoints if we're resetting the debugger
	if !reset {
		d.bps.initialize()
	}

	d.mu, err = uc.NewUnicorn(d.arch.id().uc, d.arch.initialMode().uc)
	if err != nil {
		return err
	}

	d.cs, err = cs.New(d.arch.id().cs, d.arch.initialMode().cs)
	if err != nil {
		d.mu.Close()
		return err
	}

	// TODO customize invalid instruction handling
	d.cs.SkipDataStart(nil)

	// Load memory regions
	d.mapped = make(MemRegions)
	for _, m := range d.cfg.Mem {
		if err = d.Map(m); err != nil {
			log.Debugf("Failed to map %s", m)
			return d.closeAll(err)
		}

		log.Debugf("Mapped %s", m)
	}

	// Load default register values
	var loadedPc bool = false
	for _, r := range d.cfg.RegDefs {
		log.Debugf("Loading %s", r)
		if r.Reg.IsProgramCounter() {
			loadedPc = true
			r.Value = d.arch.initialPC(r.Value)
		}

		if err := d.mu.RegWrite(r.Reg.uc, r.Value); err != nil {
			return d.closeAll(err)
		}
	}

	// If the register used as the program counter was not specified,
	// default it to the start of code memory.
	if !loadedPc {
		codeMem := d.code()

		pc, err := d.arch.register("pc")
		if err != nil {
			return d.closeAll(err)
		}
		pcVal := d.arch.initialPC(codeMem.base)

		if err := d.mu.RegWrite(pc.uc, pcVal); err != nil {
			return d.closeAll(err)
		}
	}

	// Code stepping setup
	d.step.options = uc.UcOptions{Timeout: 0, Count: 0}
	d.step.dbg = d
	codeMem := d.code()
	d.step.hook, err = d.mu.HookAdd(uc.HOOK_CODE, d.step.cb, codeMem.base, codeMem.size)
	if err != nil {
		return d.closeAll(err)
	}

	d.exInfo.dbg = d
	d.exInfo.hook, err = d.mu.HookAdd(uc.HOOK_INTR, d.exInfo.cb, 0, 0)
	if err != nil {
		return d.closeAll(err)
	}

	return nil
}

// Deinitialize Debugger and write any requested output files
func (d *Debugger) Close() error {
	var ret error = nil

	for name := range d.cfg.Mem {
		if err := d.Unmap(name); ret == nil && err != nil {
			ret = err
		}
	}

	return d.closeAll(ret)
}

func (d *Debugger) closeAll(e error) error {
	d.mu.Close()
	d.cs.Close()
	return e
}

// Code region access that must succeed
func (d *Debugger) code() MemRegion {
	if codeMem, err := d.mapped.Get("code"); err != nil {
		panic(err)
	} else {
		return codeMem
	}
}

func (d *Debugger) Map(toMap MemRegion) error {
	var prot int

	if toMap.size == 0 {
		return errors.New("Zero-length mappings are not permitted.")
	}

	if d.mapped.Contains(toMap.name) {
		return fmt.Errorf("A mapping named \"%s\" already exists.", toMap.name)
	}

	prot = 0
	if toMap.perms.Read {
		prot |= uc.PROT_READ
	}
	if toMap.perms.Write {
		prot |= uc.PROT_WRITE
	}
	if toMap.perms.Exec {
		prot |= uc.PROT_EXEC
	}

	if err := d.mu.MemMapProt(toMap.base, toMap.size, prot); err != nil {
		return err
	}

	if data, err := toMap.LoadInputData(); err == nil {
		if err := d.WriteMem(toMap.base, data); err != nil {
			return err
		}
	} else {
		return err
	}

	d.mapped.Add(toMap)
	return nil
}

func (d *Debugger) Unmap(name string) error {
	var m MemRegion
	var err error
	var ret error = nil

	if m, err = d.mapped.Get(name); err != nil {
		return err
	}

	// An output file name indicates we want to save the contents of this region
	if m.outputFile != "" {
		data, err := d.mu.MemRead(m.base, m.size)
		if err == nil {
			err = ioutil.WriteFile(m.outputFile, data, 0644)
		}

		ret = err
	}

	err = d.mu.MemUnmap(m.base, m.size)
	if ret != nil {
		ret = err
	}
	d.mapped.Remove(m.name)
	return ret
}

func (d *Debugger) Endianness() (Endianness, error) {
	regs, err := d.ReadRegAll()
	if err != nil {
		return LittleEndian, err
	}

	return d.arch.endianness(regs), nil
}

func (d *Debugger) pc() (uint64, error) {
	regs, err := d.ReadRegAll()
	if err != nil {
		return 0xdeadbeefdeadbeef, err
	}

	for _, reg := range regs {
		if reg.Reg.IsProgramCounter() {
			pc := d.arch.currentPC(reg.Value, regs)
			return pc, nil
		}
	}

	panic("Failed to locate program counter")
}

func (d *Debugger) ReadRegAll() ([]RegisterValue, error) {
	var err error

	regDefs := d.arch.registers()
	regVals := make([]RegisterValue, len(regDefs), len(regDefs))

	for i, reg := range regDefs {
		regVals[i], err = d.ReadReg(reg)
		if err != nil {
			return []RegisterValue{}, err
		}
	}

	return regVals, nil
}

func (d *Debugger) ReadReg(reg *RegisterDef) (RegisterValue, error) {
	var rv RegisterValue
	var val uint64
	var err error

	if val, err = d.mu.RegRead(reg.uc); err != nil {
		return rv, err
	}

	rv.Reg = reg
	rv.Value = val
	return rv, nil
}

func (d *Debugger) ReadRegByName(name string) (RegisterValue, error) {
	var rv RegisterValue
	if reg, err := d.arch.register(name); err == nil {
		return d.ReadReg(reg)
	} else {
		return rv, err
	}
}

func (d *Debugger) WriteReg(rv RegisterValue) error {
	if rv.Reg.IsProgramCounter() {
		if regs, err := d.ReadRegAll(); err != nil {
			return err
		} else {
			rv.Value = d.arch.currentPC(rv.Value, regs)
		}
	}
	return d.mu.RegWrite(rv.Reg.uc, rv.Value)
}

func (d *Debugger) WriteRegs(rvs []RegisterValue) error {
	for _, rv := range rvs {
		if err := d.WriteReg(rv); err != nil {
			return err
		}
	}

	return nil
}

func (d *Debugger) WriteRegByName(name string, value uint64) error {
	var rv RegisterValue
	if reg, err := d.arch.register(name); err == nil {
		rv.Reg = reg
		rv.Value = value
		return d.WriteReg(rv)
	} else {
		return err
	}
}

func (d *Debugger) ReadMem(addr, size uint64) ([]byte, error) {
	return d.mu.MemRead(addr, size)
}

func (d *Debugger) WriteMem(addr uint64, data []byte) error {
	return d.mu.MemWrite(addr, data)
}

// A negative count implies "Run until a breakpoint or exception"
// Returns (hitException, intNumber, err)
func (d *Debugger) Step(count int64) (Exception, error) {
	d.step.regs = []RegisterValue{}
	d.step.count = count
	d.exInfo.last = Exception{}

	log.Debugf("Stepping %d instructions.", d.step.count)

	pc, err := d.pc()
	if err != nil {
		return d.exInfo.last, err
	}

	err = d.mu.StartWithOptions(pc, d.code().End(), &d.step.options)
	if err != nil {
		return d.exInfo.last, err
	}

	return d.exInfo.last, d.WriteRegs(d.step.regs)
}

func (d *Debugger) Continue() (Exception, error) {
	d.step.regs = []RegisterValue{}
	d.step.count = -1
	d.exInfo.last = Exception{}

	pc, err := d.pc()
	if err != nil {
		return d.exInfo.last, err
	}

	err = d.mu.StartWithOptions(pc, d.code().End(), &d.step.options)
	if err != nil {
		return d.exInfo.last, err
	}

	return d.exInfo.last, d.WriteRegs(d.step.regs)
}

func (h *codeStep) cb(mu uc.Unicorn, addr uint64, size uint32) {
	d := h.dbg

	log.Debugf("Code step hook @ 0x%08x (%d), countdown=%d", addr, size, d.step.count)

	breakpointTriggered := d.bps.process(addr)
	if breakpointTriggered {
		log.Debugf("Breakpoint triggered @ 0x%08x", addr)
	}

	if breakpointTriggered || d.step.count == 0 {
		var err error

		// The state of PC and status registers (e.g., ARM CPSR) will change
		// after calling mu.Stop(). Back them up and restore them for the next
		// time we start.
		d.step.regs, err = d.ReadRegAll()
		if err != nil {
			log.Errorf("Failed to backup registers: %s", err)
		}

		if err = mu.Stop(); err != nil {
			log.Errorf("Failed to halt execution: %s", err)
		} else {
			log.Debugf("Stopping execution @ 0x%08x", addr)
		}
	} else if d.step.count > 0 {
		d.step.count -= 1
	}
}

func (e *exceptionInfo) cb(mu uc.Unicorn, intno uint32) {
	var instr []byte

	d := e.dbg
	instrLen := d.arch.maxInstructionSize()

	d.mu.Stop()

	pc, err := d.ReadRegByName("pc")
	if err != nil {
		panic("Failed to read pc in interrupt callback.")
	}

	regs, err := d.ReadRegAll()
	if err != nil {
		panic("Failed to read registers in interrupt callback.")
	}

	// FIXME this could fail in ARM Thumb mode if it's the last instruction
	// TODO  Implement smarter approach (and check for valid disassembly?)
	instr, err = d.ReadMem(pc.Value, uint64(instrLen))
	if err != nil {
		panic("Failed to read current instruction in interrupt callback.")
	}

	d.exInfo.last = d.arch.exception(intno, regs, instr)
}

// Disassemble `count` instructions, starting at PC
func (d *Debugger) Disassemble(count uint64) ([]Disassembly, error) {
	if rv, err := d.ReadRegByName("pc"); err != nil {
		return []Disassembly{}, nil
	} else {
		return d.DisassembleAt(rv.Value, count)
	}
}

func (d *Debugger) DisassembleAt(addr uint64, count uint64) ([]Disassembly, error) {
	var ret []Disassembly

	len := count * uint64(d.arch.maxInstructionSize())

	if code, err := d.ReadMem(addr, len); err != nil {
		return ret, nil
	} else {
		if instrs, err := d.cs.Disasm(code, addr, count); err == nil {
			for _, instr := range instrs {
				var entry Disassembly
				entry.AddressU64 = uint64(instr.Address)
				entry.Address = fmt.Sprintf("%08x", instr.Address)
				entry.Opcode = hex.EncodeToString(instr.Bytes)
				entry.Mnemonic = instr.Mnemonic
				entry.Operands = instr.OpStr
				ret = append(ret, entry)
			}
		}
	}

	return ret, nil
}

func (d *Debugger) SetBreakpoint(addr uint64) Breakpoint {
	return d.bps.add(addr)
}

func (d *Debugger) DeleteAllBreakpoints() {
	d.bps.removeAll()
}

func (d *Debugger) DeleteBreakpointsAt(addr uint64) {
	d.bps.removeAllAt(addr)
}

func (d *Debugger) DeleteBreakpoint(id int) {
	d.bps.remove(id)
}

func (d *Debugger) GetBreakpoints() BreakpointList {
	return d.bps.get()
}

func (d *Debugger) GetBreakpointsAt(addr uint64) BreakpointList {
	return d.bps.getAllAt(addr)
}
