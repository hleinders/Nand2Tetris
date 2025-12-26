package emu

import (
	"fmt"
	"strconv"

	"github.com/hleinders/hack/asm"
	"github.com/hleinders/hack/tools"
)

// Memory
var MEM []uint16

const (
	// all values in millisecs
	CLKSINGLE = -1
	CLKUNLIM  = 0
	CLK01HZ   = 100
	CLK05HZ   = 500
	CLK1HZ    = 1000
	CLK5HZ    = 5000
	CLK10HZ   = 10000
	CLK50HZ   = 50000
)

const (
	ZFLAG = 1 << iota
	NFLAG
	OFLAG
	CFLAG
	WFLAG
)

type Statusregister struct {
	register uint8
}

func (sr *Statusregister) Get(flag uint8) bool {
	return (sr.register & flag) != 0
}

func (sr *Statusregister) Set(flag uint8) {
	sr.register |= flag
}

func (sr *Statusregister) Clear(flag uint8) {
	sr.register &= ^flag
}

// Registers
var A, D uint

type HACKCPU struct {
	ALU       int16
	A         int16
	D         int16
	M         int16
	SR        Statusregister
	PC        uint16
	RAM       []int16
	ROM       []uint16
	CLOCKMODE int
	RUNNING   bool
	TRACE     bool
	TRMEMLO   uint16
	TRMEMHI   uint16
	CURINST   uint16
	TRLASTPC  uint16
}

func (cpu *HACKCPU) Init() {
	cpu.RAM = make([]int16, asm.MAXMEM)
	cpu.ROM = make([]uint16, asm.MAXADDR)
	cpu.ALU = 0
	cpu.A = 0
	cpu.D = 0
	cpu.SR = Statusregister{0}
	cpu.PC = 0
	cpu.CLOCKMODE = CLKUNLIM
	cpu.RUNNING = false
	cpu.TRACE = false
	cpu.TRMEMLO = 0
	cpu.TRMEMHI = 32
}

func (cpu *HACKCPU) Reset() {
	cpu.PC = 0
}

func (cpu *HACKCPU) LOAD(hackProgram []string) (err error) {
	var hexval int64
	var cnt int
	var instr string

	for cnt, instr = range hackProgram {
		// skip empty
		if len(instr) == 0 {
			continue
		}

		hexval, err = strconv.ParseInt(instr, 2, 32)
		if err != nil {
			return
		}

		if cnt > asm.MAXADDR {
			err = fmt.Errorf("loading rom: address overflow: %x", cnt)
			return
		}

		if hexval <= asm.MAXUINT {
			cpu.ROM[cnt] = uint16(hexval)
		} else {
			err = fmt.Errorf("loading rom: value too large for uint16: %x", hexval)
			return
		}
	}

	tools.Verbose("LOAD: loaded %d lines\n", cnt)
	return
}

func (cpu *HACKCPU) BLOAD(hexProgram []string) (err error) {
	var hexval int64
	var instr string
	var hackProgram []string

	for _, instr = range hexProgram {
		hexval, err = strconv.ParseInt(instr, 16, 64)
		if err != nil {
			return
		}
		hackProgram = append(hackProgram, fmt.Sprintf("%016b", uint16(hexval)))
	}

	err = cpu.LOAD(hackProgram)

	return
}

func (cpu *HACKCPU) RUN(hackProgram []string) (err error) {
	if len(hackProgram) > 0 {
		err = cpu.LOAD(hackProgram)
	}

	// starting clock
	cpu.Clock()

	return
}

func (cpu *HACKCPU) BRUN(hexProgram []string) (err error) {
	if len(hexProgram) > 0 {
		err = cpu.BLOAD(hexProgram)
	}

	// starting clock
	cpu.Clock()

	return
}

func (cpu *HACKCPU) UpdateA(val int16) {
	if AaddrIsValid(val) {
		cpu.A = val
	} else {
		tools.Error("address too large setting A register: %d\n", val)
		cpu.Reset()
	}
}

func (cpu *HACKCPU) UpdateD() {
	cpu.D = cpu.ALU
}

func (cpu *HACKCPU) UpdateM() {
	if MaddrIsValid(cpu.A) {
		cpu.RAM[cpu.A] = cpu.ALU
	} else {
		tools.Error("address too large writing RAM: %d\n", cpu.A)
		cpu.Reset()
	}
}

func (cpu *HACKCPU) ReadM() (memval int16) {
	if MaddrIsValid(cpu.A) {
		memval = cpu.RAM[cpu.A]
	} else {
		tools.Error("address too large accessing RAM: %d\n", cpu.A)
		cpu.Reset()
	}

	return
}

func (cpu *HACKCPU) Exec(instr uint16) (err error) {
	if (instr & 0x8000) == 0 {
		err = cpu.ExecA(instr)
	} else {
		err = cpu.ExecC(instr)
	}

	return
}

func (cpu *HACKCPU) ExecA(instr uint16) (err error) {
	cpu.A = int16(instr) & 0x7FFF
	cpu.PC += 1

	return
}

func (cpu *HACKCPU) ExecC(instr uint16) (err error) {
	dest, comp, jump, err := SplitC(instr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 1: COMP instruction: fill ALU
	err = ExecComp(cpu, comp)
	if err != nil {
		return
	}

	// 2: DEST instruction: transfer ALU to reg
	err = ExecDest(cpu, dest)
	if err != nil {
		return
	}

	// 3: JUMP instruction: set PC
	cpu.PC, err = ExecJump(cpu, jump)

	return
}

func (cpu *HACKCPU) Step() (err error) {
	cpu.CURINST = cpu.ROM[cpu.PC]

	// save old pc
	cpu.TRLASTPC = cpu.PC

	// exec instruction
	err = cpu.Exec(cpu.CURINST)

	return
}

func (cpu *HACKCPU) Clock() (err error) {
	// exec one instruction
	cpu.Step()

	// wait for clock tick
	ClockCycle(cpu.CLOCKMODE)

	// show status
	if cpu.TRACE {
		cpu.PrintStatus()
		cpu.PrintRAM(cpu.TRMEMLO, cpu.TRMEMHI)
	}

	return
}

func (cpu *HACKCPU) PrintStatus() {
	ALU := cpu.ALU
	PC := cpu.PC
	A := cpu.A
	D := cpu.D
	M := cpu.ReadM()
	N := asm.Bool2Bit(cpu.SR.Get(NFLAG))
	Z := asm.Bool2Bit(cpu.SR.Get(ZFLAG))
	O := asm.Bool2Bit(cpu.SR.Get(OFLAG))

	curStr, err := asm.DisassembleHex(cpu.CURINST)
	if err != nil {
		fmt.Println(err)
	}
	tools.Info(tools.Bold("CPU Status:")+"  Current Instuction: %016b  --> %s\n", cpu.CURINST, curStr)

	tools.Info("Register:  A %5d (%04X) [%016b] |  D %5d (%04X) [%016b] | M %5d (%04X) [%016b]\n", A, A, A, D, D, D, M, M, M)
	tools.Info("         ALU %5d (%04X) [%016b] | PC %5d (%04X) [%016b] | Flags N %1s | Z %1s | O %1s \n", ALU, ALU, ALU, PC, PC, PC, N, Z, O)
}

func (cpu *HACKCPU) PrintRAM(lowAddr, highAddr uint16) {
	tools.Info(tools.Bold("MEM Dump: %4X - %4X")+"\n", lowAddr, highAddr)

	var chunks uint16 = 16

	for l := lowAddr; l < highAddr; l += chunks {
		fmt.Printf("%4X:    ", l)

		for k := l; k < l+chunks; k++ {
			if k < highAddr {
				fmt.Printf("  %04X", cpu.RAM[k])
			} else {
				fmt.Print(k, k < highAddr)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (cpu *HACKCPU) PrintROM(lowAddr, highAddr uint16) {
	tools.Info(tools.Bold("ROM Dump: %4X - %4X")+"\n", lowAddr, highAddr)

	var chunks uint16 = 16

	for l := lowAddr; l < highAddr; l += chunks {
		fmt.Printf("%04X:    ", l)
		for k := l; k < l+chunks; k++ {
			if k < highAddr {
				fmt.Printf("  %04X", cpu.ROM[k])
			} else {
				fmt.Print(k, k < highAddr)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
