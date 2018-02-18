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

	flags := []string{}
	highlightChanges := ui.theme.ColorIfStringsDiffer

	regs, err := ui.dbg.ReadRegAll()
	if err != nil {
		return err
	}

	for _, reg := range regs {
		flags = append(flags, reg.FlagStrings()...)
	}

	if len(ui.flags.prev) == 0 {
		ui.flags.prev = flags
	}

	view.Clear()
	numFlags := len(flags)
	for i := 0; i < numFlags; i += 2 {
		line := " " + highlightChanges(flags[i], ui.flags.prev[i])

		if (i + 1) < numFlags {
			for j := len(flags[i]); j < 16; j++ {
				line += " "
			}

			line += highlightChanges(flags[i+1], ui.flags.prev[i+1])
		}

		fmt.Fprintln(view, line)
	}

	if ui.flags.tainted {
		ui.flags.prev = ui.flags.curr
		ui.flags.curr = flags
		ui.flags.tainted = false
	}

	return nil
}
