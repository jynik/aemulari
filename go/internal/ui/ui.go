package ui

import (
	"github.com/jroimartin/gocui"

	"../arch"
	dbg "../debugger"
	"./theme"
)

type Ui struct {
	g     *gocui.Gui
	views Views

	dbg *dbg.Debugger
	pc  uint64

	disasm DisassemblyInfo
	regs   RegInfo
	flags  FlagInfo
	mem    MemInfo
	hist   CommandHistory

	theme theme.Theme

	quit bool
}

func Create(dbg *dbg.Debugger) (*Ui, error) {
	var ui Ui
	var err error
	var rv arch.RegisterValue
	var rvs []arch.RegisterValue

	rvs, err = dbg.ReadRegAll()
	if err != nil {
		return nil, err
	}

	ui.initializeViews("%08x", len(rvs))

	ui.theme, err = theme.New("default", dbg.RegisterRegexp())
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
	if rv, err = dbg.ReadRegByName("pc"); err != nil {
		return nil, err
	} else {
		ui.pc = rv.Value
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
