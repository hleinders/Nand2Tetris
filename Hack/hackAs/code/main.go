package code

// With many thanks fetched and modified from https://github.com/sato11/the-hack-assembler

import (
	"fmt"
)

// Code is an interface to which modular functionality is provided.
type Code struct{}

var destinations = map[string]string{
	"null": "000",
	"0":    "000",
	"M":    "001",
	"D":    "010",
	"MD":   "011",
	"A":    "100",
	"AM":   "101",
	"AD":   "110",
	"AMD":  "111",
}

var compute = map[string]string{
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"-D":  "0001111",
	"-A":  "0110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"D+A": "0000010",
	"D-A": "0010011",
	"A-D": "0000111",
	"D&A": "0000000",
	"D|A": "0010101",
	"M":   "1110000",
	"!M":  "1110001",
	"-M":  "1110011",
	"M+1": "1110111",
	"M-1": "1110010",
	"D+M": "1000010",
	"D-M": "1010011",
	"M-D": "1000111",
	"D&M": "1000000",
	"D|M": "1010101",
	// Undocumented:
	"D!&A": "0000001", // NAND
	"D!&M": "1000001",
	"D!|A": "0010100", // NOR
	"D!|M": "1010100",
}

var jumps = map[string]string{
	"null": "000",
	"0":    "000",
	"JGT":  "001",
	"JEQ":  "010",
	"JGE":  "011",
	"JLT":  "100",
	"JNE":  "101",
	"JLE":  "110",
	"JMP":  "111",
}

// New returns the api interface.
func New() *Code {
	return &Code{}
}

// Dest code iss the "binary" code of the dest mnemonic.
func (c *Code) Dest(token string) string {
	var opcode string

	if opcode = destinations[token]; opcode != "" {
		return opcode
	}

	return "000"
}

// Comp code is the "binary" code of the comp mnemonic.
func (c *Code) Comp(token string) (string, error) {
	var opcode string

	if opcode = compute[token]; opcode != "" {
		return opcode, nil
	}

	return "0000000", fmt.Errorf("unknown compute instruction '%s'", token)

}

// Jump code iss the "binary" code of the jump mnemonic.
func (c *Code) Jump(token string) string {
	var opcode string

	if opcode = jumps[token]; opcode != "" {
		return opcode
	}

	return "000"
}
