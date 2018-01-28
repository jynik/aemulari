package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"

	"../arch"
)

type RegInfo struct {
	curr    []arch.RegisterValue
	prev    []arch.RegisterValue
	tainted bool // Track if register values may have chaned
}

func (ui *Ui) updateRegView(view *gocui.View) error {
	var err error

	if len(ui.regs.curr) != 0 && ui.regs.tainted {
		if len(ui.regs.curr) != len(ui.regs.prev) {
			ui.regs.prev = make([]arch.RegisterValue, len(ui.regs.curr))
		}

		copy(ui.regs.prev, ui.regs.curr)
		ui.regs.tainted = false
	}

	if ui.regs.curr, err = ui.dbg.ReadRegAll(); err != nil {
		return err
	}

	if len(ui.regs.prev) == 0 {
		ui.regs.prev = make([]arch.RegisterValue, len(ui.regs.curr))
		copy(ui.regs.prev, ui.regs.curr)
	}

	view.Clear()
	for i := 0; i < len(ui.regs.curr); i += 2 {
		if i+1 < len(ui.regs.curr) {
			fmt.Fprintf(view, " %s    %s\n",
				ui.theme.ColorIfStringsDiffer(ui.regs.curr[i].String(),
					ui.regs.prev[i].String()),
				ui.theme.ColorIfStringsDiffer(ui.regs.curr[i+1].String(),
					ui.regs.prev[i+1].String()))

		} else {
			fmt.Fprintf(view, " %s\n",
				ui.theme.ColorIfStringsDiffer(ui.regs.curr[i].String(),
					ui.regs.prev[i].String()))
		}
	}

	return nil
}