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

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute a hack program",
	Long: `With 'hack execute' you can run a hack program, either in single steps
or continously with a given clock speed. At the end, the RAM content
is dumped with the given ranges.`,
	Run: ExecExecute,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("execute called")
	// },
}

func init() {
	rootCmd.AddCommand(executeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// executeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// executeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ExecExecute(cmd *cobra.Command, args []string) {
	fmt.Println("execute called")
	var hackProgram []string
	var err error

	MAXSTEPS := 10000

	for _, filename := range args {

		tools.Info("Running %s\n", filename)

		tools.Verbose("%s %s\n", tools.Bold(" Pass 0:"), "reading input file")

		if ext := asm.GetFileExtension(filename); ext == ".hex" {
			hackProgram, err = asm.ReadObjFile(filename)
		} else {
			hackProgram, err = asm.ReadHackFile(filename)
		}
		if err != nil {
			tools.Error(err.Error())
		}

		tools.Verbose("%s %s\n", tools.Bold(" Pass 1:"), "initialize emulator")

		hackCPU := emu.HACKCPU{}
		hackCPU.Init()
		err = hackCPU.LOAD(hackProgram)
		if err != nil {
			fmt.Println(err)
		}

		hackCPU.PrintROM(0, 64)
		hackCPU.RUNNING = true

		for i := 0; i < MAXSTEPS; i++ {
			hackCPU.Step()
			// hackCPU.PrintStatus()
			// hackCPU.PrintRAM(0, 32)
			if hackCPU.RUNNING == false {
				fmt.Printf("Done after %d steps\n", i)
				break
			}
		}
		hackCPU.PrintStatus()
		hackCPU.PrintRAM(0, 48)
	}
}
