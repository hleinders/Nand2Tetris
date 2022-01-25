package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	flag "github.com/spf13/pflag"
	"golang.org/x/term"

	"github.com/fatih/color"
)

const (
	appVersion = "1.0 (2022-01-20)"
	appAuthor  = "Harald Leinders"
	appEMail   = "harald@leinders.de"
)

// color funcs
var bold = color.New(color.Bold).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

// Return values
const (
	OK = iota
	ErrUndef
	ErrConfFile
	ErrReadSource
	ErrCleanSource
	ErrParseSource
	ErrTranslateSource
	ErrWriteFile
)

// FlagType is an Object containing all needed flags
type FlagType struct {
	help, debug, verbose, version, noColor bool
	configFile, createConfig               string
}

func isTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func ttyHasColor() bool {
	return os.Getenv("TERM") != "dumb" && isTTY()
}

func init() {
	// Use available Cores for Goroutines
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Detect terminal
	if ttyHasColor() {
		color.NoColor = false
	} else {
		color.NoColor = true
	}

	// Detect locale for printing:
	// if runtime.GOOS == "windows" {
	// 	asciiChars()
	// } else if runtime.GOOS == "darwin" {
	// 	sbFullBlock = "█"
	// 	sbShadedBlock = "▒"
	// 	sbLeftChar = "┫"
	// 	sbRightChar = "┣"
	// }
}

func usage() {
	fmt.Fprintf(os.Stderr, bold(yellow("\nUsage:    %s [options] asm-file\n")), filepath.Base(os.Args[0]))
	fmt.Fprintln(os.Stderr, bold("\nOptions:"))
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
}

func version() {
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "Version:     %s\n", appVersion)
	fmt.Fprintf(os.Stderr, "Go version:  %s\n", runtime.Version())
	fmt.Fprintf(os.Stderr, "Go compiler: %s\n", runtime.Compiler)
	fmt.Fprintf(os.Stderr, "Binary type: %s (%s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(os.Stderr, "Author:      %s (%s)\n", appAuthor, appEMail)
}

func check(e error, rcode int) {
	if e != nil {
		fmt.Fprintf(os.Stderr, red("*** Error: %+v\n"), e)
		os.Exit(rcode)
	}
}
