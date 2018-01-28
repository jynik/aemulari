package ui

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/jroimartin/gocui"

	"./theme"
)

type MemInfo struct {
	addr uint64
	data []byte

	paddr uint64 // Previous address
	pdata []byte // Previous state of data

	tainted bool // Has data been potentially tainted?
}

func (ui *Ui) updateMemView(view *gocui.View) error {
	_, height := view.Size()

	if height < 1 {
		return errors.New("Memory view not large enough to draw")
	}

	dataLen := uint64(height * 16)
	if uint64(len(ui.mem.data)) == dataLen && ui.mem.tainted {
		if len(ui.mem.data) != len(ui.mem.pdata) {
			ui.mem.pdata = make([]byte, len(ui.mem.data))
		}

		copy(ui.mem.pdata, ui.mem.data)
		ui.mem.tainted = false
	}

	tmp, err := ui.dbg.ReadMem(ui.mem.addr, dataLen)
	if err != nil {
		return err
	}
	ui.mem.data = tmp

	if len(ui.mem.pdata) != len(ui.mem.data) {
		ui.mem.pdata = make([]byte, len(ui.mem.data))
		copy(ui.mem.pdata, ui.mem.data)
	}

	view.Clear()
	// TODO get fmt from ui.dbg
	fmt.Fprintf(view, "%s", ui.hexdump(ui.mem.addr, "%08x", ui.theme, ui.mem))

	return nil
}

// Create a hexdump (akin to hexdump -C) with a configurable address & format,
// and highlighting differences between data `x` and `other`. Set other to nil
// if there's no data to compare.
func (ui Ui) hexdump(addr uint64, addrFmt string, theme theme.Theme, m MemInfo) string {
	var dump string
	var i, j, count uint64

	count = uint64(len(m.data))

	for i = 0; i < count; i += 16 {

		// Address
		line := " " + ui.theme.ColorAddress(addrFmt, addr+i) + "  "

		// First set of 8 bytes
		for j := i; j < (i + 8); j += 1 {
			if j < count {
				str := fmt.Sprintf("%02x ", m.data[j])
				line += theme.ColorIfBytesDiffer(str, m.data[j], m.pdata[j])
			} else {
				line += "   "
			}
		}

		line += " "

		// Second set of 8 bytes
		for j := i + 8; j < (i + 16); j += 1 {
			if j < count {
				str := fmt.Sprintf("%02x ", m.data[j])
				line += theme.ColorIfBytesDiffer(str, m.data[j], m.pdata[j])
			} else {
				line += "   "
			}
		}

		line += " "

		for j = i; j < i+16; j += 1 {
			if j < count {
				if m.data[j] <= unicode.MaxASCII && unicode.IsPrint(rune(m.data[j])) {
					str := fmt.Sprintf("%c", m.data[j])
					line += theme.ColorIfBytesDiffer(str, m.data[j], m.pdata[j])
				} else {
					line += theme.ColorIfBytesDiffer(".", m.data[j], m.pdata[j])
				}
			}
		}

		line += "\n"
		dump += line
	}

	return dump
}
