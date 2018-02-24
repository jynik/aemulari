package ui

import (
	"github.com/jroimartin/gocui"

	ae "../../../aemulari.v0"
	"./theme"
)

type Ui struct {
	g     *gocui.Gui
	views Views

	dbg *ae.Debugger
	pc  uint64

	disasm DisassemblyInfo
	regs   RegInfo
	flags  FlagInfo
	mem    MemInfo
	hist   CommandHistory

	theme theme.Theme

	quit bool
}

func Create(arch *ae.Architecture, dbg *ae.Debugger) (*Ui, error) {
	var ui Ui

	regs, err := dbg.ReadRegAll()
	if err != nil {
		return nil, err
	}

	ui.initializeViews("%08x", len(regs))

	ui.theme, err = theme.New("default", (*arch).RegisterRegexp())
	if err != nil {
		return nil, err
	}

	ui.g, err = gocui.NewGui(gocui.Output256)
	if err != nil {
		return nil, err
	}

	ui.g.Cursor = true
	ui.g.SetManagerFunc(ui.update)

	err = ui.g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, ui.quitRequest)
	if err != nil {
		return nil, err
	}

	ui.dbg = dbg
	if reg, err := ui.dbg.ReadRegByName("pc"); err != nil {
		return nil, err
	} else {
		ui.pc = reg.Value
	}

	ui.mem.addr = ui.pc
	ui.update(ui.g)

	return &ui, nil
}

func (ui *Ui) quitRequest(gui *gocui.Gui, view *gocui.View) error {
	ui.quit = true
	return gocui.ErrQuit
}

func (ui *Ui) Close() {
	ui.g.Close()
}

func (ui *Ui) Run() error {
	if err := ui.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (ui *Ui) update(gui *gocui.Gui) error {
	for _, view := range ui.views {
		view.Update(gui)
	}

	if ui.quit {
		return gocui.ErrQuit
	} else {
		gui.SetCurrentView(vCommands)
		return nil
	}
}
