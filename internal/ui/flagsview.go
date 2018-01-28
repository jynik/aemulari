package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type FlagInfo struct {
	curr    []string
	prev    []string
	tainted bool
}

func (ui *Ui) updateFlagsView(view *gocui.View) error {

	newFlags := []string{}
	highlightChanges := ui.theme.ColorIfStringsDiffer

	regs, err := ui.dbg.ReadRegAll()
	if err != nil {
		return err
	}

	for _, reg := range regs {
		for _, flag := range reg.Reg.Flags {
			nv := flag.GetNamedString(reg.Value)
			newFlags = append(newFlags, nv)
		}
	}

	if len(ui.flags.prev) == 0 {
		ui.flags.prev = newFlags
	}

	view.Clear()
	numFlags := len(newFlags)
	for i := 0; i < numFlags; i += 2 {
		line := " " + highlightChanges(newFlags[i], ui.flags.prev[i])

		if (i + 1) < numFlags {
			for j := len(newFlags[i]); j < 16; j++ {
				line += " "
			}

			line += highlightChanges(newFlags[i+1], ui.flags.prev[i+1])
		}

		fmt.Fprintln(view, line)
	}

	if ui.flags.tainted {
		ui.flags.prev = ui.flags.curr
		ui.flags.curr = newFlags
		ui.flags.tainted = false
	}

	return nil
}
