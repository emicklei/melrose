package term

import (
	"log"

	tvp "github.com/emicklei/tviewplus"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// startUI blocks
func startUI(mon *Monitor) {
	textBg := tcell.NewRGBColor(25, 28, 32)
	dropBg := tcell.NewRGBColor(20, 23, 27)
	tview.Styles.PrimaryTextColor = tcell.ColorGray
	tview.Styles.ContrastBackgroundColor = dropBg

	app := tview.NewApplication()

	foc := tvp.NewFocusGroup(app)

	bpm := tvp.NewReadOnlyTextView(app, mon.BPM)

	inputDevice := tvp.NewDropDownView(foc, mon.InputDeviceList)
	inputDevice.SetLabel("  in ")
	inputDevice.SetTextOptions("", "", "", "▼", "---")
	pitchOnly := tvp.NewCheckboxView(foc, mon.EchoReceivedPitchOnly).SetLabel("pitch only ")
	pitchOnly.SetCheckedString("Y")

	inputSection := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(inputDevice, 0, 1, false).
		// space
		AddItem(tview.NewBox().SetBorderPadding(0, 0, 1, 0), 1, 1, false).
		AddItem(pitchOnly, 12, 0, false)

	outputDevice := tvp.NewDropDownView(foc, mon.OutputDeviceList)
	outputDevice.SetTextOptions("", "", "", "▼", "---")
	outputDevice.SetLabel(" out ")

	biab := tvp.NewReadOnlyTextView(app, mon.BIAB)

	sent := tvp.NewReadOnlyTextView(app, mon.Sent)
	sent.SetBackgroundColor(textBg)

	received := tvp.NewReadOnlyTextView(app, mon.Received)
	received.SetBackgroundColor(textBg)

	console := tvp.NewReadOnlyTextView(app, mon.Console)
	console.SetTextColor(tcell.ColorLightGray)
	console.SetBackgroundColor(textBg)

	clear := tvp.NewButtonView(foc).SetLabel("clear all")
	clear.SetSelectedFunc(func() {
		mon.Sent.Set("")
		mon.Received.Set("")
		mon.Console.Set("")
		foc.GetApplication().SetFocus(nil) // loose it
	})

	settings := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(NewStaticView(" [white]Melrōse "), 0, 1, false).
		AddItem(tview.NewBox().SetBorderPadding(0, 0, 1, 0), 1, 1, false).
		AddItem(NewStaticView("[yellow]BPM "), 4, 0, false).
		AddItem(bpm, 8, 0, false).
		AddItem(tview.NewBox().SetBorderPadding(0, 0, 1, 0), 1, 1, false).
		AddItem(NewStaticView("[yellow]BIAB "), 5, 0, false).
		AddItem(biab, 1, 1, false).
		AddItem(tview.NewBox().SetBorderPadding(0, 0, 1, 0), 1, 1, false).
		AddItem(clear, 11, 0, false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(settings, 1, 1, true).

		// console
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(NewStaticView(" [yellow]console"), 1, 1, false).
		AddItem(console, 0, 4, false).

		// sent
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(outputDevice, 1, 0, false).
		AddItem(sent, 0, 2, false).

		// received
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(inputSection, 1, 0, false).
		AddItem(received, 0, 2, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		log.Println(err)
	}
}

func NewStaticView(label string) *tview.TextView {
	return tview.NewTextView().SetDynamicColors(true).SetText(label)
}
