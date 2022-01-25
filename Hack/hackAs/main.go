package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hleinders/MyLog"
	"github.com/hleinders/hackas/parser"
	flag "github.com/spf13/pflag"
	viper "github.com/spf13/viper"
)

func main() {
	var flags FlagType
	var lg MyLog.Log

	// Initial setup:
	viper.GetViper()
	viper.SetConfigName(".hack/hackas.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	flag.Usage = usage

	// Set defaults:
	viper.SetDefault("Verbose", false)
	viper.SetDefault("NoColor", false)

	// Check args
	// Bools
	flag.BoolVarP(&flags.help, "help", "h", false, "show this help")
	flag.BoolVarP(&flags.debug, "debug", "d", false, "debug mode")
	flag.BoolVarP(&flags.verbose, "verbose", "v", false, "verbose mode")
	flag.BoolVarP(&flags.version, "version", "V", false, "show version info")
	flag.BoolVarP(&flags.noColor, "mono", "m", false, "do not use colors (monochrom mode)")

	// Parameter
	flag.StringVarP(&flags.configFile, "conf", "c", "", "config `file` to use (default: ($HOME|.)/.hack/hackas.yml")
	flag.StringVarP(&flags.createConfig, "create-conf", "C", "", "write config template to `file`")

	flag.CommandLine.MarkHidden("debug")

	flag.Parse()

	// bind only a subset of flags:
	// viper.BindPFlags(flag.CommandLine)
	viper.BindPFlag("Verbose", flag.CommandLine.Lookup("verbose"))
	viper.BindPFlag("NoColor", flag.CommandLine.Lookup("mono"))

	// read defaults from config file:
	if flags.configFile != "" {
		viper.SetConfigFile(flags.configFile)
		err := viper.ReadInConfig() // Find and read the config file
		check(err, ErrConfFile)
	} else {
		// don't panic if default config not found
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				check(err, ErrConfFile)
			}
		}
	}

	// == Init Logger with Defaults:
	lg.Init(os.Stdout, os.Stderr)

	if flags.help || len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(OK)
	}

	if flags.version {
		version()
		os.Exit(0)
	}

	if flags.createConfig != "" {
		viper.Set("Verbose", false)
		viper.Set("NoColor", false)
		viper.SafeWriteConfigAs(flags.createConfig)
		os.Exit(0)
	}

	if viper.GetBool("NoColor") {
		color.NoColor = true
	}

	if flags.debug {
		fmt.Printf(red("*** Debug: Settings:\n    %+v\n"), viper.AllSettings())
	}

	lg.SetModeBool(MyLog.LgVerbose, viper.GetBool("Verbose"))
	lg.SetModeBool(MyLog.LgDebug, flags.debug)
	lg.SetModeBool(MyLog.LgColor, !color.NoColor)
	lg.SetInteractive()
	lg.SetNoPrefix()

	lg.VerboseInfo("\n\nGolang HACK Assembler %s", appVersion)

	// All remaining args should be asm files
	for _, inFile := range flag.Args() {
		// parse source, get machine code and symbol table

		p, err := parser.New(&lg, inFile)
		check(err, ErrReadSource)

		mCode, err := p.Parse()
		check(err, ErrParseSource)

		outFile := inFile[0:len(inFile)-len("asm")] + "hack"
		err = os.WriteFile(outFile, []byte(mCode), 0644)
		check(err, ErrWriteFile)

		if viper.GetBool("Verbose") {
			var lbl string

			st := p.GetSymbolTable()
			pSource := p.GetPureSource()

			mc := strings.Split(mCode, "\n")
			ps := strings.Split(pSource, "\n")
			// us := st.GetUserSymbols()
			ad := st.GetUserLabelAddresses()

			if len(mc) != len(ps) {
				fmt.Printf("Error: source and code have different lengths!")
				os.Exit(ErrUndef)
			}

			lg.VerboseInfo("\nAssembly successful!\n")

			fmt.Println("Symbol table:")
			fmt.Println("=============================")
			st.PrettyPrint()

			fmt.Printf("\n\n")
			fmt.Println("Machine Code      Line Nr.    Label              Assenmbler")
			fmt.Println("-----------------------------------------------------------")
			for i, m := range mc {
				if m == "" {
					continue
				}

				ac := ps[i]
				lbl = ""

				if v, ok := ad[uint16(i)]; ok {
					lbl = v
				}

				fmt.Printf("%s    %6d     %-16s  %s\n", m, i, lbl, ac)
			}
			fmt.Printf("\n")
		}
	}

	os.Exit(OK)
}
