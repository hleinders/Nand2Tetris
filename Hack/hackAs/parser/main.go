package parser

// this package provides the Parser bject and it's API functions.
// Some of the reoutines are a bit ugly with golang, but I tried
// to keep it close to the implementation instructions from
// "The Elements of Computing Systems" (Nand2Tetris)

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hleinders/MyLog"
	code "github.com/hleinders/hackas/code"
	symbols "github.com/hleinders/hackas/symbols"
)

const (
	UNKNOWN = iota
	A_INSTRUCTION
	C_INSTRUCTION
	L_INSTRUCTION
)

const (
	LabelFmt = `\(([A-Za-z_.$:][A-Za-z0-9_.$:]*)\)`
	SymblFmt = `@([A-Za-z0-9_.$:]+)`
	InstrFmt = `(?P<dest>[A-Z]*)=?(?P<comp>[A-Za-z0-9=\-+&|!]+);?(?P<jump>[A-Z]*)` // --> (dest)(comp)(jmp)
)

type Parser struct {
	rawSource   string // Source as read from inpout
	cleanSource string // Source with all comments and whitespace removed
	pureSource  string // Source with all labels resolved and removed
	symTable    *symbols.SymbolTable
	reader      []string
	logger      *MyLog.Log
}

// New returns the api interface.
func New(lg *MyLog.Log, fname string) (*Parser, error) {
	lg.Verbose("\nReading Source: %s", fname)
	raw, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("could not read source file '%s'", fname)
	}

	// Yes, I know... it's a litte bit off the implementation instructions,
	// but so much easier later on
	lg.Verbose("Cleaning...")
	clean := cleanSrc(raw)

	return &Parser{rawSource: string(raw), cleanSource: clean, logger: lg}, nil
}

func (p *Parser) GetSymbolTable() *symbols.SymbolTable {
	return p.symTable
}

func (p *Parser) GetPureSource() string {
	return p.pureSource
}

// "hasMoreLines" checks if more lines are available
func (p *Parser) hasMoreLines() bool {
	return len(p.reader) > 0
}

// "advance" gets next input line
func (p *Parser) advance() string {

	next := p.reader[0]
	if len(p.reader) > 1 {
		p.reader = p.reader[1:]
	} else {
		p.reader = nil
	}

	return next
}

// "instructionType" detects the instruction type
func (p *Parser) instructionType(mnem string) int {
	if mnem == "" {
		return UNKNOWN
	}

	fChar := strings.TrimSpace(mnem)[0]

	switch fChar {
	case '@':
		return A_INSTRUCTION
	case '(':
		return L_INSTRUCTION
	default:
		return C_INSTRUCTION
	}
}

func (p *Parser) symbol(mnem string) string {
	var match []string
	reSmbl := regexp.MustCompile(SymblFmt)

	if reSmbl.Match([]byte(mnem)) {
		match = reSmbl.FindStringSubmatch(mnem)
		if len(match) > 1 {
			return match[1]
		}
	}
	panic(fmt.Sprintf("error: Empty Symbol %+v found", match))
}

func (p *Parser) splitC(mnem string) (string, string, string, error) {
	reC := regexp.MustCompile(InstrFmt)
	match := reC.FindStringSubmatch(mnem)

	if len(match) == 4 {
		return match[1], match[2], match[3], nil

	}

	return "", "", "", fmt.Errorf("parser error: could not parse c instruction '%s'", mnem)
}

func (p *Parser) dest(mnem string) (string, error) {
	destStr, _, _, err := p.splitC(mnem)

	return destStr, err
}

func (p *Parser) comp(mnem string) (string, error) {
	_, compStr, _, err := p.splitC(mnem)

	return compStr, err
}

func (p *Parser) jump(mnem string) (string, error) {
	_, _, jumpStr, err := p.splitC(mnem)

	return jumpStr, err
}

// main parser function
func (p *Parser) Parse() (string, error) {
	var mlang []string
	var opcode, dest, comp, jump string
	var val int
	var err error

	// Pass 1: resolve labels
	// creates the parse objects oureSrc and synbolTable
	p.logger.Verbose("Pass 1: resolving labels")
	err = p.resolve()
	if err != nil {
		return "", err
	}

	// normally, I would iterate over the source lines an do the job.
	// But, I just obey the instructions and make (and use) some, for golang
	// quite unneccessary, api funnctions. That's why "resolve" creates the reader

	cd := code.New()

	// Pass 2: translate the source
	p.logger.Verbose("Pass 2: parsing & compiling")
	for p.hasMoreLines() {
		mnemonic := p.advance()
		it := p.instructionType(mnemonic)
		opcode = ""
		val = 0
		switch it {
		case L_INSTRUCTION:
			continue // already resolved in pass 1
		case A_INSTRUCTION:
			sm := p.symbol(mnemonic)
			if val, err = strconv.Atoi(sm); err == nil {
				// is a number, convert to "binary" string:
				opcode = fmt.Sprintf("%016b", uint16(val))
			} else {
				val = int(p.symTable.GetAddress(sm))
				opcode = fmt.Sprintf("%016b", uint16(val))
			}
		case C_INSTRUCTION:
			// normally, it would just be:
			// destStr, compStr, jumpStr, err := p.splitC(mnemonic)
			// but again, we follow the implementation instructions and
			// complete the API :-)

			destStr, err := p.dest(mnemonic)
			if err != nil {
				return "", err
			}
			compStr, err := p.comp(mnemonic)
			if err != nil {
				return "", err
			}
			jumpStr, err := p.jump(mnemonic)
			if err != nil {
				return "", err
			}

			dest = cd.Dest(destStr)
			comp, err = cd.Comp(compStr)
			if err != nil {
				return "", err
			}
			jump = cd.Jump(jumpStr)

			opcode = fmt.Sprintf("111%s%s%s", comp, dest, jump)
		}
		mlang = append(mlang, opcode)
	}

	return strings.Join(mlang, "\n"), nil
}

// helper functions
func cleanSrc(src []byte) string {
	buf := []byte(src)

	// Pass 1: Remove comments
	buf = regexp.MustCompile(`(?m://.*)`).ReplaceAll(buf, nil)

	// Pass 2: Remove empty
	buf = regexp.MustCompile(`(?m:^(?:\s*(?:\r?\n|\r))+)`).ReplaceAll(buf, nil)

	// Pass 2: Remove whitespace
	buf = regexp.MustCompile(`(?m:[\t ]*)`).ReplaceAll(buf, nil)

	return string(buf)
}

func (p *Parser) resolve() error {
	var addrCntr uint16
	var result []string

	symTbl := symbols.New()

	// Init address counter
	addrCntr = 0

	// find all labels
	// Must begin with a char an may contain chars, letters and [_.$:]
	reLbl := regexp.MustCompile(LabelFmt)

	for _, line := range strings.Split(p.cleanSource, "\n") {
		if reLbl.Match([]byte(line)) {
			match := reLbl.FindStringSubmatch(line)
			if len(match) > 1 {
				lbl := match[1]
				symTbl.AddEntry(lbl, addrCntr)
			} else {
				return fmt.Errorf("error: Empty Label %+v found", match)
			}
		} else {
			result = append(result, line)
			addrCntr++
		}
	}

	p.pureSource = strings.Join(result, "\n")
	p.symTable = symTbl
	p.reader = strings.Split(p.pureSource, "\n")

	return nil
}
