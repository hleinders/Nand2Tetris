package emu

import (
	"strconv"

	"github.com/hleinders/hack/tools"
)

const (
	MBit = 0b0001
	DBit = 0b0010
	ABit = 0b0100
)

func ExecDest(cpu *HACKCPU, destbits string) (err error) {
	destval, err := strconv.ParseInt(destbits, 2, 8)

	tools.Debug("Bitstr: %s --> Bitval %x\n", destbits, destval)

	if (destval & ABit) != 0 {
		tools.Debug("%03b: A-Register <-- ALU (%x)\n", destval, cpu.ALU)
		cpu.UpdateA(cpu.ALU)
	}

	if (destval & DBit) != 0 {
		tools.Debug("%03b: D-Register <-- ALU (%x)\n", destval, cpu.ALU)
		cpu.UpdateD()
	}

	if (destval & MBit) != 0 {
		tools.Debug("%03b: M-Register <-- ALU (%x)\n", destval, cpu.ALU)
		cpu.UpdateM()
	}

	return
}
