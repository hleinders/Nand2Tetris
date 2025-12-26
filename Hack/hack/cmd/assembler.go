/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hleinders/hack/asm"
	"github.com/hleinders/hack/tools"
	"github.com/spf13/cobra"
)

// assemblerCmd represents the assembler command
var assemblerCmd = &cobra.Command{
	Use:     "assembler source.asm [-o outfile.hack]",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"assem", "asm"},
	Short:   "This invokes the hack assembler",
	Long: `The hack assembler translates hack mnemonics into hack objectcode.
The source file should have the extension ".asm", or automatic outfiles
are named "source.ext.hack" instead of "source.hack".
The objectcode can be written in the normal bit format as ".hack" file or
(when using the -H option) as a series of hex values as ".hex" file.
Use the "-o" option to write to another file than "source.hack"`,
	Run: ExecAssembler,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("assembler called")
	// },
}

func init() {
	rootCmd.AddCommand(assemblerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// assemblerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// assemblerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	assemblerCmd.Flags().BoolP("hexdump", "H", false, "create hex object file instead of hack dump")
	assemblerCmd.Flags().BoolP("stripped", "S", false, "save the intermediate, stripped source (.S)")
	assemblerCmd.Flags().BoolP("symbol-table", "T", false, "save the symbol table (.sym)")
	assemblerCmd.Flags().StringP("outfile", "o", "", "use `filename` instead of basename of source")
}

func ExecAssembler(cmd *cobra.Command, args []string) {
	for _, filename := range args {
		err := asm.Assemble(cmd.Flags(), filename)
		if err != nil {
			tools.Error(err.Error())
		}
	}
}
