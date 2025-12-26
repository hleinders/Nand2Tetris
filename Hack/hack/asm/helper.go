/*
Copyright © 2025 harald@leinders.de
*/

package asm

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Hack language helper functions
func StringToInt(s string) (int64, error) {
	// base = 0 → automatische Erkennung der Basis
	// bitSize = 64 → Ergebnis passt in int64
	return strconv.ParseInt(s, 0, 64)
}

func Bool2Bit(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func CheckSize15(value int64) (rval uint16, err error) {
	if value >= 0 && value <= MAXUINT15 {
		rval = uint16(value)
	} else {
		err = fmt.Errorf("value too large: %d > %d", value, MAXUINT15)
	}

	return
}

func CheckSymbolName(sym string) (err error) {
	var hackSymbolRegex = regexp.MustCompile(`^[A-Za-z_.$:][A-Za-z0-9_.$:]*$`)

	if sym == "" {
		err = fmt.Errorf("symbol name must not be empty")
	}

	// Regex-Regeln prüfen
	if !hackSymbolRegex.MatchString(sym) {
		err = fmt.Errorf("invalid symbol name %s: must start with a letter or [_.$:] and contain only letters, digits and [_.$:]", sym)
	}

	return
}

func SymTab2Slice(symtab SymbolTable, sortByVal bool) (symslice []string) {
	m := symtab.stmap

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	if sortByVal {
		sort.Slice(keys, func(i, j int) bool { return m[keys[i]] < m[keys[j]] })
	} else {
		sort.Strings(keys)
	}

	for _, k := range keys {
		symslice = append(symslice, fmt.Sprintf("%s: %d", k, m[k]))
	}

	return
}

func GetFileExtension(path string) string {
	return filepath.Ext(path)
}

func NotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func ChangeFileExtension(path string, newExt string) (string, error) {
	info, err := os.Stat(path)
	// if err != nil {
	// 	return "", err // Datei/Verzeichnis existiert nicht oder kein Zugriff
	// }

	if err == nil && info.IsDir() {
		return "", fmt.Errorf("angegebener Pfad ist ein Verzeichnis")
	}

	// Extension normalisieren
	if !strings.HasPrefix(newExt, ".") {
		newExt = "." + newExt
	}

	dir := filepath.Dir(path)
	base := filepath.Base(path)
	oldExt := filepath.Ext(base)

	nameWithoutExt := strings.TrimSuffix(base, oldExt)
	newName := nameWithoutExt + newExt

	return filepath.Join(dir, newName), nil
}

func Hack2hex(input []string) ([]string, error) {
	var out []string
	var err error

	for _, val := range input {
		hval, err := strconv.ParseUint("0b"+val, 0, 64)
		if err != nil {
			err = fmt.Errorf("convert to hex failed: %s", err)
		}
		out = append(out, fmt.Sprintf("%04X", hval))
	}

	if len(input) != len(out) {
		err = fmt.Errorf("something went wrong: convert to hex failed")
	}

	return out, err
}

func ReverseInstructionMap(input InstructionMap) (out InstructionMap) {
	out = make(InstructionMap, len(input))

	for k, v := range input {
		out[v] = k
	}

	return
}

// returns program as hex dump
func LoadProgram(fname string, makeHex bool) (hackProgram []string, err error) {
	ext := GetFileExtension(fname)

	switch ext {
	case ".asm", ".S":
		hackProgram, err = OnePassAssembler(fname)
	case ".hex", ".bin":
		hackProgram, err = ReadObjFile(fname)
	case ".hack":
		hackProgram, err = ReadHackFile(fname)
	default:
		err = fmt.Errorf("unknown file extension: %s", ext)
	}

	if err != nil {
		return
	}

	if makeHex {
		hackProgram, err = Hack2hex(hackProgram)
	}
	return
}
