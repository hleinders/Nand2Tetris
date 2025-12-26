/*
Copyright © 2025 harald@leinders.de
*/

package tools

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var (
	Yellow = color.New(color.FgYellow).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Green  = color.New(color.FgGreen).SprintFunc()
	Bold   = color.New(color.Bold).SprintFunc()
)

func Debug(format string, a ...any) {
	if viper.GetBool("debug") {
		fmt.Printf(Red("*** DEB: ")+format, a...)
	}
}

func Error(format string, a ...any) {
	fmt.Printf(Red("*** ERROR: ")+format+"\n", a...)
}

func Warn(format string, a ...any) {
	fmt.Printf(Yellow("*** WARN: ")+format, a...)
}

func Info(format string, a ...any) {
	fmt.Printf(Green(format), a...)
}

func Verbose(format string, a ...any) {
	if viper.GetBool("verbose") {
		fmt.Printf(format, a...)
	}
}
