package ui

import (
	"strings"

	"github.com/jroimartin/gocui"
)

func (ui *Ui) handleInput(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {

	curX, curY := v.Cursor()
	atLeftBound := curX <= 2

	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case !atLeftBound && (key == gocui.KeyBackspace || key == gocui.KeyBackspace2):
		v.EditDelete(true)
	case !atLeftBound && key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyArrowUp:
		ui.redrawCommandsView(v, ui.hist.PrevEntryLine())
	case key == gocui.KeyArrowDown:
		ui.redrawCommandsView(v, ui.hist.NextEntryLine())
	case key == gocui.KeyCtrlU:
		ui.redrawCommandsView(v, "")
	case key == gocui.KeyEnter:
		if line, err := v.Line(curY); err != nil {
			return
		} else {
			promptLen := len(ui.theme.CmdPrompt())

			// Strip prompt and excess space
			if len(line) >= promptLen {
				line = line[promptLen:]
			}

			// Trim whitespace and null terminator
			line = strings.Trim(line, " \r\n\x00")

			// Empty line - execute the last item in the history
			if len(line) == 0 {
				histLen := len(ui.hist.entries)
				if histLen != 0 {
					line = ui.hist.entries[histLen-1].line
				} else {
					// Nothing in history to run
					return
				}
			}

			output, suppressHistory, err := ui.handleCommand(line)

			if output != "" {
				ui.appendConsole("\n" + output)
			}

			if err != nil {
				ui.appendConsole("\n" + ui.theme.ErrorMessage(err))
			}

			if !suppressHistory {
				histStatusSymbol := ui.theme.CmdSuccessSymbol()
				if err != nil {
					histStatusSymbol = ui.theme.CmdFailureSymbol()
				}

				ui.hist.Append(histStatusSymbol, line)
			}
			ui.redrawCommandsView(v, "")
		}
	}
}

func (ui *Ui) resetCommandPrompt(v *gocui.View, fill string) {
	_, height := v.Size()
	if height-1 < 0 {
		return
	}

	v.SetCursor(0, height-1)

	for _, c := range ui.theme.CmdPrompt() {
		v.EditWrite(c)
	}

	for _, c := range fill {
		v.EditWrite(c)
	}
}

func (ui *Ui) redrawCommandsView(v *gocui.View, fill string) {
	// Subtract off a line for the current prompt and a row of separation
	_, height := v.Size()
	height -= 2
	if height < 0 {
		return
	}

	if v == nil {
		return
	}
	v.Clear()

	start := 0
	count := ui.hist.Size()

	if count > height {
		start = count - height
	}

	for i := start; i < count; i++ {
		v.Write([]byte(ui.hist.EntryString(i)))
	}

	ui.resetCommandPrompt(v, fill)
}

func (ui *Ui) updateCommandsView(v *gocui.View) error {
	if !v.Editable {
		v.Editable = true
		v.Editor = gocui.EditorFunc(ui.handleInput)
		ui.redrawCommandsView(v, "")
	}

	return nil
}
