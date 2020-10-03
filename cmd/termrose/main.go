package main

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	mon := NewMonitor()

	bpm := NewInputView(app, mon.BPM)
	bpm.SetLabel("BPM:")

	beat := NewTextView(app, mon.Beat)
	beat.SetTextColor(tcell.ColorLightCyan)
	beat.SetBackgroundColor(tcell.NewRGBColor(33, 37, 46))

	sent := NewTextView(app, mon.Sent)
	sent.SetBackgroundColor(tcell.NewRGBColor(33, 37, 46))

	received := NewTextView(app, mon.Sent)
	received.SetBackgroundColor(tcell.NewRGBColor(25, 28, 32))

	inputDevice := NewDropDownView(app, mon.InputDeviceList)
	inputDevice.SetLabel("Input:")
	inputDevice.SetFieldWidth(40)

	outputDevice := NewDropDownView(app, mon.OutputDeviceList)
	outputDevice.SetLabel("Output:")
	outputDevice.SetFieldWidth(40)

	settings := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(beat, 4, 1, true).
		// spacer left
		AddItem(tview.NewBox().SetBorderPadding(0, 0, 1, 0), 1, 1, false).
		AddItem(bpm, 0, 1, true).
		AddItem(inputDevice, 0, 1, true).
		AddItem(outputDevice, 0, 1, true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(settings, 1, 1, true).
		// spacer
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(sent, 0, 2, false).
		// spacer
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(received, 0, 2, false)

	go func() {
		mon.SetBPM(120)
		mon.SetBeat(1)
		mon.AppendSent("A B C")
		time.Sleep(1 * time.Second)
		mon.SetBPM(140)
		mon.SetInputDevices([]string{"abc", "def"})
		mon.SetOutputDevices([]string{"abfdsafdsaf fd f ds f dsa fd saf dsa fds a fdsc", "defdsafdsafdsafdsafdsafdsafdsafdsafdsafdsafdsf"})
		for i := 0; i < 100; i++ {
			mon.AppendSent("(D E F)")
		}
	}()

	go func() {
		for {
			for i := 1; i < 5; i++ {
				mon.SetBeat(i)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	if err := app.SetRoot(flex, true).SetFocus(flex).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
