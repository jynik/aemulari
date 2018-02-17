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

// A Debugger, after being created via NewDebugger(), may be used to
// execute (via emulation) and inspect a program.
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

// A single disassembled instruction separated into its components
type Disassembly struct {
	AddressU64 uint64		// Address of the instruction, as a uint64
	Address    string		// Address of the instruction, as a string
	Opcode     string		// String representation of the binary opcode
	Mnemonic   string		// String representation of the instruction mnemonic
	Operands   string		// String representation of the instruction operands
}

// Returns true if two instructions are the same, and false otherwise.
// Note that this implies exactly the same, not just semantically the same.
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

// Create a Debugger using the provided Architecture and DebuggerConfig.
func NewDebugger(a Architecture, c DebuggerConfig) (*Debugger, error) {
	var d Debugger

	if err := d.init(a, c, false); err != nil {
		return nil, err
	}

	return &d, nil
}

// Reset the debugger to its original state. If `keepMappings` is true,
// the current state of memory mappyings will be retained. This is often
// useful if mappings have been added over the course of debugging.
//
// Otherwise, if `keepMappings` is false, existing memory mappings will be
// unmapped, writing memory contents to any configured output files. Then,
// mappings will be re-initialized based upon the configuration specified in
// the DebugConfiguration provided when the Debugger was created.
func (d *Debugger) Reset(keepMappings bool) error {
	var err, firstError error
	d.mu.Close()
	d.cs.Close()

	// Reset to original configuration, but keep memory mappings
	// This is intended to allow us to map regions as we discover they're
	// necessary, and not have to re-do that with each reset.
	newConfig := d.cfg
	if keepMappings {
		newConfig.Mem = d.mapped
	} else {
		for _, r := range d.mapped {
			// Don't treat the error as fatal; allow the reset to complete
			err := d.Unmap(r.name)
			if err != nil && firstError == nil {
				firstError = err
			}
		}
	}

	err = d.init(d.arch, newConfig, true)
	if err != nil {
		if firstError == nil {
			return err
		}
	}

	return firstError
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

// Deinitialize a Debugger and unmap memory regions, writing their contents
// to output files if configured to do so.
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

// Map a memory region described by `toMap`. If the MemRegion's `inputFile`
// field is non-empty, the contents of the associated file will be used to
// initialize the region. If the MemRegion's `outputFile` field is non-empty,
// the contents of memory will be written to this file when unmapped by a call
// to Debugger.Unmap(), or Debugger.Reset(false).
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

// Unmapped the memory region named `name`. If the `outputFile` field specified when
// the region was mapped was non-empty, the contents of the memory will be written
// to this file.
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

// Retrieve the emulated processor's current Endianness.
func (d *Debugger) Endianness() (Endianness, error) {
	regs, err := d.ReadRegAll()
	if err != nil {
		return LittleEndian, err
	}

	return d.arch.endianness(regs), nil
}

// Retrieve the current program counter value.
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

// Retrieve the current state all registers.
func (d *Debugger) ReadRegAll() ([]RegisterValue, error) {
	var err error

	regDefs := d.arch.registers()
	regVals := make([]RegisterValue, len(regDefs), len(regDefs))

	for i, reg := range regDefs {
		regVals[i], err = d.readReg(reg)
		if err != nil {
			return []RegisterValue{}, err
		}
	}

	return regVals, nil
}

// Retrieve the current state of the register described by `reg`.
func (d *Debugger) readReg(reg *RegisterDef) (RegisterValue, error) {
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

// Read a register and update the its value in the provided RegisterValue
func (d *Debugger) ReadReg(rv *RegisterValue) error {
	reg, err := d.readReg(rv.Reg)
	if err != nil {
		return err
	}
	rv.Value = reg.Value
	return nil
}

// Retrieve the current state of a register, specified by its name
func (d *Debugger) ReadRegByName(name string) (RegisterValue, error) {
	var rv RegisterValue
	if reg, err := d.arch.register(name); err == nil {
		return d.readReg(reg)
	} else {
		return rv, err
	}
}

// Update the value of a single register.
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

// Update the values of a set of registers.
func (d *Debugger) WriteRegs(rvs []RegisterValue) error {
	for _, rv := range rvs {
		if err := d.WriteReg(rv); err != nil {
			return err
		}
	}

	return nil
}

// Update the value of a single register, specified by name.
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

// Read `size` bytes of memory starting at `addr`.
func (d *Debugger) ReadMem(addr, size uint64) ([]byte, error) {
	return d.mu.MemRead(addr, size)
}

// Write `data` to memory at the address specified by `addr`
func (d *Debugger) WriteMem(addr uint64, data []byte) error {
	return d.mu.MemWrite(addr, data)
}

// Execute `count` instructions and then return.
// A negative count implies "Run until a breakpoint or exception"
// Returns (hitException, intNumber, err)
func (d *Debugger) Step(count int64) (Exception, error) {
	if count <= 0 {
		return Exception{}, errors.New("Debugger.Step() requires that count >= 1.")
	}

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

// Start or continue execution in the debugger. Upon hitting a breakpoint
// or encountering an execption, this function will return. The error
// value should be tested first to determine if the debugger encountered
// an unexpected error. Next, the returned Exeception's Occurred() method
// should be called to determine if an exception occured. If so, the
// Exception.String() method may be used to retrieve information about the
// exception.
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

// Code step callback
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

// Interrupt callback
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

// Disassemble `count` instructions, starting at the current program counter
func (d *Debugger) Disassemble(count uint64) ([]Disassembly, error) {
	if rv, err := d.ReadRegByName("pc"); err != nil {
		return []Disassembly{}, nil
	} else {
		return d.DisassembleAt(rv.Value, count)
	}
}

// Disassemble `count` instructions, starting at the address specified by `addr`.
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

// Set a breakpoint at the specified address. It will automatically
// be assigned an ID.
func (d *Debugger) SetBreakpoint(addr uint64) Breakpoint {
	return d.bps.add(addr)
}

// Delete all existing breakpoints.
func (d *Debugger) DeleteAllBreakpoints() {
	d.bps.removeAll()
}

// Delete all breakpoints at the specified address/
func (d *Debugger) DeleteBreakpointsAt(addr uint64) {
	d.bps.removeAllAt(addr)
}

// Delete the breakpoint associatedw with the specified ID.
func (d *Debugger) DeleteBreakpoint(id int) {
	d.bps.remove(id)
}

// Get a list of all breakpoints.
func (d *Debugger) GetBreakpoints() BreakpointList {
	return d.bps.get()
}

// Get a list of breakpoints set at the specified address.
func (d *Debugger) GetBreakpointsAt(addr uint64) BreakpointList {
	return d.bps.getAllAt(addr)
}
