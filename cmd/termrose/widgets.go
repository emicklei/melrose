package main

import (
	"github.com/rivo/tview"
)

func NewTextView(app *tview.Application, h *StringHolder) *tview.TextView {
	w := tview.NewTextView()
	w.SetChangedFunc(func() {
		app.Draw()
	})
	h.AddDependent(func(old, new string) {
		app.QueueUpdate(func() {
			w.SetText(new)
		})
	})
	return w
}

func NewInputView(app *tview.Application, h *StringHolder) *tview.InputField {
	w := tview.NewInputField()
	w.SetChangedFunc(func(text string) {
		//app.Draw()
	})
	h.AddDependent(func(old, new string) {
		app.QueueUpdate(func() {
			w.SetText(new)
		})
	})
	return w
}

func NewDropDownView(app *tview.Application, h *StringListSelectionHolder) *tview.DropDown {
	w := tview.NewDropDown()
	h.AddDependent(func(old, new []string) {
		app.QueueUpdate(func() {
			w.SetOptions(new, func(text string, index int) {
				h.Selection = text
				h.SelectionIndex = index
			})
		})
	})
	return w
}
