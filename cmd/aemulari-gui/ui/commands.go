package ui

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

type cmd struct {
	names []string
	min   int
	max   int
	exec  func(*Ui, cmd, []string) (string, error)
	help  string

	suppressHistory bool // Do not pace this command in history
	mayTaintRegs    bool
	mayTaintMem     bool

	matchedName string // Filled in by lookupCmd to denote which
	// command alias was matched
}

func lookupCmd(toFind string) (cmd, error) {
	c := strings.ToLower(toFind)

	for _, entry := range cmdList {
		for _, name := range entry.names {
			if strings.HasPrefix(name, c) {
				// XXX: Hack to avoid initialization loop
				if entry.names[0] == "help" && entry.exec == nil {
					entry.exec = cmdHelp
				}

				entry.matchedName = name
				return entry, nil
			}
		}
	}

	return cmd{}, fmt.Errorf("No such command: %s", toFind)
}

func (ui *Ui) handleCommand(line string) (string, bool, error) {
	args := strings.Split(line, " ")

	cmd, err := lookupCmd(args[0])
	if err != nil {
		return "", false, err
	}

	nargs := len(args)
	if nargs < cmd.min {
		return "", true,
			fmt.Errorf("%s: Too few arguments. At least %d are required.",
				cmd.matchedName, cmd.min)
	} else if nargs > cmd.max {
		return "", true,
			fmt.Errorf("%s: Too many arguments. The max is %d.",
				cmd.matchedName, cmd.max)
	}

	if cmd.mayTaintRegs {
		ui.regs.tainted = true
	}

	if cmd.mayTaintMem {
		ui.mem.tainted = true
	}

	output, err := cmd.exec(ui, cmd, args)
	if err != nil {
		if err == gocui.ErrQuit {
			return "", false, err
		}
		err = fmt.Errorf("%s: %s", cmd.matchedName, err.Error())
	}

	if rv, err := ui.dbg.ReadRegByName("pc"); err != nil {
		err = fmt.Errorf("Failed to re-read program counter.")
	} else {
		ui.pc = rv.Value
	}

	return output, cmd.suppressHistory, err
}

func lowerTrim(s string) string {
	return strings.ToLower(strings.Trim(s, " \t\r\n\x00"))
}

func matches(toMatch, s string) bool {
	return strings.HasPrefix(toMatch, lowerTrim(s))
}

/*******************************************************************************
 * Command table
 ******************************************************************************/

// Keep this sorted in the order of preferred completion and help display order
var cmdList []cmd = []cmd{
	{
		names:        []string{"regwrite", "rw"},
		min:          3,
		max:          3,
		exec:         cmdRegWrite,
		mayTaintRegs: true,
		help: "<regsister> <value>\n" +
			"Write <value> to <register>.",
	},

	{
		names:       []string{"memwrite", "mw"},
		min:         3,
		max:         4096, // Arbitrary "good enough" value
		exec:        cmdMemWrite,
		mayTaintMem: true,
		help: "<address> <value> [value] ... [value]\n" +
			"Write one or more values, converted to target endianness, to <address>.\n" +
			"<value> may be one of:\n" +
			"  A base 10 or base 16 value. This may be positive or negative.\n" +
			"  {<hex sequence>} such as: {0102deadbeef0405}\n" +
			"  A fixed-length signed or unsigned value via fn(<x>) where fn is:\n" +
			"    i8() u8(), u16(), i16(), i32(), u32(), i64(), u64()\n",
	},

	// TODO: memdump <region> <filename>

	// TODO: memload <region> <filename>
	{
		names:        []string{"continue", "run"},
		min:          1,
		max:          1,
		exec:         cmdContinue,
		mayTaintRegs: true,
		mayTaintMem:  true,
		help:         "\nContinue execution until a breakpoint or exception occurs.\n",
	},

	{
		names:        []string{"step"},
		min:          1,
		max:          2,
		exec:         cmdStep,
		mayTaintRegs: true,
		mayTaintMem:  true,
		help:         "[count]\nExecute a single or [count] instructions.",
	},

	{
		names: []string{"breakpoint"},
		min:   1,
		max:   2,
		exec:  cmdBreak,
		help: "[address]\n" +
			"Set a breakpoint at PC or [address], if specified.",
	},

	// Show breakpoints
	{
		names:       []string{"show", "display"},
		min:         2,
		max:         3,
		exec:        cmdShow,
		mayTaintMem: true,
		help: "<item> [per-item args]\nShow or display the specified <item>.\n" +
			"Available items:\n" +
			"	breakpoints" +
			"	memory <address>",
	},

	{
		names: []string{"delete"},
		min:   1,
		max:   3,
		exec:  cmdDelete,
		help: "[all | [<id|address> <value>]]\n" +
			" - With no arguments, this deletes any breakpoints at PC.\n" +
			" - If run with \"all\", all breakpoints are removed.\n" +
			" - Providing <id> and a <value> removes the associated breakpoint.\n" +
			" - Specifying <address> and <value> removes all breakpoints at <address>.\n",
	},

	{
		names:           []string{"clear"},
		min:             1,
		max:             2,
		exec:            cmdClear,
		suppressHistory: true,
		help:            "[console|commands]\nClear Console, Commands, or both.",
	},

	{
		names: []string{"reset"},
		min:   1,
		max:   1,
		exec:  cmdReset,
		help:  "\nResets the debugger. Mappings and breakpoints are kept as-is.",
	},

	{
		names:           []string{"quit", "exit"},
		min:             1,
		max:             1,
		exec:            cmdQuit,
		suppressHistory: true,
		help:            "\nExit the application. Alternatively, use Ctrl-Q.",
	},

	{
		names: []string{"help", "halp!", "wtf"},
		min:   1,
		max:   2,
		/* exec assigned later to avoid initialization loop */

		help: "<command>\nShow the help text for <command>\n",
	},
}

/*******************************************************************************
 * Command implementations
 ******************************************************************************/

// Keep these alphabetical, please!

func cmdBreak(ui *Ui, cmd cmd, args []string) (string, error) {
	var addr uint64
	var err error

	if len(args) < 2 {
		addr = ui.pc
	} else {
		addr, err = strconv.ParseUint(args[1], 0, 64)
		if err != nil {
			return "", err
		}
	}

	bp := ui.dbg.SetBreakpoint(addr)

	// FIXME use dbg-supplied address format
	return fmt.Sprintf("Added breakpoint %d at 0x%08x", bp.ID, bp.Address), nil
}

func cmdContinue(ui *Ui, cmd cmd, args []string) (string, error) {
	exception, err := ui.dbg.Continue()
	if err != nil {
		return "", err
	}

	if exception.Occurred() {
		return "Halted due to exception: " + exception.String(), nil
	}
	return "", nil
}

func cmdClear(ui *Ui, cmd cmd, args []string) (string, error) {
	cleared := false
	alen := len(args)

	vCmd, err := ui.g.View(vCommands)
	if err != nil {
		return "", err
	}

	if alen < 2 || strings.HasPrefix("console", lowerTrim(args[1])) {
		if v, err := ui.g.View(vConsole); err != nil {
			return "", err
		} else {
			v.Clear()
			vCmd.Clear()
			cleared = true
		}
	}

	if alen < 2 || strings.HasPrefix("commands", lowerTrim(args[1])) {
		ui.hist = CommandHistory{}
		vCmd.Clear()
		cleared = true
	}

	if !cleared {
		return "", fmt.Errorf("\"%s\" is not a valid argument.", args[1])
	}

	return "", nil
}

func cmdDelete(ui *Ui, cmd cmd, args []string) (string, error) {
	if len(args) == 1 {
		ui.dbg.DeleteBreakpointsAt(ui.pc)
		// FIXME need arch-specific address format
		return fmt.Sprintf("Removed breakpoints at 0x%08x.", ui.pc), nil

	} else if len(args) == 2 && matches("all", args[1]) {
		ui.dbg.DeleteAllBreakpoints()
		return "Removed all breakpoints.", nil

	} else if len(args) == 3 && matches("address", args[1]) {
		addr, err := strconv.ParseUint(args[2], 0, 64)
		if err != nil {
			return "", fmt.Errorf("\"s\" is not a valid breakpoint address.", args[2])
		}
		ui.dbg.DeleteBreakpointsAt(addr)

		// FIXME need arch-specific address format
		return fmt.Sprintf("Removed breakpoints at 0x%08x.", addr), nil

	} else if len(args) == 3 && matches("id", args[1]) {
		id, err := strconv.ParseInt(args[2], 0, 32)
		if err != nil {
			return "", fmt.Errorf("\"s\" is not a valid breakpoint ID.", args[2])
		}

		ui.dbg.DeleteBreakpoint(int(id))
		return fmt.Sprintf("Removed breakpoint %d.", id), nil
	} else {
		return "", errors.New("Invalid usage. See \"help delete\".")
	}
}

func cmdHelp(ui *Ui, cmd cmd, args []string) (string, error) {

	if len(args) < 2 {
		var helpText string = "Available commands:\n"
		for _, c := range cmdList {
			helpText += "  " + c.names[0] + "\n"
		}
		return helpText, nil
	} else {
		c, err := lookupCmd(args[1])
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s %s\n", c.matchedName, c.help), nil
	}
}

func cmdStep(ui *Ui, cmd cmd, args []string) (string, error) {
	var err error
	var count int64 = 1

	if len(args) > 1 {
		count, err = strconv.ParseInt(args[1], 0, 64)
		if err != nil || count <= 0 {
			return "", fmt.Errorf("\"%s\" is not a valid step size.", args[1])
		}
	}

	exception, err := ui.dbg.Step(count)
	if err != nil {
		return "", err
	}

	if exception.Occurred() {
		return "Halted due to exception: " + exception.String(), nil
	}

	if count == 1 {
		return fmt.Sprintf("Stepped 1 instruction."), nil
	} else {
		return fmt.Sprintf("Stepped %d instructions.", count), nil
	}
}

func cmdMemWrite(ui *Ui, cmd cmd, args []string) (string, error) {
	var data []byte

	endianness, err := ui.dbg.Endianness()
	if err != nil {
		return "", err
	}

	addr, err := strconv.ParseUint(args[1], 0, 64)
	if err != nil {
		return "", fmt.Errorf("\"%s\" is not a valid address.", args[1])
	}

	for _, value := range args[2:] {
		bytes, err := parseValue(value, endianness)
		if err != nil {
			return "", err
		}

		data = append(data, bytes...)
	}

	return "", ui.dbg.WriteMem(addr, data[0:])
}

func cmdQuit(ui *Ui, cmd cmd, args []string) (string, error) {
	ui.quit = true
	return "", nil
}

func cmdShow(ui *Ui, cmd cmd, args []string) (string, error) {
	what := lowerTrim(args[1])

	if strings.HasPrefix("breakpoints", what) {
		var ret string
		bps := ui.dbg.GetBreakpoints()
		for _, bp := range bps {
			ret += bp.String() + "\n"
		}
		return ret, nil
	} else if strings.HasPrefix("memory", what) {
		if len(args) < 3 {
			return "", fmt.Errorf("A required <address> argument was not provided.")
		}

		newAddr, err := strconv.ParseUint(args[2], 0, 64)
		if err != nil {
			return "", fmt.Errorf("\"%s\" is not a valid memory address.", args[2])
		}

		if view, err := ui.g.View(vMem); err != nil {
			return "", err
		} else {
			ui.mem.addr = newAddr
			return "", ui.updateMemView(view)
		}

		return "", nil
	}

	return "", fmt.Errorf("\"%s\" is not a valid item.", args[1])
}

func cmdRegWrite(ui *Ui, cmd cmd, args []string) (string, error) {
	if val, err := strconv.ParseUint(args[2], 0, 64); err != nil {
		return "", fmt.Errorf("\"%s\" is not a valid register value.", args[2])
	} else {
		return "", ui.dbg.WriteRegByName(args[1], val)
	}
}

func cmdReset(ui *Ui, cmd cmd, args []string) (string, error) {
	return "", ui.dbg.Reset(true)
}
