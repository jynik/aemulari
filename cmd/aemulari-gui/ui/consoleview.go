package ui

import "github.com/jroimartin/gocui"

func (ui *Ui) writeConsole(text string) {
	v, err := ui.g.View(vConsole)
	if v != nil && err == nil {
		v.Clear()
		v.Write([]byte(text))
	}
}

func (ui *Ui) appendConsole(text string) {
	v, err := ui.g.View(vConsole)
	if v != nil && err == nil {
		v.Write([]byte(text))
	}
}

func (ui *Ui) updateConsoleView(v *gocui.View) error {
	if !v.Wrap {
		v.Wrap = true
		v.Autoscroll = true
	}
	return nil
}
