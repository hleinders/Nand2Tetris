/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/hleinders/hack/asm"
	"github.com/hleinders/hack/emu"
	"github.com/hleinders/hack/tools"
	"github.com/spf13/cobra"
)

var clock float32

// simulatorCmd represents the simulator command
var simulatorCmd = &cobra.Command{
	Use:     "simulator",
	Aliases: []string{"simul", "sim"},
	Short:   "watch a hack CPU executing a program (asm,hack,hex)",
	Long: `Look at the hack cpu while a program is beeing executed. All registers,
the cpu status flags an the relevant portions of ROM (the program)
and RAM are shown.
You can execute the program with single steps or a given clock rate.
`,
	Run: ExecSimulator,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("simulator called")
	// },
}

func init() {
	rootCmd.AddCommand(simulatorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// simulatorCmd.PersistentFlags().String("foo", "", "A help for foo")
	simulatorCmd.PersistentFlags().Float32VarP(&clock, "clock", "c", 1.0, "initial clock rate")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// simulatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ExecSimulator(cmd *cobra.Command, args []string) {
	for _, filename := range args {
		program, err := asm.LoadProgram(filename, true)
		if err != nil {
			fmt.Println(err)
		}
		err = emu.CreateSimulator(program, clock)
		if err != nil {
			tools.Error("%s\n", err)
		}
	}
}
