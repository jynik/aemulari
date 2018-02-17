package main

import (
	"fmt"
	"os"

	"../../aemulari"
	"../common"
	"./ui"
)

func main() {
	var gui *ui.Ui
	var dbg *aemulari.Debugger
	var err error

	common.InitCommonFlags()

	arch, cfg, err := common.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	dbg, err = aemulari.NewDebugger(arch, cfg)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if gui, err = ui.Create(arch, dbg); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	gui.Run()
	gui.Close()
}
