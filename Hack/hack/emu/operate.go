package emu

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/hleinders/hack/asm"
)

func AaddrIsValid(val int16) bool {
	return (val >= 0) && (val <= asm.MAXADDR)
}

func MaddrIsValid(val int16) bool {
	return (val >= 0) && (val <= asm.MAXMEM)
}

func SplitC(instr uint16) (dest, comp, jump string, err error) {
	cstr := fmt.Sprintf("%016b", instr)
	if len(cstr) < 16 {
		err = fmt.Errorf("problem parsing C instruction")
		return
	}
	comp, dest, jump = cstr[3:10], cstr[10:13], cstr[13:16]

	return
}

func waitForInput() {
	bufio.NewReader(os.Stdin).ReadBytes(' ')
}

func ClockCycle(clkmode int) {
	// special modes:
	switch clkmode {
	case CLKUNLIM:
		return
	case CLKSINGLE:
		waitForInput()
		return
	}

	// all other are time controlled:
	time.Sleep(time.Duration(clkmode * int(time.Millisecond)))
}
