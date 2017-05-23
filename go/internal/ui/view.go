package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type View struct {
	// View name
	name string

	// Normalized X offset in [0, 1.0]. (% of width)
	x float32

	/* Pin to right of specified window.
	 * Specify this instead of x, or leave empty.  */
	rightOf string

	/* Normalized width in [0, 1.0] If rightOf is used, this is the percentage
	 * of the remaining width.  Otherwise, it's the percentage with respect to
	 * the entire GUI width. */
	width float32

	// Preferred absolute width, in characters
	prefWidth float32

	// Normalized Y Offset in [0, 1.0]. (% of height)
	y float32

	/* Pin to bottom of specified window. Specify this
	 * instead of y or leave empty. */
	below string

	// Normalized height in [0, 1.0]
	height float32

	// Preferred absolute height, in characters
	prefHeight float32

	// Parent UI
	ui *Ui

	// "Redraw" update callback
	UpdateCb ViewUpdateCb
}

type Views map[string]View
type ViewUpdateCb func(view *gocui.View) error

// View names
const vDisasm = "{ Disassembly }"
const vReg = "{ Registers }"
const vFlags = "{ Flags } "
const vMem = "{ Memory }"
const vConsole = "{ Console }"
const vCommands = "{ Commands }"

/* Set in View.width and View.height  to indicate the view should fill the
 * remaining width or height of the GUI region. */
const fillRemaining = 0.0

func (v View) guiWidth() float32 {
	w, _ := v.ui.g.Size()
	return float32(w)
}

func (v View) guiHeight() float32 {
	_, h := v.ui.g.Size()
	return float32(h)
}

func (v View) calculateWidth(x1 float32) float32 {
	var width float32

	if v.prefWidth != 0 {
		width = v.prefWidth
	} else if v.width == fillRemaining {
		width = v.guiWidth() - x1
	} else if v.rightOf != "" {
		if _, _, x, _, err := v.ui.g.ViewPosition(v.rightOf); err == nil {
			width = v.width * (v.guiWidth() - float32(x+1))
		}
	} else {
		width = v.width * v.guiWidth()
	}

	return width
}

func (v View) calculateHeight(y1 float32) float32 {
	var height float32
	if v.prefHeight != 0 {
		height = v.prefHeight + 2
	} else if v.height == fillRemaining {
		height = v.guiHeight() - y1
	} else if v.below != "" {
		if _, _, _, y, err := v.ui.g.ViewPosition(v.below); err == nil {
			height = v.height * (v.guiHeight() - float32(y-1))
		}
	} else {
		height = v.height * v.guiHeight()
	}

	return height
}

func (v View) calculateX1(gui *gocui.Gui) float32 {
	x1 := v.x * v.guiWidth()
	if v.rightOf != "" {
		if _, _, x, _, err := gui.ViewPosition(v.rightOf); err == nil {
			x1 = float32(x + 1)
		}
	}
	return x1
}

func (v View) calculateY1(gui *gocui.Gui) float32 {
	y1 := v.y * v.guiHeight()
	if v.below != "" {
		if _, _, _, y, err := gui.ViewPosition(v.below); err == nil {
			y1 = float32(y + 1)
		}
	}
	return y1
}

func (v *View) Update(gui *gocui.Gui) error {
	x1 := v.calculateX1(gui)
	x2 := x1 + v.calculateWidth(x1)
	y1 := v.calculateY1(gui)
	y2 := y1 + v.calculateHeight(y1)

	updatedView, err := gui.SetView(v.name, int(x1), int(y1), int(x2)-1, int(y2)-1)

	if err != nil {
		return fmt.Errorf("Failed to SetView(%s, ...): %s", v.name, err)
	} else {
		updatedView.Title = v.name
		if v.UpdateCb != nil {
			v.UpdateCb(updatedView)
		}
	}

	return nil
}

func (ui *Ui) initializeViews(addrFmt string, numRegs int) {
	testStr := fmt.Sprintf(addrFmt, 0)
	leftSideMaxWidth := 72 + len(testStr)
	flagsWidth := 30
	regWidth := leftSideMaxWidth - flagsWidth

	regHeight := float32((numRegs + 1) / 2)

	ui.views = Views{

		vFlags: View{
			name:       vFlags,
			x:          0.0,
			y:          0.0,
			prefWidth:  float32(flagsWidth),
			prefHeight: regHeight,
			UpdateCb:   ui.updateFlagsView,
			ui:         ui,
		},

		vReg: View{
			name:       vReg,
			rightOf:    vFlags,
			y:          0.0,
			prefWidth:  float32(regWidth),
			prefHeight: regHeight,
			UpdateCb:   ui.updateRegView,
			ui:         ui,
		},

		vMem: View{
			name:      vMem,
			x:         0.0,
			below:     vReg,
			prefWidth: float32(leftSideMaxWidth),
			height:    0.70,
			UpdateCb:  ui.updateMemView,
			ui:        ui,
		},

		vCommands: View{
			name:      vCommands,
			x:         0.0,
			below:     vMem,
			prefWidth: float32(leftSideMaxWidth),
			height:    fillRemaining,
			UpdateCb:  ui.updateCommandsView,
			ui:        ui,
		},

		vDisasm: View{
			name:     vDisasm,
			rightOf:  vMem,
			y:        0.0,
			width:    fillRemaining,
			height:   0.65,
			UpdateCb: ui.updateDisasmView,
			ui:       ui,
		},

		vConsole: View{
			name:     vConsole,
			rightOf:  vCommands,
			below:    vDisasm,
			width:    fillRemaining,
			height:   fillRemaining,
			UpdateCb: ui.updateConsoleView,
			ui:       ui,
		},
	}
}
