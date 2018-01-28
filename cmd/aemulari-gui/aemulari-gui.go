package main

import (
	"fmt"
	"os"

	"../../internal/cmdline"
	"../../internal/debugger"
	"../../internal/log"
	"../../internal/ui"
)

func main() {
	var gui *ui.Ui
	var dbg *debugger.Debugger
	var err error

	log.Init()

	cmdline.InitCommonFlags()

	cfg, err := cmdline.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	dbg, err = debugger.New(cfg)
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
