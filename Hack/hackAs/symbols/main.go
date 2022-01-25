package symbol

import (
	"fmt"
	"sort"
	"strings"
)

// SymbolTable is an interface to which modular functionality is provided.
type SymbolTable struct {
	symbol map[string]uint16
}

const (
	symbolBase = uint16(0x0010)
	symbolMax  = uint16(0x3fff)
)

var symbolCurrent = symbolBase

var predefined map[string]uint16 = map[string]uint16{
	"R0": 0x0000, "R1": 0x0001, "R2": 0x0002, "R3": 0x0003, "R4": 0x0004, "R5": 0x0005, "R6": 0x0006, "R7": 0x0007,
	"R8": 0x0008, "R9": 0x0009, "R10": 0x000a, "R11": 0x000b, "R12": 0x000c, "R13": 0x000d, "R14": 0x000e, "R15": 0x000f,
	"SP": 0x0000, "LCL": 0x0001, "ARG": 0x0002, "THIS": 0x0003, "THAT": 0x0004, "SCREEN": 0x4000, "KBD": 0x6000,
}

// New returns the api interface.
func New() *SymbolTable {
	nt := make(map[string]uint16)
	for k, v := range predefined {
		nt[k] = v
	}

	st := SymbolTable{
		symbol: nt,
	}

	return &st
}

func (st *SymbolTable) contains(sbl string) bool {
	return st.symbol[sbl] != 0
}

// AddEntry is a function to add an entry to the symbol table
func (st *SymbolTable) AddEntry(sbl string, address uint16) {
	st.symbol[sbl] = address
}

// If sbl was a Label, it has been set in pass1 via
// AddEntry with an Address pointing to ROM starting at 0,
// otherwise it's a variable with a RAM address starting at 16
func (st *SymbolTable) GetAddress(sbl string) uint16 {
	if !st.contains(sbl) {
		st.AddEntry(sbl, symbolCurrent)
		symbolCurrent++
		if symbolCurrent >= symbolMax {
			panic("ERR: Symbol Table overflow! Abort!")
		}
	}

	return st.symbol[sbl]
}

func (st *SymbolTable) GetUserSymbols() map[string]uint16 {
	us := make(map[string]uint16)
	for k := range st.symbol {
		if _, ok := predefined[k]; !ok {
			us[k] = st.symbol[k]
		}
	}
	return us
}

func (st *SymbolTable) GetUserLabelAddresses() map[uint16]string {
	ad := make(map[uint16]string)
	for k, v := range st.GetUserSymbols() {
		if strings.ToUpper(k) == k {
			ad[v] = k
		}
	}

	return ad
}

func (st SymbolTable) PrettyPrint() {
	t := st.GetUserSymbols()
	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%-16s  --->  %d\n", k, t[k])
	}
}
