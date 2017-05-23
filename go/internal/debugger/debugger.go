package debugger

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"

	cs "github.com/bnagy/gapstone"
	"github.com/op/go-logging"
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"

	"../arch"
)

// Top-level debugger object
type Debugger struct {
	mu     uc.Unicorn
	cs     cs.Engine
	cfg    Config
	mapped MemRegions
	step   codeStep
	bps    Breakpoints
	exInfo exceptionInfo // CPU Exception handling
}

// Debugger configuration
type Config struct {
	Arch    arch.Arch            // Architecture definition to use
	RegDefs []arch.RegisterValue // Default register values
	Mem     MemRegions           // Memory region configuration
}

// Disassembly
type Disassembly struct {
	AddressU64 uint64
	Address    string
	Opcode     string
	Mnemonic   string
	Operands   string
}

func (d Disassembly) Equals(other Disassembly) bool {
	return d.AddressU64 == other.AddressU64 &&
		d.Address == other.Address &&
		d.Opcode == other.Opcode &&
		d.Mnemonic == other.Mnemonic &&
		d.Operands == other.Operands
}

// Exception handling
type exceptionInfo struct {
	dbg  *Debugger
	hook uc.Hook
	last arch.Exception
}

// Data used to implement single-stepping execution
type codeStep struct {
	dbg     *Debugger
	pc      uint64
	count   int64
	hook    uc.Hook
	options uc.UcOptions
}

var log = logging.MustGetLogger("")

// Instantiate and configure a new Debugger
func New(c Config) (*Debugger, error) {
	var d Debugger
	d.cfg = c

	if err := d.init(d.cfg.Mem, false); err != nil {
		return nil, err
	}

	return &d, nil
}

func (d *Debugger) Reset() error {
	d.mu.Close()
	d.cs.Close()
	return d.init(d.mapped, true)
}

func (d *Debugger) init(toMap MemRegions, reset bool) error {
	var err error

	// Keep existing breakpoints if we're resetting the debugger
	if !reset {
		d.bps.Initialize()
	}

	archType := d.cfg.Arch.Type()
	archMode := d.cfg.Arch.Defaults().Mode

	d.mu, err = uc.NewUnicorn(archType.Uc, archMode.Uc)
	if err != nil {
		return err
	}

	d.cs, err = cs.New(archType.Cs, archMode.Cs)
	if err != nil {
		d.mu.Close()
		return err
	}

	// TODO customize invalid instruction handling
	d.cs.SkipDataStart(nil)

	// Load memory regions
	d.mapped = make(MemRegions)
	for _, m := range toMap {
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
		if err := d.mu.RegWrite(r.Reg.Uc(), r.Value); err != nil {
			return d.closeAll(err)
		} else if r.Reg.IsProgramCounter() {
			loadedPc = true
			d.step.pc = r.Value
		}
	}

	// If the register used as the program counter was not specified,
	// default it to the start of code memory.
	if !loadedPc {
		codeMem := d.code()

		pc, err := d.cfg.Arch.Register("pc")
		if err != nil {
			return d.closeAll(err)
		}

		if err := d.mu.RegWrite(pc.Uc(), codeMem.base); err != nil {
			return d.closeAll(err)
		}

		d.step.pc = codeMem.base
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

func (d *Debugger) Endianness() (arch.Endianness, error) {
	regs, err := d.ReadRegAll()
	if err != nil {
		return arch.LittleEndian, err
	}

	return d.cfg.Arch.Endianness(regs), nil
}

func (d *Debugger) ReadRegAll() ([]arch.RegisterValue, error) {
	var err error

	regDefs := d.cfg.Arch.Registers()
	regVals := make([]arch.RegisterValue, len(regDefs), len(regDefs))

	for i, reg := range regDefs {
		regVals[i], err = d.ReadReg(reg)
		if err != nil {
			return []arch.RegisterValue{}, err
		}
	}

	return regVals, nil
}

func (d *Debugger) ReadReg(reg *arch.RegisterDef) (arch.RegisterValue, error) {
	var rv arch.RegisterValue
	var val uint64
	var err error

	if val, err = d.mu.RegRead(reg.Uc()); err != nil {
		return rv, err
	}

	rv.Reg = reg
	rv.Value = val
	return rv, nil
}

func (d *Debugger) ReadRegByName(name string) (arch.RegisterValue, error) {
	var rv arch.RegisterValue
	if reg, err := d.cfg.Arch.Register(name); err == nil {
		return d.ReadReg(reg)
	} else {
		return rv, err
	}
}

func (d *Debugger) WriteReg(rv arch.RegisterValue) error {
	err := d.mu.RegWrite(rv.Reg.Uc(), rv.Value)
	if err == nil && rv.Reg.IsProgramCounter() {
		d.step.pc = rv.Value
	}
	return err
}

func (d *Debugger) WriteRegByName(name string, value uint64) error {
	var rv arch.RegisterValue
	if reg, err := d.cfg.Arch.Register(name); err == nil {
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
func (d *Debugger) Step(count int64) (arch.Exception, error) {
	d.step.count = count
	d.exInfo.last = arch.Exception{}

	log.Debugf("Stepping %d instructions.", d.step.count)

	err := d.mu.StartWithOptions(d.step.pc, d.code().End(), &d.step.options)
	err2 := d.WriteRegByName("pc", d.step.pc)

	log.Debugf("PC @ 0x%08x after Step()", d.step.pc)

	if err != nil {
		return d.exInfo.last, err
	} else {
		return d.exInfo.last, err2
	}
}

func (d *Debugger) Continue() (arch.Exception, error) {
	d.step.count = -1
	d.exInfo.last = arch.Exception{}

	err := d.mu.StartWithOptions(d.step.pc, d.code().End(), &d.step.options)
	if err != nil {
		return d.exInfo.last, err
	}

	return d.exInfo.last, d.WriteRegByName("pc", d.step.pc)
}

func (h *codeStep) cb(mu uc.Unicorn, addr uint64, size uint32) {
	d := h.dbg
	d.step.pc = addr

	log.Debugf("Code step hook @ 0x%08x (%d), countdown=%d", addr, size, d.step.count)

	breakpointTriggered := d.bps.Process(addr)
	if breakpointTriggered {
		log.Debugf("Breakpoint triggered @ 0x%08x", addr)
	}

	if breakpointTriggered || d.step.count == 0 {
		if err := mu.Stop(); err != nil {
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
	instrLen := d.cfg.Arch.MaxInstrLen()

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

	d.exInfo.last = d.cfg.Arch.Exception(intno, regs, instr)
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

	// FIXME integer overflow
	len := count * uint64(d.cfg.Arch.MaxInstrLen())

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
	return d.bps.Add(addr)
}

func (d *Debugger) DeleteAllBreakpoints() {
	d.bps.RemoveAll()
}

func (d *Debugger) DeleteBreakpointsAt(addr uint64) {
	d.bps.RemoveAllAt(addr)
}

func (d *Debugger) DeleteBreakpoint(id int) {
	d.bps.Remove(id)
}

func (d *Debugger) GetBreakpoints() BreakpointList {
	return d.bps.Get()
}

func (d *Debugger) GetBreakpointsAt(addr uint64) BreakpointList {
	return d.bps.GetAt(addr)
}

// Returns a *Regexp for matching register names and aliases
func (d *Debugger) RegisterRegexp() *regexp.Regexp {
	return d.cfg.Arch.RegisterRegexp()
}
