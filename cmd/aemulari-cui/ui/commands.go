package ui

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"

	ae "../../../aemulari.v0"
	"github.com/jroimartin/gocui"
)

const linesep = "" +
"-------------------------------------------------------------------------------\n"

type cmd struct {
	names   []string
	min     int
	max     int
	exec    func(*Ui, cmd, []string) (string, error)
	summary string
	details string

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
		if cmd.min <= 2 {
			err = fmt.Errorf("%s: Required argument is missing.",
					cmd.matchedName)
		} else {
			err = fmt.Errorf("%s: Too few arguments. At least %d are required.",
					cmd.matchedName, cmd.min - 1)
		}
		return "", true, err
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
		names:        []string{"continue", "run"},
		min:          1,
		max:          1,
		exec:         cmdContinue,
		mayTaintRegs: true,
		mayTaintMem:  true,
		summary:	  "Execute until a breakpoint or exception occurs",
		details:	  "\n" +
		"\n" +
		"Execute until a breakpoint or exception occurs.",
	},

	{
		names:        []string{"step"},
		min:          1,
		max:          2,
		exec:         cmdStep,
		mayTaintRegs: true,
		mayTaintMem:  true,
		summary:	  "Execute 1 or more instructions",
		details:      "[count]\n" +
		"\n" +
		"Execute a single or [count] instructions.\n",
	},

	{
		names: []string{"breakpoint"},
		min:   1,
		max:   2,
		exec:  cmdBreak,
		summary: "Set a breakpoint",
		details: "[address]\n" +
			"\n" +
			"Set a breakpoint at PC or [address], if specified.\n",
	},

	{
		names: []string{"delete"},
		min:   1,
		max:   3,
		exec:  cmdDelete,
		summary: "Delete specified breakpoints",
		details: "[all | [id|address <value>]]\n" +
			"\n" +
			"Notes:" +
			" - With no arguments, this deletes any breakpoints at PC.\n" +
			" - If run with \"all\", all breakpoints are removed.\n" +
			" - Providing `id` and a <value> removes the associated breakpoint.\n" +
			" - Specifying `address` and <value> removes all breakpoints at <value>.\n",
	},


	{
		names:        []string{"rw"},
		min:          3,
		max:          3,
		exec:         cmdRegWrite,
		mayTaintRegs: true,
		summary: "Write a value to a register",
		details: "<register> <value>\n" +
			"Write <value> to <register>\n" +
			"\n" +
			"<value> may be one of:\n" +
			"  - A base 10 or base 16 value. This may be positive or negative.\n" +
			"  - {<hex sequence>} such as: {0102deadbeef0405}\n" +
			"  - A fixed-length signed or unsigned value via fn(<x>) where fn is:\n" +
			"     i8() u8(), u16(), i16(), i32(), u32(), i64(), u64()\n" +
			"\n" +
			"Examples:\n" +
			" rw r0 0x1b4d1dea\n" +
			" rw r0 i16(-7)\n" +
			" rw r0 {deadbeef}\n",
	},

	{
		names:       []string{"mw"},
		min:         3,
		max:         4096, // Arbitrary "good enough" value
		exec:        cmdMemWrite,
		mayTaintMem: true,

		summary: "Write data to the specified memory address",
		details: "<address> <value> [value] ... [value]\n" +
			"Write one or more values, converted to target endianness, to <address>.\n" +
			"\n" +
			"<value> may be one of:\n" +
			"  - A base 10 or base 16 value. This may be positive or negative.\n" +
			"  - {<hex sequence>} such as: {0102deadbeef0405}\n" +
			"  - A fixed-length signed or unsigned value via fn(<x>) where fn is:\n" +
			"     i8() u8(), u16(), i16(), i32(), u32(), i64(), u64()\n" +
			"\n" +
			"Examples:\n" +
			" mw 0x1ab000 0x1b4d1dea\n" +
			" mw 0x1ab000 i16(-7)\n" +
			" mw 0x1ab000 {deadbeef} 0xbadc0de\n",
	},

	{
		names: []string{"map"},
		min:		2,
		max:		7,
		exec:  cmdMemMap,
		mayTaintMem: true,
		summary: "Map and configure a memory region.",
		details: "<name>:<addr>:<size>:[permissions]:[input file]:[output_file]\n" +
			"\n" +
			"This command is identical to the -m/--mem <region> command line option,\n" +
			"and may be used to map, configure, and initialize a memory region.\n" +
			"However, in this UI, a space may be used instead of the ':' separator.\n" +
			"\n" +
			"If [output_file] is specified, the memory contents will be written\n" +
			"when the region is unmapped - either manually or when the debugger\n" +
			"is closed.\n",
	},

	{
		names: []string{"unmap"},
		exec:  cmdMemUnmap,
		min:		2,
		max:		4096,  // Arbitrary "good enough" value
		mayTaintMem: true,
		summary: "Unmap a memory region",
		details: "<name> [name] ... [name]\n" +
			"\n" +
			"Unmap one or more memory regions, by name.\n" +
			"\n" +
			"If a region was configured with an output file, the asscioted memory\n" +
			"will be written to this file before unmapping the region.\n",
	},

	{
		names:		[]string{"dumpmem", "savemem"},
		min:		3,
		max:		4,
		exec:		cmdDumpMem,
		summary:	"Save the contents of a memory region to a file",
		details:	"<filename> <region name>\n" +
			"               <filename> <addr> <length>\n" +
			"\n" +
			"Save the contents of a memory region, specified by name or\n" +
			"by address and length, to a file.",
	},


	{
		names:       []string{"display", "show"},
		min:         2,
		max:         3,
		exec:        cmdDisplay,
		mayTaintMem: true, // Not really taint, but we need to force redraw of Memory window
		summary: "Display information about the specified item(s)",
		details: "<item> [per-item args]\nShow or display the specified <item>.\n" +
			"\n" +
			"Available items:\n" +
			"	breakpoints" +
			"	memory <address>" +
			"	mapped [name]\n",
	},

	{
		names:           []string{"clear"},
		min:             1,
		max:             2,
		exec:            cmdClear,
		suppressHistory: true,
		summary:		 "Clear console or command window",
		details:         "[console|commands]\n" +
		"\n" +
		"Clear Console, Commands, or both.",
	},

	{
		names:   []string{"reset"},
		min:     1,
		max:     1,
		exec:    cmdReset,
		summary: "Reset the debugger state",
		details: "\n" +
		"\n" +
		"Resets the debugger. Memory mappings and breakpoints are kept as-is.",
	},

	{
		names:           []string{"quit", "exit"},
		min:             1,
		max:             1,
		exec:            cmdQuit,
		suppressHistory: true,
		summary:		 "Exit the program",
		details:         "\n" +
		"\n" +
		"Exit the program. Alternatively, use Ctrl-Q.",
	},

	{
		names: []string{"help"},
		min:   1,
		max:   2,
		/* exec assigned later to avoid initialization loop */

		summary: "Describe the specified command",
		details: "<command>\n" +
		"\n" +
		"Show the help text for <command>\n",
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
			v.SetOrigin(0, 0)
			cleared = true
		}
	}

	if alen < 2 || strings.HasPrefix("commands", lowerTrim(args[1])) {
		ui.hist = CommandHistory{}
		vCmd.Clear()
		vCmd.SetOrigin(0, 0)
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

func cmdDisplay(ui *Ui, cmd cmd, args []string) (string, error) {
	var ret string
	what := lowerTrim(args[1])

	if strings.HasPrefix("breakpoints", what) {
		bps := ui.dbg.GetBreakpoints()

		ret += "\nBreakpoints\n"
		ret += linesep

		for _, bp := range bps {
			ret += bp.String() + "\n"
		}
		return ret, nil

	} else if strings.HasPrefix("mapped", what) {
		ret += "\nMemory Mapped Regions\n"
		ret += linesep

		regions := ui.dbg.Mapped()
		for _, region := range regions {
			var haveMatch bool
			if len(args) > 2 {
				for _, name := range args[2:] {
					haveMatch = name == region.Name()
					if haveMatch {
						break
					}
				}
			} else {
				haveMatch = true
			}

			if haveMatch {
				ret += region.String() + "\n"
			}
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
			ui.mem.pdata = []byte{}

			// FIXME This shouldn't require a double-kick to prevent it from
			//		 incorrectly highlighting changes when we point the view at
			//		 a different memory location
			ui.updateMemView(view)
			ui.mem.pdata = []byte{}

			return "", ui.updateMemView(view)
		}
		return "", nil
	}

	return "", fmt.Errorf("\"%s\" is not a valid item.", args[1])
}

func cmdDumpMem(ui *Ui, cmd cmd, args []string) (string, error) {
	var addr, size uint64
	var name, filename string
	var err error

	filename = args[1]

	if len(args) == 3 {
		name = args[2]
		err = ui.dbg.DumpMemRegion(filename, name)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Wrote contents of \"%s\" region to %s\n", name, filename), nil

	}

	addr, err = strconv.ParseUint(args[2], 0, 64)
	if err != nil {
		return "", fmt.Errorf("Invalid address (%s)", args[2])
	}

	size, err = strconv.ParseUint(args[3], 0, 64)
	if err != nil {
		return "", fmt.Errorf("Invalid size (%s)", args[3])
	}

	err = ui.dbg.DumpMem(filename, addr, size)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Wrote contents of [0x%08x - 0x%08x] to %s\n", addr, addr+size-1, filename), nil
}

func cmdHelp(ui *Ui, cmd cmd, args []string) (string, error) {

	if len(args) < 2 {
		var helpText string

		helpText += "Available commands:\n"

		for _, c := range cmdList {
			helpText += fmt.Sprintf("  %-12s %s\n", c.names[0], c.summary)
		}

		helpText += "\n"
		helpText += "Notes:\n"
		helpText += " - Cycle through command history using the up and down keys.\n"
		helpText += " - Entering an empty command will run the previous command.\n"
		helpText += "     This is useful when stepping through a program.\n"
		helpText += " - Only a subset of command names is actually required.\n"
		helpText += "     Commands are matched in the order presented above.\n"

		return helpText, nil
	} else {
		c, err := lookupCmd(args[1])
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Usage: %s %s\n", c.matchedName, c.details), nil
	}
}


func cmdMemMap(ui *Ui, cmd cmd, args []string) (string, error) {
	regionStr := strings.Join(args[1:], ":")

	region, err := ae.NewMemRegion(regionStr)
	if err != nil {
		return "", err
	}

	return "", ui.dbg.Map(region)
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

func cmdMemUnmap(ui *Ui, cmd cmd, args []string) (string, error) {
	var ret string
	for _, name := range args[1:] {
		if name == "code" {
			ret += "Code region may not be unmapped. Use dumpmem to save it to disk."
			continue
		}

		err := ui.dbg.Unmap(name)
		if err != nil {
			return "", err
		}
	}

	return ret, nil
}



func cmdQuit(ui *Ui, cmd cmd, args []string) (string, error) {
	ui.quit = true
	return "", nil
}

func cmdRegWrite(ui *Ui, cmd cmd, args []string) (string, error) {
	var regVal uint64

	endianness, err := ui.dbg.Endianness()
	if err != nil {
		return "", err
	}

	bytes, err  := parseValue(args[2], endianness)
	if err != nil {
		return "", fmt.Errorf("\"%s\" is not a valid register value.", args[2])
	}

	if len(bytes) > 8 {
		return "", fmt.Errorf("\"%s\" exceeds the maximum register size.", args[2])
	} else if len(bytes) < 8 {
		padLen := 8 - len(bytes)
		padding := make([]byte, padLen)
		bytes = append(bytes, padding...)
	}

	if endianness == ae.BigEndian {
		regVal = binary.BigEndian.Uint64(bytes)
	} else {
		regVal = binary.LittleEndian.Uint64(bytes)
	}

	return "", ui.dbg.WriteRegByName(args[1], regVal)
}

func cmdReset(ui *Ui, cmd cmd, args []string) (string, error) {
	return "", ui.dbg.Reset(true)
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
