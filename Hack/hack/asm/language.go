/*
Copyright © 2025 harald@leinders.de
*/

package asm

// Hack language specifications

import (
	"fmt"
	"maps"
	"os"
	"sort"
	"text/tabwriter"
)

// Global
const (
	MAXADDR = 0x7FFF
	MAXMEM  = 0x6000
	MAXUINT = 0xFFFF
)

// Symbol Table
const USERMEMBEGIN = uint16(16)
const HACKPREFIX = "111"
const MAXUINT15 = 0b0111111111111111

type SymbolMap map[string]uint16

var DEFAULTSYMBOLMAP = SymbolMap{
	"R0":     0,
	"R1":     1,
	"R2":     2,
	"R3":     3,
	"R4":     4,
	"R5":     5,
	"R6":     6,
	"R7":     7,
	"R8":     8,
	"R9":     9,
	"R10":    10,
	"R11":    11,
	"R12":    12,
	"R13":    13,
	"R14":    14,
	"R15":    15,
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"SCREEN": 16384,
	"KBD":    24576,
}

type SymbolTable struct {
	stmap SymbolMap
	stcnt uint16
}

var DEFAULTSYMTABLE = SymbolTable{
	DEFAULTSYMBOLMAP,
	USERMEMBEGIN,
}

func (s SymbolTable) Clone() (newtable SymbolTable) {
	newmap := make(SymbolMap)
	maps.Copy(newmap, s.stmap)
	newcnt := s.stcnt
	newtable = SymbolTable{
		newmap,
		newcnt,
	}

	return
}

func (s SymbolTable) Lookup(sym string) (value uint16, found bool) {
	value, found = s.stmap[sym]

	return
}

func (s *SymbolTable) AddLabel(sym string, addr uint16) (err error) {
	if a, found := s.Lookup(sym); found {
		err = fmt.Errorf("label already exists: %s --> %d", sym, a)
	} else {
		(*s).stmap[sym] = addr
	}

	return
}

func (s *SymbolTable) AddSymbol(sym string) (addr uint16, err error) {
	var found bool

	if err = CheckSymbolName(sym); err == nil {
		if addr, found = s.Lookup(sym); !found {
			addr = s.stcnt
			(*s).stmap[sym] = addr
			(*s).stcnt++
		}
	}

	return
}

func (s SymbolTable) PrintSymboltable(batch int) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	keys := make([]string, 0, len(s.stmap))
	for k := range s.stmap {
		keys = append(keys, k)
	}
	// sort.Strings(keys)
	sort.Slice(keys, func(i, j int) bool { return s.stmap[keys[i]] < s.stmap[keys[j]] })

	cnt := 1
	for _, key := range keys {
		fmt.Fprintf(w, "%8s : %-6d\t", key, s.stmap[key])
		if cnt == batch {
			fmt.Fprintln(w, "")
			cnt = 1
		} else {
			cnt++
		}
	}
	fmt.Fprintln(w, "")
	w.Flush()
}

type InstructionMap map[string]string

func (im InstructionMap) Lookup(instr string) (bitmask string, err error) {
	// TODO: SORT REGISTER!
	bitmask, found := im[instr]
	if !found {
		err = fmt.Errorf("unknown instruction: %s", instr)
	}

	return
}

// ALU destinations
var ALUDEST = InstructionMap{
	"":    "000",
	"M":   "001",
	"D":   "010",
	"DM":  "011",
	"A":   "100",
	"AM":  "101",
	"AD":  "110",
	"ADM": "111",
}

// ALU operations
var ALUCOMP = InstructionMap{
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
	"BRK": "1111111",
}

// ALU Jumps
var ALUJUMP = InstructionMap{
	"":    "000",
	"JGT": "001",
	"JEQ": "010",
	"JGE": "011",
	"JLT": "100",
	"JNE": "101",
	"JLE": "110",
	"JMP": "111",
}
