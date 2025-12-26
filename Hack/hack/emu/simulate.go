package emu

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/hleinders/hack/asm"
	"github.com/rivo/tview"
)

var globMsg string

type TvApplication struct {
	CPU     *HACKCPU
	App     *tview.Application
	MainScr *tview.Grid
	HelpScr *tview.Grid
	ConfScr *tview.Grid
	RgD     *tview.TextView
	RgAM    *tview.TextView
	ALU     *tview.TextView
	PC      *tview.TextView
	ROM     *tview.List
	RAM     *tview.Table
	Header  *tview.TextView
	FooterL *tview.TextView
	FooterR *tview.TextView
	STEP    bool
	CLOCK   float32
}

func (tva *TvApplication) Update() bool {
	if tva.CPU.RUNNING {
		err := tva.CPU.Step()
		if err != nil {
			globMsg = "CPU ERROR"
			return false
		}
	} else {
		tva.STEP = true
	}

	return tva.UpdateScreen()
}

func (tva *TvApplication) UpdateScreen() bool {
	var flagN, flagZ string

	D := tva.CPU.D
	A := tva.CPU.A
	M := tva.CPU.ReadM()
	ALU := tva.CPU.ALU
	PC := tva.CPU.TRLASTPC
	curStr, _ := asm.DisassembleHex(tva.CPU.CURINST)
	jumped := ""

	if tva.CPU.SR.Get(NFLAG) {
		flagN = "[red]N[-:-:-:-]"
	} else {
		flagN = "N"
	}
	if tva.CPU.SR.Get(ZFLAG) {
		flagZ = "[red]Z[-:-:-:-]"
	} else {
		flagZ = "Z"
	}

	if (tva.CPU.PC != PC+1) && (PC != 0) {
		jumped = "[red]JMP[-:-:-:-]"
	} else {
		jumped = "   "
	}

	tva.RgD.Clear()
	fmt.Fprintf(tva.RgD, " D:  [[green]%016b[white]]\n", uint16(D))
	fmt.Fprintf(tva.RgD, "      0x%04X    %+06d", uint16(D), D)

	tva.RgAM.Clear()
	fmt.Fprintf(tva.RgAM, " A:  [[green]%016b[white]]\n", uint16(A))
	fmt.Fprintf(tva.RgAM, "      0x%04X     %05d\n", uint16(A), A)
	fmt.Fprintf(tva.RgAM, " M:  [[green]%016b[white]]\n", uint16(M))
	fmt.Fprintf(tva.RgAM, "      0x%04X    %+06d", uint16(M), M)

	tva.ALU.Clear()
	fmt.Fprintf(tva.ALU, " ISA:  [[::b]%14s[-:-:-:-]  ]\n\n", curStr)
	fmt.Fprintf(tva.ALU, " ALU:  [[green]%016b[white]]     (%s)\n", uint16(ALU), flagZ)
	fmt.Fprintf(tva.ALU, "        0x%04X    %+06d      (%s)", uint16(ALU), ALU, flagN)

	tva.PC.Clear()
	fmt.Fprintf(tva.PC, " PC:   [[green]%016b[white]]   [%s]\n", PC, jumped)
	fmt.Fprintf(tva.PC, "        0x%04X     %05d", uint16(PC), PC)

	// Header and Footer:
	tva.Header.Clear()
	if tva.STEP {
		if tva.CPU.RUNNING {
			fmt.Fprint(tva.Header, "\n  [::b]HACK CPU ready[-:-:-:-]")
		} else {
			fmt.Fprint(tva.Header, "\n  [::b]HACK CPU stopped![-:-:-:-]")
		}
	} else {
		fmt.Fprint(tva.Header, "\n  [::b]HACK CPU running....[-:-:-:-] ")
	}

	tva.FooterL.Clear()

	if globMsg != "" {
		fmt.Fprintf(tva.FooterL, "\n  [red]CPU-Error:[-:-:-:-] %s", globMsg)
		globMsg = ""
	} else {
		var stepMode string
		if tva.STEP {
			stepMode = "ON"
		} else {
			stepMode = "OFF"
		}
		fmt.Fprintf(tva.FooterL, "\n  Single Step Mode: [green]%s[white] / Clock: %f Hz", stepMode, tva.CLOCK)
	}

	lidx := int(PC)
	lrom := tva.ROM.GetItemCount()
	_, _, _, height := tva.ROM.GetRect()

	if lrom > height {
		offset, _ := tva.ROM.GetOffset()
		if height+offset-lidx <= 3 {
			offset += height - 6
			if offset+height > lrom {
				// already at end?
				offset = lrom - height
			}

			tva.ROM.SetOffset(offset, 0)
		} else if lidx-offset < 3 {
			// jumped back more than n=height lines?
			offset = lidx - 3

			tva.ROM.SetOffset(offset, 0)
		}

		// fmt.Fprintf(tva.Footer, "\nIdx: %d / Height: %d / Len: %d / Offset: %d", lidx, height, lrom, offset)
	}
	tva.ROM.SetCurrentItem(lidx)

	return false
}

func (tva *TvApplication) RAMPageUp() {
	var newOffset int
	curOffset, _ := tva.RAM.GetOffset()
	_, _, _, height := tva.RAM.GetRect()
	if curOffset-height > 0 {
		newOffset = curOffset - height
	} else {
		newOffset = 0
	}
	tva.RAM.SetOffset(newOffset, 0)
}

func (tva *TvApplication) RAMPageDown() {
	var newOffset int
	curOffset, _ := tva.RAM.GetOffset()
	_, _, _, height := tva.RAM.GetRect()
	if curOffset+height < asm.MAXMEM {
		newOffset = curOffset + height
	} else {
		newOffset = asm.MAXMEM - height
	}
	tva.RAM.SetOffset(newOffset, 0)
}

func (tva *TvApplication) GRAutoUpdate(done <-chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			if tva.STEP {
				tva.App.Draw()
				return
			}
			T := (1.0 / tva.CLOCK) * 1000.0 // set to ms
			time.Sleep(time.Duration(T) * time.Millisecond)
			tva.Update()
			tva.App.Draw()
		}
	}
}

// --------------------------------------------------------------------------------------------------------------------
type TableData struct {
	orgData []int16
	columns int
	tview.TableContentReadOnly
}

func (d *TableData) GetCell(row, column int) *tview.TableCell {
	if column%d.columns == 0 {
		// column header
		return tview.NewTableCell(fmt.Sprintf("  [green]0x%04X:[white]    %04X", uint16(row*d.columns), d.orgData[uint16(row*d.columns)+uint16(column)])).SetExpansion(1)

	}
	return tview.NewTableCell(fmt.Sprintf("%04X", d.orgData[uint16(row*d.columns)+uint16(column)])).SetExpansion(1)
}

func (d *TableData) GetRowCount() int {
	return asm.MAXMEM / d.columns
}

func (d *TableData) GetColumnCount() int {
	return d.columns
}

// ===============================================================================================================

func CreateSimulator(program []string, initialClock float32) error {
	app := tview.NewApplication()

	// Main Screen -----------------------------------------------------------------------------------------------
	regD := tview.NewTextView().SetDynamicColors(true) //.SetTitle(" Register D ").SetTitleAlign(tview.AlignLeft)
	regD.SetBorder(true).SetTitle(" Register D ").SetTitleAlign(tview.AlignLeft)

	regAM := tview.NewTextView().SetDynamicColors(true)
	regAM.SetBorder(true).SetTitle(" Register D ").SetTitleAlign(tview.AlignLeft)

	alu := tview.NewTextView().SetDynamicColors(true)
	alu.SetBorder(true).SetTitle(" Arithmetic Logical Unit ").SetTitleAlign(tview.AlignLeft)

	pc := tview.NewTextView().SetDynamicColors(true)
	pc.SetBorder(true).SetTitle(" Program Counter ").SetTitleAlign(tview.AlignLeft)

	RAM := tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false)

	RamGRID := tview.NewGrid().
		SetBorders(true).
		SetRows(0).
		SetColumns(0).
		AddItem(RAM, 0, 0, 1, 1, 0, 0, false)

	ROM := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedTextColor(tcell.ColorWhite).
		SetSelectedBackgroundColor(tcell.ColorDarkGreen).
		SetHighlightFullLine(true)

	for i, val := range program {
		hexval, _ := strconv.ParseInt(val, 16, 64)
		m, _ := asm.DisassembleHex(uint16(hexval))
		pline := fmt.Sprintf("%6d:    %-12s    %016b     %04X\n", i, m, uint16(hexval), uint16(hexval))
		ROM.AddItem(pline, "", 0, nil)
	}

	regBlock := tview.NewGrid().
		SetRows(4, 6).
		SetColumns(26).
		SetBorders(false).
		AddItem(regD, 0, 0, 1, 1, 0, 0, false).
		AddItem(regAM, 1, 0, 1, 1, 0, 0, false)

	aluBlock := tview.NewGrid().
		SetRows(6, 4).
		SetColumns(38).
		SetBorders(false).
		AddItem(alu, 0, 0, 1, 1, 0, 0, false).
		AddItem(pc, 1, 0, 1, 1, 0, 0, false)

	RAMBlock := tview.NewGrid().
		SetRows(0).
		SetColumns(64).
		SetBorders(false).
		AddItem(RamGRID, 0, 0, 1, 1, 0, 0, false)

	ROMBlock := tview.NewGrid().
		SetRows(0).
		SetColumns(0).
		SetBorders(false).
		AddItem(ROM, 0, 0, 1, 1, 0, 0, false)

	LeftBlock := tview.NewGrid().
		SetRows(10, 0).
		SetColumns(26, 42).
		SetBorders(false).
		AddItem(regBlock, 0, 0, 1, 1, 0, 0, false).
		AddItem(aluBlock, 0, 1, 1, 1, 0, 0, false).
		AddItem(RAMBlock, 1, 0, 1, 2, 0, 0, false)

	RightBlock := tview.NewGrid().
		SetRows(0).
		SetColumns(0).
		SetBorders(false).
		AddItem(ROMBlock, 0, 0, 1, 1, 0, 0, false)

	Header := tview.NewTextView().SetText("\n  [::b]HACK CPU ready[-:-:-:-]").SetDynamicColors(true)
	FooterL := tview.NewTextView().SetText("").SetDynamicColors(true)
	FooterR := tview.NewTextView().SetText("").SetDynamicColors(true).SetWordWrap(false).SetWrap(false)

	// help section:
	fmt.Fprintf(FooterR, ` Keys:  [yellow](r,s):[white] toggle single step/run mode | [yellow]R:[white] Reset
	    [yellow]arrow keys:[white] change clock | [yellow](Space,n):[white] next step
	    [yellow]PgUp,PgDn:[white] scroll RAM | [yellow](q,ESC):[white] quit`)

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(64, 0).
		SetBorders(true).
		AddItem(Header, 0, 0, 1, 2, 0, 0, false).
		AddItem(FooterL, 2, 0, 1, 1, 0, 0, false).
		AddItem(FooterR, 2, 1, 1, 1, 0, 70, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(LeftBlock, 1, 0, 1, 1, 0, 0, false).
		AddItem(RightBlock, 1, 1, 1, 1, 0, 70, false)

	// put everything together ----------------------------------------------------------------------------------

	hackCPU := HACKCPU{}
	hackCPU.Init()
	err := hackCPU.BLOAD(program)
	if err != nil {
		return err
	}

	app.SetRoot(grid, true).SetFocus(grid)

	tvApp := &TvApplication{
		CPU:     &hackCPU,
		App:     app,
		MainScr: grid,
		RgD:     regD,
		RgAM:    regAM,
		ALU:     alu,
		PC:      pc,
		ROM:     ROM,
		RAM:     RAM,
		Header:  Header,
		FooterL: FooterL,
		FooterR: FooterR,
		CLOCK:   initialClock,
		STEP:    true,
	}

	// prepare RAM data
	ramData := &TableData{orgData: tvApp.CPU.RAM, columns: 8}
	RAM.SetContent(ramData)

	// set hackCPU to running
	tvApp.CPU.RUNNING = true

	// initial update:
	tvApp.Update()

	// handle key press
	tvApp.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		done := make(chan bool, 10)

		switch event.Key() {
		case tcell.KeyCtrlC:
			tvApp.App.Stop()
		case tcell.KeyEsc: //, tcell.KeyEnter:
			tvApp.App.Stop()
		case tcell.KeyUp:
			if event.Modifiers()&tcell.ModShift != 0 {
				tvApp.CLOCK += 0.1
			} else {
				tvApp.CLOCK *= 2.0
			}
			tvApp.UpdateScreen()
		case tcell.KeyDown:
			if event.Modifiers()&tcell.ModShift != 0 {
				if tvApp.CLOCK > 0.1 {
					tvApp.CLOCK -= 0.1
				}
			} else {
				tvApp.CLOCK /= 2.0
			}
			tvApp.UpdateScreen()
		case tcell.KeyLeft:
			tvApp.CLOCK = 1.0
			tvApp.UpdateScreen()
		case tcell.KeyRight:
			tvApp.CLOCK = 5.0
			tvApp.UpdateScreen()
		case tcell.KeyPgUp:
			tvApp.RAMPageUp()
		case tcell.KeyPgDn:
			tvApp.RAMPageDown()
		}

		switch event.Rune() {
		case 'q':
			tvApp.App.Stop()
		case ' ', 'n':
			tvApp.STEP = true
			tvApp.Update()
		case 'r', 's':
			// toggle single step
			tvApp.STEP = !tvApp.STEP
			if !tvApp.STEP {
				go tvApp.GRAutoUpdate(done)
			} else {
				done <- true
			}
			tvApp.UpdateScreen()
		case 'R':
			tvApp.STEP = true
			tvApp.CPU.Reset()
			globMsg = "RESET (press space)"
			tvApp.CPU.RUNNING = true
			tvApp.CLOCK = 1.0
			tvApp.Update()
		default:
			tvApp.UpdateScreen()
		}

		// return event
		return nil
	})

	if err := tvApp.App.Run(); err != nil {
		panic(err)
	}

	return nil
}
