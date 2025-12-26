package emu

import "fmt"

// The jump function returns PC++ (no jump) or the new PC address (jump)
func ExecJump(cpu *HACKCPU, jumpbits string) (newpc uint16, err error) {

	// default: no jump
	newpc = cpu.PC + 1
	if !AaddrIsValid(int16(newpc)) || !AaddrIsValid(cpu.A) {
		err = fmt.Errorf("address overflow! A=%4X | PC=%04X", cpu.A, newpc)
		return
	}

	ZERO := cpu.SR.Get(ZFLAG)
	NGTV := cpu.SR.Get(NFLAG)

	switch jumpbits {
	case "000": // nojump
		return
	case "001": // "JGT":
		if !ZERO && !NGTV {
			newpc = uint16(cpu.A)
		}
	case "010": // "JEQ":
		if ZERO {
			newpc = uint16(cpu.A)
		}
	case "011": // "JGE":
		if !NGTV {
			newpc = uint16(cpu.A)
		}
	case "100": // "JLT":
		if !ZERO && NGTV {
			newpc = uint16(cpu.A)
		}
	case "101": // "JNE":
		if !ZERO {
			newpc = uint16(cpu.A)
		}
	case "110": // "JLE":
		if ZERO || NGTV {
			newpc = uint16(cpu.A)
		}
		// "111": "JMP":
	case "111": // "JMP":
		newpc = uint16(cpu.A)
	}

	return
}
