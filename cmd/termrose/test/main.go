package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	inputDevice := tview.NewDropDown().
		AddOption("abc", nil).
		AddOption("def", nil).
		SetFieldWidth(10)

	outputDevice := tview.NewDropDown().
		AddOption("123", nil).
		AddOption("456", nil).
		SetFieldWidth(10)

	settings := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(inputDevice, 0, 1, true).
		AddItem(outputDevice, 0, 1, true)

	top := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(settings, 0, 1, true)

	if err := app.SetRoot(top, true).SetFocus(top).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
