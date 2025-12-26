package emu

import "fmt"

type InstructionFunc func(cpu *HACKCPU) (result int16, err error)

type InstructionFuncMap map[string]InstructionFunc

// Operations
var ISACOMP = InstructionFuncMap{
	"0101010": comp0101010, // "0":
	"0111111": comp0111111, // "1":
	"0111010": comp0111010, // "-1":
	"0001100": comp0001100, // "D":
	"0110000": comp0110000, // "A":
	"0001101": comp0001101, // "!D":
	"0110001": comp0110001, // "!A":
	"0001111": comp0001111, // "-D":
	"0110011": comp0110011, // "-A":
	"0011111": comp0011111, // "D+1":
	"0110111": comp0110111, // "A+1":
	"0001110": comp0001110, // "D-1":
	"0110010": comp0110010, // "A-1":
	"0000010": comp0000010, // "D+A":
	"0010011": comp0010011, // "D-A":
	"0000111": comp0000111, // "A-D":
	"0000000": comp0000000, // "D&A":
	"0010101": comp0010101, // "D|A":
	"1110000": comp1110000, // "M":
	"1110001": comp1110001, // "!M":
	"1110011": comp1110011, // "-M":
	"1110111": comp1110111, // "M+1":
	"1110010": comp1110010, // "M-1":
	"1000010": comp1000010, // "D+M":
	"1010011": comp1010011, // "D-M":
	"1000111": comp1000111, // "M-D":
	"1000000": comp1000000, // "D&M":
	"1010101": comp1010101, // "D|M":
	"1111111": comp1111111, // "BRK":
}

func ExecComp(cpu *HACKCPU, compbits string) (err error) {
	var result int16

	isacomp, found := ISACOMP[compbits]
	if !found {
		err = fmt.Errorf("no such comp function: %s", compbits)
		return
	}

	result, err = isacomp(cpu)
	cpu.ALU = result

	if cpu.ALU == 0 {
		cpu.SR.Set(ZFLAG)
	} else {
		cpu.SR.Clear(ZFLAG)
	}
	if cpu.ALU < 0 {
		cpu.SR.Set(NFLAG)
	} else {
		cpu.SR.Clear(NFLAG)
	}

	return
}

// "BRK": "1111111"
func comp1111111(cpu *HACKCPU) (result int16, err error) {
	cpu.RUNNING = false
	return
}

// "0":   "0101010",
func comp0101010(cpu *HACKCPU) (result int16, err error) {
	result = 0
	return
}

// "1":   "0111111",
func comp0111111(cpu *HACKCPU) (result int16, err error) {
	result = 1
	return
}

// "-1":  "0111010",
func comp0111010(cpu *HACKCPU) (result int16, err error) {
	result = -1
	return
}

// "D":   "0001100",
func comp0001100(cpu *HACKCPU) (result int16, err error) {
	result = cpu.D
	return
}

// "A":   "0110000",
func comp0110000(cpu *HACKCPU) (result int16, err error) {
	result = cpu.A
	return
}

// "!D":  "0001101",
func comp0001101(cpu *HACKCPU) (result int16, err error) {
	result = ^(cpu.D)
	return
}

// "!A":  "0110001",
func comp0110001(cpu *HACKCPU) (result int16, err error) {
	result = ^(cpu.A)
	return
}

// "-D":  "0001111",
func comp0001111(cpu *HACKCPU) (result int16, err error) {
	result = -(cpu.D)
	return
}

// "-A":  "0110011",
func comp0110011(cpu *HACKCPU) (result int16, err error) {
	result = -(cpu.A)
	return
}

// "D+1": "0011111",
func comp0011111(cpu *HACKCPU) (result int16, err error) {
	result = cpu.D + 1
	return
}

// "A+1": "0110111",
func comp0110111(cpu *HACKCPU) (result int16, err error) {
	result = cpu.A + 1
	return
}

// "D-1": "0001110",
func comp0001110(cpu *HACKCPU) (result int16, err error) {
	result = cpu.D - 1
	return
}

// "A-1": "0110010",
func comp0110010(cpu *HACKCPU) (result int16, err error) {
	result = cpu.A - 1
	return
}

// "D+A": "0000010",
func comp0000010(cpu *HACKCPU) (result int16, err error) {
	result = cpu.D + cpu.A
	return
}

// "D-A": "0010011",
func comp0010011(cpu *HACKCPU) (result int16, err error) {
	result = cpu.D - cpu.A
	return
}

// "A-D": "0000111",
func comp0000111(cpu *HACKCPU) (result int16, err error) {
	result = cpu.A - cpu.D
	return
}

// "D&A": "0000000",
func comp0000000(cpu *HACKCPU) (result int16, err error) {
	result = cpu.A & cpu.D
	return
}

// "D|A": "0010101",
func comp0010101(cpu *HACKCPU) (result int16, err error) {
	result = cpu.A | cpu.D
	return
}

// "M":   "1110000",
func comp1110000(cpu *HACKCPU) (result int16, err error) {
	result = cpu.ReadM()
	return
}

// "!M":  "1110001",
func comp1110001(cpu *HACKCPU) (result int16, err error) {
	result = ^(cpu.ReadM())
	return
}

// "-M":  "1110011",
func comp1110011(cpu *HACKCPU) (result int16, err error) {
	result = -(cpu.ReadM())
	return
}

// "M+1": "1110111",
func comp1110111(cpu *HACKCPU) (result int16, err error) {
	result = cpu.ReadM() + 1
	return
}

// "M-1": "1110010",
func comp1110010(cpu *HACKCPU) (result int16, err error) {
	result = cpu.ReadM() - 1
	return
}

// "D+M": "1000010",
func comp1000010(cpu *HACKCPU) (result int16, err error) {
	result = cpu.D + cpu.ReadM()
	return
}

// "D-M": "1010011",
func comp1010011(cpu *HACKCPU) (result int16, err error) {
	result = cpu.D - cpu.ReadM()
	return
}

// "M-D": "1000111",
func comp1000111(cpu *HACKCPU) (result int16, err error) {
	result = cpu.ReadM() - cpu.D
	return
}

// "D&M": "1000000",
func comp1000000(cpu *HACKCPU) (result int16, err error) {
	result = cpu.D & cpu.ReadM()
	return
}

// "D|M": "1010101",
func comp1010101(cpu *HACKCPU) (result int16, err error) {
	result = cpu.D | cpu.ReadM()
	return
}
