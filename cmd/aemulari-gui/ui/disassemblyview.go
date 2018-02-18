package ui

import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"

	ae "../../../aemulari"
)

type DisassemblyInfo struct {
	curr DisassemblyList
	prev DisassemblyList
}

type DisassemblyList struct {
	addr    uint64           // Address of first instruction in disassembly
	entries []ae.Disassembly // Disassembly list
}

func (ui *Ui) refreshDisasmView() error {
	if v, err := ui.g.View(vDisasm); err != nil {
		return err
	} else {
		return ui.updateDisasmView(v)
	}
}

func (d DisassemblyList) Contains(addr uint64) (bool, ae.Disassembly) {
	for _, e := range d.entries {
		if e.AddressU64 == addr {
			return true, e
		}
	}
	return false, ae.Disassembly{}
}

func (ui Ui) getBpSymbolAt(addr uint64) string {
	bps := ui.dbg.GetBreakpointsAt(addr)
	if len(bps) == 0 {
		return " "
	}

	if bps.Enabled() {
		return ui.theme.ArmedBreakpointSymbol()
	}

	return ui.theme.DisabledBreakpointSymbol()
}

func (ui Ui) getPcSymbolAt(addr uint64) string {
	if addr == ui.pc {
		return ui.theme.CurrentInstructionSymbol()
	}
	return " "
}

func (ui Ui) getLineAnnotations(addr uint64) string {
	annotations := ui.getBpSymbolAt(addr)
	annotations += ui.getPcSymbolAt(addr)
	return annotations
}

func (ui *Ui) updateDisasmView(view *gocui.View) error {
	var err error
	var line string

	_, height := view.Size()

	if height < 1 {
		return errors.New("Disassembly view not large enough to draw")
	}

	hasPc, _ := ui.disasm.curr.Contains(ui.pc)
	if !hasPc {
		ui.disasm.curr.addr = ui.pc
	}

	// Always re-read in case of self-modifying code
	ui.disasm.curr.entries, err = ui.dbg.DisassembleAt(ui.disasm.curr.addr, uint64(height))
	if err != nil {
		return err
	}

	prevLen := len(ui.disasm.prev.entries)
	currLen := len(ui.disasm.curr.entries)

	if ui.disasm.prev.addr != ui.disasm.curr.addr || prevLen != currLen {
		ui.disasm.prev.entries = make([]ae.Disassembly, currLen)
		ui.disasm.prev.addr = ui.disasm.curr.addr
		copy(ui.disasm.prev.entries, ui.disasm.curr.entries)
	}

	view.Clear()

	for i, e := range ui.disasm.curr.entries {
		annotation := ui.getLineAnnotations(e.AddressU64)

		if e.Equals(ui.disasm.prev.entries[i]) {
			// FIXME Get this from ui.dbg
			addrFmt := "%08x"
			line = fmt.Sprintf("%s <%s>  %s %s\n",
				ui.theme.ColorAddress(addrFmt, e.AddressU64),
				ui.theme.ColorOpcode(e.Opcode),
				ui.theme.ColorMnemonic(e.Mnemonic),
				ui.theme.ColorOperands(e.Operands))
		} else {
			line = fmt.Sprintf("%s <%s>  %s %s\n", e.Address, e.Opcode, e.Mnemonic, e.Operands)
			line = ui.theme.ColorModifiedInstruction(line)
		}

		fmt.Fprintf(view, annotation+line)
	}

	return nil
}
