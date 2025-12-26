/*
Copyright © 2025 harald@leinders.de
*/

package asm

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hleinders/hack/tools"
	"github.com/spf13/pflag"
)

// init reverse lookup tables

var (
	ALUDESTREV InstructionMap = ReverseInstructionMap(ALUDEST)
	ALUCOMPREV InstructionMap = ReverseInstructionMap(ALUCOMP)
	ALUJUMPREV InstructionMap = ReverseInstructionMap(ALUJUMP)
)

func ReadHackFile(filename string) (content []string, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	content = strings.Split(string(data), "\n")

	return
}

func ReadObjFile(filename string) (content []string, err error) {
	var hack string
	var val int64

	hexfile, err := ReadHackFile(filename)
	if err != nil {
		return
	}

	for lcounter, line := range hexfile {
		val, err = strconv.ParseInt(line, 16, 16)
		if err != nil {
			err = fmt.Errorf("%s - line #%d", err.Error(), lcounter)
			return
		} else {
			hack = fmt.Sprintf("@%16b", val)
			content = append(content, hack)
		}
	}

	return
}

func DisassembleHex(hexval uint16) (lineOut string, err error) {

	hack := fmt.Sprintf("%016b", hexval)
	lineOut, err = DisassembleHack(hack)

	return
}

func DisassembleHack(lineIn string) (lineOut string, err error) {
	var val int64
	var comp, dest, jump string

	// Linein MUST be a bitmap as string of 16 char with
	if len(lineIn) != 16 {
		err = fmt.Errorf("line length mismatch: %s", lineIn)
		return
	}

	if strings.HasPrefix(lineIn, "0") { // A-Instruction
		val, err = strconv.ParseInt(lineIn, 2, 16)
		if err != nil {
			return
		} else {
			lineOut = fmt.Sprintf("@%d", val)
		}
		tools.Debug("%s: @%s\n", lineIn, lineIn[1:16])
	} else {
		// C-Instruction:
		// 111    1 111111 111  111
		// Prefix A Comp   Dest Jmp
		// comp, dest, jump := lineIn[3:10], lineIn[10:13], lineIn[13:16]
		tools.Debug("%s: (%s) c: %s d: %s j: %s\n", lineIn, lineIn[0:3], lineIn[3:10], lineIn[10:13], lineIn[13:16])
		comp, err = ALUCOMPREV.Lookup(lineIn[3:10])
		if err != nil {
			err = fmt.Errorf("COMP: %s (%s): %s", lineIn[3:10], lineIn, err.Error())
			return
		}

		dest, err = ALUDESTREV.Lookup(lineIn[10:13])
		if err != nil {
			err = fmt.Errorf("DEST: %s - line #%s: %s", lineIn[10:13], lineIn, err.Error())
			return
		}

		jump, err = ALUJUMPREV.Lookup(lineIn[13:16])
		if err != nil {
			err = fmt.Errorf("JUMP: %s - line #%s: %s", lineIn[12:16], lineIn, err.Error())
			return
		}

		if dest != "" {
			dest = dest + "="
		}
		if jump != "" {
			jump = ";" + jump
		}

		lineOut = dest + comp + jump
	}

	return
}

func DisassembleProg(sourceIn []string) (sourceOut []string, err error) {
	var mnemonic string

	for lcounter, line := range sourceIn {
		if line == "" {
			continue
		}

		mnemonic, err = DisassembleHack(line)
		if err != nil {
			err = fmt.Errorf("%s - line #%d", err.Error(), lcounter)
		}

		sourceOut = append(sourceOut, mnemonic)
	}

	return
}

func DisAssemble(flags *pflag.FlagSet, filename string) (err error) {
	var objectCode, asmCode []string

	tools.Info("DisAssemble %s\n", filename)
	tools.Verbose("%s %s\n", tools.Bold(" Pass 0:"), "reading input file")

	if ext := GetFileExtension(filename); ext == ".hex" {
		objectCode, err = ReadObjFile(filename)
	} else {
		objectCode, err = ReadHackFile(filename)
	}
	if err != nil {
		return
	}

	tools.Verbose("%s %s\n", tools.Bold(" Pass 1:"), "disassembling source")
	asmCode, err = DisassembleProg(objectCode)

	outfile, _ := flags.GetString("outfile")
	force, _ := flags.GetBool("force")

	if outfile != "" {
		if NotExists(outfile) || force {
			tools.Info("Writing assembler file %s\n", outfile)
			WriteAsmFile(outfile, asmCode, true)
		} else {
			err = fmt.Errorf("file %s already exists", outfile)
		}
	} else {
		fmt.Printf("%s\n", tools.Bold(" Result:"))
		for lcount, line := range asmCode {
			fmt.Printf("  %5d:\t%s\n", lcount, line)
		}
	}

	return
}
