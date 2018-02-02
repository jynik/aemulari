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

	cfg, err := common.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	dbg, err = aemulari.NewDebugger(cfg)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if gui, err = ui.Create(dbg); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	gui.Run()
	gui.Close()
}
