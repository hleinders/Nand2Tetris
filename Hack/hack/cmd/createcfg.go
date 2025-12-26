/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCfgFile string

// createcfgCmd represents the createcfg command
var createcfgCmd = &cobra.Command{
	Use:   "createcfg",
	Short: "Creates a default hack configuration file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("createcfg called")
	// },
	Run: ExecCreatecfg,
}

func init() {
	rootCmd.AddCommand(createcfgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createcfgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createcfgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createcfgCmd.Flags().StringVar(&newCfgFile, "create-config", "", "create `config file`")
}

func ExecCreatecfg(cmd *cobra.Command, args []string) {
	fmt.Println("createcfg called")
	if newCfgFile != "" {
		viper.SafeWriteConfigAs(newCfgFile)
	} else {
		fmt.Println(viper.AllSettings())
	}
}
