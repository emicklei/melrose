package term

import (
	"log"

	tvp "github.com/emicklei/tviewplus"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// startUI blocks
func startUI(mon *Monitor) {
	app := tview.NewApplication()

	foc := tvp.NewFocusGroup(app)

	bpm := tvp.NewTextView(app, mon.BPM)

	inputDevice := tvp.NewDropDownView(foc, mon.InputDeviceList)
	inputDevice.SetLabel("Input:")
	inputDevice.SetFieldWidth(40)

	outputDevice := tvp.NewDropDownView(foc, mon.OutputDeviceList)
	outputDevice.SetLabel("Output:")
	outputDevice.SetFieldWidth(40)

	beat := tvp.NewTextView(app, mon.Beat)
	beat.SetTextColor(tcell.ColorLightCyan)
	beat.SetBackgroundColor(tcell.NewRGBColor(33, 37, 46))

	sent := tvp.NewTextView(app, mon.Sent)
	sent.SetBackgroundColor(tcell.NewRGBColor(33, 37, 46))

	received := tvp.NewTextView(app, mon.Received)
	received.SetBackgroundColor(tcell.NewRGBColor(25, 28, 32))

	console := tvp.NewTextView(app, mon.Console)
	console.SetBackgroundColor(tcell.NewRGBColor(25, 28, 32))

	settings := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tvp.NewStaticView(" Melrōse "), 0, 1, false).
		AddItem(beat, 4, 1, false).
		AddItem(tview.NewBox().SetBorderPadding(0, 0, 1, 0), 1, 1, false).
		AddItem(bpm, 3, 0, false).
		// spacer right
		AddItem(tview.NewBox().SetBorderPadding(0, 0, 1, 0), 1, 1, false).
		AddItem(inputDevice, 0, 1, false).
		AddItem(outputDevice, 0, 1, false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(settings, 1, 1, true).

		// sent
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(tvp.NewStaticView(" [yellow]sent"), 1, 1, false).
		AddItem(sent, 0, 2, false).

		// received
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(tvp.NewStaticView(" [yellow]received"), 1, 1, false).
		AddItem(received, 0, 2, false).

		// console
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(tvp.NewStaticView(" [yellow]console"), 1, 1, false).
		AddItem(console, 0, 4, false)

	if err := app.SetRoot(flex, true).SetFocus(foc.GetFocus()).EnableMouse(true).Run(); err != nil {
		log.Println(err)
	}
}
