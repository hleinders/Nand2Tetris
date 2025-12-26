/*
Copyright © 2025 harald@leinders.de
*/

package asm

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hleinders/hack/tools"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Hack language assembler functions
func ReadAsmFile(filename string) (content string, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	content = string(data)
	return
}

func WriteFile(filename string, content []string) (err error) {
	bytearray := strings.Join(content, "\n") + "\n"
	err = os.WriteFile(filename, []byte(bytearray), 0644)

	return
}

func WriteAsmFile(filename string, asmcode []string, indent bool) (err error) {
	var indentCode []string

	if indent {
		for _, line := range asmcode {
			indentCode = append(indentCode, "\t"+line)
		}
	} else {
		indentCode = asmcode
	}
	err = WriteFile(filename, indentCode)

	return
}

func WriteObjFile(filename string, hackcode []string, useHex bool) (err error) {
	if useHex {
		hackcode, err = Hack2hex(hackcode)
		if err != nil {
			return
		}
	}
	err = WriteFile(filename, hackcode)

	return
}

func WriteStrippedFile(filename string, hackcode []string) (err error) {
	err = WriteFile(filename, hackcode)

	return
}

func WriteSymbolTable(filename string, symtab SymbolTable) (err error) {
	sl := SymTab2Slice(symtab, false)
	err = WriteFile(filename, sl)

	return
}

// CleanFile is a function to strip all non language stuff from the source
// 1: remove all comments.
// 2: remove all whitespaces
// 3: remove all remaining empty lines
func CleanFile(content string) (cleanedContent []string, err error) {
	reDelCmts := regexp.MustCompilePOSIX(`//.*\n`)
	reDelWS := regexp.MustCompilePOSIX(`[\t\v\f\r ]+`)
	reDelEmpty := regexp.MustCompilePOSIX(`[\n\r]{2,}`)

	tools.Debug("RAW:      " + strings.ReplaceAll(content, "\n", " ") + "\n")

	tempContent := reDelCmts.ReplaceAllString(content, "\n")
	tools.Debug("DelCmts:  " + strings.ReplaceAll(tempContent, "\n", " ") + "\n")

	tempContent = reDelWS.ReplaceAllString(tempContent, "")
	tempContent = reDelEmpty.ReplaceAllString(tempContent, "\n")
	tools.Debug("DelWS:    " + strings.ReplaceAll(tempContent, "\n", " ") + "\n")

	cleanedContent = strings.Split(strings.TrimSpace(tempContent), "\n")

	return
}

func Resolve(content []string) (resolvedContent []string, symtable SymbolTable, err error) {
	var icnt uint16

	symtable = DEFAULTSYMTABLE.Clone()
	icnt = 0

	reLabl := regexp.MustCompile(`^\((.+)\)`)

	for _, line := range content {
		m := reLabl.FindStringSubmatch(line)
		if m != nil {
			all := m[0]
			sym := m[1]

			tools.Debug("L: %-5s\tvalue: %8s\taddr: %d\n", all, sym, icnt)

			err = symtable.AddLabel(sym, icnt)
			if err != nil {
				return
			}

			tools.Verbose("     %-8s: %d\n", sym, icnt)
		} else {
			resolvedContent = append(resolvedContent, line)
			icnt++
		}
	}

	return
}

func ParseA(instruction string, symtable *SymbolTable) (opcode, resolvedLine string, err error) {
	var num64 int64
	var num16 uint16

	reA := regexp.MustCompile(`^[@](.+)$`)

	a := reA.FindStringSubmatch(instruction)
	if a != nil {

		all := a[0]
		val := a[1]

		if num64, err = StringToInt(val); err == nil {
			// is a number value
			num16, err = CheckSize15(num64)
			if err != nil {
				// but value too large
				err = fmt.Errorf("instruction %s: %s", instruction, err.Error())
				return
			} else {
				// value is good
				opcode = fmt.Sprintf("%016b", num16)
				resolvedLine = fmt.Sprintf("@%d", num16)
			}
		} else {
			// is a variable name
			num16, err = symtable.AddSymbol(val)
			if err != nil {
				return
			} else {
				// variable is good:
				opcode = fmt.Sprintf("%016b", num16)
				resolvedLine = fmt.Sprintf("%s(=%d)", instruction, num16)
			}
		}

		tools.Debug("A: %-5s\tvalue: %s\t\t%s\n", all, val, opcode)
	} else {
		err = fmt.Errorf("error parsing A-instruction %s: %s", instruction, err)
	}

	return
}

func ParseC(instruction string) (opcode, resolvedLine string, err error) {
	var destBits, compBits, jumpBits string
	reC := regexp.MustCompile(`^(?:([^@=;][^()=;]*)=)?([^@;()][^;]*)(?:;(.+))?$`)

	c := reC.FindStringSubmatch(instruction)
	if c != nil {
		all := c[0]
		dest := strings.ToUpper(c[1])
		comp := strings.ToUpper(c[2])
		jump := strings.ToUpper(c[3])

		destBits, err = ALUDEST.Lookup(dest)
		if err != nil {
			return
		}
		compBits, err = ALUCOMP.Lookup(comp)
		if err != nil {
			return
		}
		jumpBits, err = ALUJUMP.Lookup(jump)
		if err != nil {
			return
		}

		opcode = HACKPREFIX + compBits + destBits + jumpBits
		resolvedLine = strings.ToUpper(instruction)

		tools.Debug("C: %-5s\tdest: %-5s\tcomp: %-5s\tjump: %s\t\t%s\n", all, dest, comp, jump, opcode)
	} else {
		err = fmt.Errorf("error parsing C-instruction %s: %s", instruction, err)
	}

	return
}

func ParseLine(instruction string, symtable *SymbolTable) (opcode, resolvedLine string, err error) {
	if strings.HasPrefix(instruction, "@") {
		opcode, resolvedLine, err = ParseA(instruction, symtable)
	} else {
		opcode, resolvedLine, err = ParseC(instruction)
	}

	return
}

func Parse(content []string, symtable *SymbolTable) (hackProgram []string, err error) {
	// Input for Parse must already has been cleaned and resolved!
	// It may only contain A or C instructions
	// A-Instruction: @value or @variable
	// C-Instruction: dest=comp;jump
	var opcode, resolvedLine string

	for lcounter, instr := range content {
		opcode, resolvedLine, err = ParseLine(instr, symtable)

		if err != nil {
			tools.Error("Error parsing instruction: %s (line %d). Abort!", instr, lcounter)
			break
		}

		hexval, _ := strconv.ParseUint("0b"+opcode, 0, 64)

		tools.Verbose("%6d:    %-12s    %s     %04X\n", lcounter, resolvedLine, opcode, hexval)

		hackProgram = append(hackProgram, opcode)
	}

	return
}

func OnePassAssembler(fname string) (hackProgram []string, err error) {
	raw, _ := ReadAsmFile(fname)
	clean, err := CleanFile(raw)
	if err != nil {
		return
	}
	resolved, symtable, err := Resolve(clean)
	if err != nil {
		return
	}
	hackProgram, err = Parse(resolved, &symtable)

	return
}

func Assemble(flags *pflag.FlagSet, filename string) (err error) {
	raw, _ := ReadAsmFile(filename)

	tools.Info("Assembling %s\n", filename)

	tools.Verbose("%s %s\n", tools.Bold(" Pass 0:"), "preparing input file")
	clean, err := CleanFile(raw)
	if err != nil {
		return
	}

	tools.Verbose("\n%s %s\n", tools.Bold(" Pass 1:"), "resolving labels")
	resolved, symtable, err := Resolve(clean)
	if err != nil {
		return
	}

	tools.Verbose("\n%s %s\n", tools.Bold(" Pass 2:"), "parsing source")
	bytecode, err := Parse(resolved, &symtable)
	if err != nil {
		return
	}

	if viper.GetBool("verbose") {
		fmt.Println(tools.Bold("\n" + " Symbol Table:"))
		symtable.PrintSymboltable(4)
	}

	outfile, _ := flags.GetString("outfile")
	hexdump, _ := flags.GetBool("hexdump")
	newext := ".hack"
	if outfile == "" {
		outfile = filename
	}
	if hexdump {
		newext = ".hex"
	}

	outfile, err = ChangeFileExtension(outfile, newext)
	if err != nil {
		return
	}

	tools.Info("Writing object file %s\n", outfile)
	err = WriteObjFile(outfile, bytecode, hexdump)
	if err != nil {
		tools.Error(err.Error())
	}

	if set, _ := flags.GetBool("stripped"); set {
		outfile, err = ChangeFileExtension(outfile, ".S")
		if err != nil {
			return
		}

		tools.Info("Writing stripped source to %s\n", outfile)
		err = WriteStrippedFile(outfile, resolved)
	}

	if set, _ := flags.GetBool("symbol-table"); set {
		outfile, err = ChangeFileExtension(outfile, ".sym")
		if err != nil {
			return
		}

		tools.Info("Writing ssymbol table to %s\n", outfile)
		err = WriteSymbolTable(outfile, symtable)
	}

	return
}
