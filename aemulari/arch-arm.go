package aemulari

import (
	"fmt"
	cs "github.com/bnagy/gapstone"
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

type archArm struct {
	archBase
}

// Defined per usercorn/qemu/target-arm/cpu.h
const (
	_             = iota
	arm_excp_udef // Undefined instruction
	arm_excp_swi  // Software interrupt (SVC / SWI)
	arm_excp_prefetch_abort
	arm_excp_data_abort
	arm_excp_irq
	arm_excp_fiq
	arm_excp_bkpt        // Software breakpoint (bkpt #imm)
	arm_excp_exit        // QEMU - intercept return from v7M exception
	arm_excp_kernel_trap // QEMU - intercept kernel commpage access
	arm_excp_strex       // QEMU - intercept strex
	arm_excp_hyp_call    // Hypervisor call
	arm_excp_hyp_trap    // Hypervisor trap
	arm_excp_smc         // Secure mode call
	arm_excp_virq
	arm_excp_vfiq
)

/* Unicorn/QEMU ARM interrupt number to brief string description
 *		See unicorn/qemu/taget-arm/internals.h
 */
var excpStr map[uint32]string = map[uint32]string{
	arm_excp_udef:           "Undefined Instruction",
	arm_excp_swi:            "Software Interrupt",
	arm_excp_prefetch_abort: "Prefetch Abort",
	arm_excp_data_abort:     "Data Abort",
	arm_excp_irq:            "IRQ",
	arm_excp_fiq:            "FIQ",
	arm_excp_bkpt:           "Breakpoint",
	arm_excp_exit:           "Emulator v7M exception exit",
	arm_excp_kernel_trap:    "Emulator interception of kernel commpage",
	arm_excp_strex:          "Emulator interception of strex",
	arm_excp_hyp_call:       "Hypervisor Call",
	arm_excp_hyp_trap:       "Hypervisor Trap",
	arm_excp_smc:            "Secure Monitor Call",
	arm_excp_virq:           "Virtual IRQ",
	arm_excp_vfiq:           "Virtual FIQ",
}

// Per: http://infocenter.arm.com/help/index.jsp?topic=/com.arm.doc.dui0473m/dom1359731136117.html

var arm_r0 RegisterDef = RegisterDef{
	name: "r0",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R0,
}

var arm_r1 RegisterDef = RegisterDef{
	name: "r1",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R1,
}

var arm_r2 RegisterDef = RegisterDef{
	name: "r2",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R2,
}

var arm_r3 RegisterDef = RegisterDef{
	name: "r3",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R3,
}

var arm_r4 RegisterDef = RegisterDef{
	name: "r4",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R4,
}

var arm_r5 RegisterDef = RegisterDef{
	name: "r5",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R5,
}

var arm_r6 RegisterDef = RegisterDef{
	name: "r6",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R6,
}

var arm_r7 RegisterDef = RegisterDef{
	name: "r7",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R7,
}

var arm_r8 RegisterDef = RegisterDef{
	name: "r8",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R8,
}

var arm_r9 RegisterDef = RegisterDef{
	name: "r9",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R9,
}

var arm_r10 RegisterDef = RegisterDef{
	name: "r10",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R10,
}

var arm_r11 RegisterDef = RegisterDef{
	name: "r11",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R11,
}

var arm_r12 RegisterDef = RegisterDef{
	name: "r12",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R12,
}

var arm_r13 RegisterDef = RegisterDef{
	name: "sp",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R13,
}

var arm_r14 RegisterDef = RegisterDef{
	name: "lr",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R14,
}

var arm_r15 RegisterDef = RegisterDef{
	name: "pc",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_R15,
	pc:   true,
}

var arm_cpsr RegisterDef = RegisterDef{
	name: "cpsr",
	mask: 0xffffffff,
	fmt:  "0x%08x",
	uc:   uc.ARM_REG_CPSR,
	Flags: []Flag{
		{
			name: "N",
			desc: "Negative: 1 = result was negative, 0 = result was positive",
			lsb:  31,
			mask: (1 << 31),
			fmt:  "%d",
		},

		{
			name: "Z",
			desc: "Zero: 1 = result was 0, 0 = nonzero result",
			lsb:  30,
			mask: (1 << 30),
			fmt:  "%d",
		},

		{
			name: "C",
			desc: "Carry: 1 = carry in last operation, 0 = No carry",
			lsb:  29,
			mask: (1 << 29),
			fmt:  "%d",
		},

		{
			name: "V",
			desc: "Overflow: 1 = overflow in last operation, 0 = no overflow",
			lsb:  28,
			mask: (1 << 28),
			fmt:  "%d",
		},

		{
			name: "Q",
			desc: "Underflow: 1= underflow in last operation, 0 = no underflow",
			lsb:  27,
			mask: (1 << 27),
			fmt:  "%d",
		},

		{
			name: "J",
			desc: "Jazelle: 1 = Jazelle state, 0 = ARM/Thumb",
			lsb:  24,
			mask: (1 << 24),
			fmt:  "%d",
		},

		{
			name: "GE",
			desc: "Greater Than or Equal (SIMD): 1's denote result >=, 0's denote result <",
			lsb:  16,
			mask: (0xf << 16),
			fmt:  "0x%x",
		},

		{
			name: "E",
			desc: "Data Endianness: 1 = big, 0 = little",
			lsb:  9,
			mask: (1 << 9),
			fmt:  "%d",
		},

		{
			name: "A",
			desc: "Abort: 1 = disable imprecise aborts, 0 = aborts enabled",
			lsb:  8,
			mask: (1 << 8),
			fmt:  "%d",
		},

		{
			name: "I",
			desc: "IRQ Disable: 1 = interrupts disabled, 0 = interrupts enabled",
			lsb:  7,
			mask: (1 << 7),
			fmt:  "%d",
		},

		{
			name: "F",
			desc: "FIQ Disable: 1 = FIQ interrupts disabled, 0 = FIQ interrupts enabled",
			lsb:  6,
			mask: (1 << 6),
			fmt:  "%d",
		},

		{
			name: "T",
			desc: "Thumb State: 1 = Thumb state, 0 = ARM state",
			lsb:  5,
			mask: (1 << 5),
			fmt:  "%d",
		},

		{
			name: "M",
			desc: "Mode: 0x10 = User, 0x11 = FIQ, 0x12 = IRQ, 0x13 = Supervisor, 0x17 = Abort, 0x1b = Undefined, 0x1f = System",
			lsb:  0,
			mask: 0x1f,
			fmt:  "0x%02x",
		},
	},
}

func armConstructor(mode string) (Architecture, error) {
	var modeInfo processorMode

	switch mode {
	case "arm", "":
		modeInfo = processorMode{uc.MODE_ARM, cs.CS_MODE_ARM}
	case "thumb", "thumb2":
		modeInfo = processorMode{uc.MODE_THUMB, cs.CS_MODE_THUMB}
	default:
		return nil, fmt.Errorf("Invalid Arm mode specified (\"%s\")", mode)
	}

	arm := &archArm{
		archBase{
			processor:   processorType{uc.ARCH_ARM, cs.CS_ARCH_ARM},
			mode:        modeInfo,
			maxInstrLen: 4,
		},
	}

	arm.RegisterMap.add([]string{"r0", "a1"}, &arm_r0)
	arm.RegisterMap.add([]string{"r1", "a2"}, &arm_r1)
	arm.RegisterMap.add([]string{"r2", "a3"}, &arm_r2)
	arm.RegisterMap.add([]string{"r3", "a4"}, &arm_r3)
	arm.RegisterMap.add([]string{"r4", "v1"}, &arm_r4)
	arm.RegisterMap.add([]string{"r5", "v2"}, &arm_r5)
	arm.RegisterMap.add([]string{"r6", "v3"}, &arm_r6)
	arm.RegisterMap.add([]string{"r7", "v4"}, &arm_r7)
	arm.RegisterMap.add([]string{"r8", "v5"}, &arm_r8)
	arm.RegisterMap.add([]string{"r9", "v6", "sb"}, &arm_r9)
	arm.RegisterMap.add([]string{"r10", "v7", "sl"}, &arm_r10)
	arm.RegisterMap.add([]string{"r11", "v8", "fp"}, &arm_r11)
	arm.RegisterMap.add([]string{"r12", "ip"}, &arm_r12)
	arm.RegisterMap.add([]string{"sp", "r13"}, &arm_r13)
	arm.RegisterMap.add([]string{"lr", "r14"}, &arm_r14)
	arm.RegisterMap.add([]string{"pc", "r15"}, &arm_r15)
	arm.RegisterMap.add([]string{"cpsr", "r16"}, &arm_cpsr)

	return arm, nil
}

func (a *archArm) initialPC(pc uint64) (uint64) {
	switch a.mode.Uc {
	case uc.MODE_ARM:
		return pc
	case uc.MODE_THUMB:
		return pc | 0x1
	default:
		panic("Unexpected processor mode in arm.InitialPC()")
	}
}

func (a *archArm) currentPC(pc uint64, rvs []RegisterValue) uint64 {
	for _, rv := range rvs {
		if rv.Reg.name == "cpsr" {
			// Test THUMB bit
			if (rv.Value & (1 << 5)) != 0 {
				return pc | 0x1, nil
			}

			return pc
		}
	}

	panic("arm.CurrentPC() was not passed CPSR.")
}

func (a *archArm) endianness(rvs []RegisterValue) Endianness {
	for _, rv := range rvs {
		if rv.Reg.name == "cpsr" {
			// Test CPSR Endianness-bit
			if (rv.Value & (1 << 9)) != 0 {
				return BigEndian
			}

			return LittleEndian
		}
	}

	panic("arm.Endianness() was not passed CPSR.")
}

func (a *archArm) exception(intno uint32, regs []RegisterValue, instr []byte) (ex exception) {
	// TODO: Check if we're in Thumb mode
	thumb := false
	havePc := false

	if (thumb && len(instr) < 2) || (!thumb && len(instr) < 4) {
		panic("Arm.Exception was provided too few instruction bytes.")
	}

	for _, r := range regs {
		if r.Reg.Name() == "pc" {
			havePc = true
			ex.pc = r.Value
			break
		}
	}

	if !havePc {
		panic("pc was not in the register set provided to Arm.Exception()")
	}

	ex.intno = intno

	switch intno {
	case arm_excp_bkpt:
		var bkpt uint
		if thumb {
			bkpt = uint(instr[0])
		} else {
			bkpt = (uint(instr[2]&0xf) << 12) | (uint(instr[1]) << 4) | uint(instr[0]&0xf)
		}
		ex.desc = fmt.Sprintf("%s #0x%04x (%d)", excpStr[intno], bkpt, bkpt)

	default:
		if str, found := excpStr[intno]; found {
			ex.desc = str
		} else {
			ex.desc = fmt.Sprintf("Unknown exception (%d) occurred at pc=0x%08x", intno, ex.pc)
		}
	}

	return ex
}
