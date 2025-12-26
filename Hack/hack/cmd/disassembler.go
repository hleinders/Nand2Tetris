/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hleinders/hack/asm"
	"github.com/hleinders/hack/tools"
	"github.com/spf13/cobra"
)

// disassemblerCmd represents the disassembler command
var disassemblerCmd = &cobra.Command{
	Use:     "disassembler",
	Aliases: []string{"disasm", "disassem"},
	Short:   "Disassemble an existing hack programm (.hex, .hack)",
	Long: `with 'hack disassembler' you may recreate a hack assembler program
with mnemonics from an existing binary hack program. Due to the fact,
that a hack program doesn't include a symbol table, there is no way
to restore labels or variable names, so this has to be done manually.`,
	Run: ExecDisAssembler,
}

func init() {
	rootCmd.AddCommand(disassemblerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// disassemblerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// disassemblerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	disassemblerCmd.Flags().StringP("outfile", "o", "", "write disassembled code to `filename`")
	disassemblerCmd.Flags().BoolP("force", "f", false, "force overwrite file")
}

func ExecDisAssembler(cmd *cobra.Command, args []string) {
	for _, filename := range args {
		err := asm.DisAssemble(cmd.Flags(), filename)
		if err != nil {
			tools.Error(err.Error())
		}
	}
}
