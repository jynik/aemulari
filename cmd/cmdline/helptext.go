package cmdline

// I'd rather have help text look the way I like it, rather than fight with a framework.

const FlagStr_arch = "" +
	"  -a, --arch <arch>           Architecture to emulate. (default: arm)\n"

const Details_arch = "" +
	"\nSupported Architectures and Initial Modes:\n" +
	"  arm          32-bit Arm\n" +
	"  arm:thumb    32-bit Arm in Thumb mode\n"

const FlagStr_mem = "" +
	"  -m, --mem <region>          Memory region to map and optionally load or dump.\n"

const Details_mem = "" +
	"\nMemory Mapped Regions:\n" +
	"  Memory mapped regions are specified using the following syntax:\n" +
	"\n" +
	"    <name>:<addr>:<size>:[perms]:[input file]:[output file]\n" +
	"\n" +
	"  - An executable region named \"code\" is required.\n" +
	"  - Only the <name>, <address>, and <size> fields are required. When including\n" +
	"      an optional field, all preceding optional fields must be specified well.\n" +
	"  - The [:perms] field specifies the access permissions of a region using the\n" +
	"      characters 'r', 'w', and 'x' for read, write, and execute, respectively.\n" +
	"  - If specified, the contents of [input file] will be loaded into the region.\n" +
	"  - If specified, a memory region's contents will be written to <output file>.\n"

const FlagStr_regs = "" +
	"  -r, --reg <name>=<value>    Assigns the initial value of a register.\n"

const FlagStr_breakpoint = "" +
	"  -b, --break <addr>          Set a breakpoint at the specified address.\n"

const FlagStr_printRegs = "" +
	"  -R, --print-regs [style]    Print registers after execution completes.\n" +
	"                               Style options: pretty (default), list\n"

const FlagStr_printHexdump = "" +
	"  -d, --hexdump <name>        Print a hexdump of a memory region after execution\n" +
	"                <addr:size>    completes. The region may be specified by name or \n" +
	"                               by an address and size.\n"

const FlagStr_help = "" +
	"  -h, --help                  Show this text and exit.\n"

const Notes = "" +
	"\nNotes:\n" +
	" - Numeric parameters may be specified in decimal format or in hex format,\n" +
	"     prefixed with \"0x\" (e.g., 0x1b4d1d3a).\n"
