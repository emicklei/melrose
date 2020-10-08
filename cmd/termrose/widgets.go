package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func NewStaticLabel(label string) *tview.TextView {
	w := tview.NewTextView()
	w.SetDynamicColors(true)
	w.SetText(label)
	return w
}

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

func NewInputView(f *FocusGroup, h *StringHolder) *tview.InputField {
	w := tview.NewInputField()
	w.SetChangedFunc(func(text string) {
		//app.Draw()
	})
	f.Add(w)
	w.SetDoneFunc(func(k tcell.Key) {
		f.HandleDone(w, k)
	})
	h.AddDependent(func(old, new string) {
		f.GetApplication().QueueUpdate(func() {
			w.SetText(new)
		})
	})
	return w
}

func NewDropDownView(f *FocusGroup, h *StringListSelectionHolder) *tview.DropDown {
	w := tview.NewDropDown()
	f.Add(w)
	w.SetDoneFunc(func(k tcell.Key) {
		f.HandleDone(w, k)
	})
	h.AddDependent(func(old, new []string) {
		f.GetApplication().QueueUpdate(func() {
			w.SetOptions(new, func(text string, index int) {
				h.Selection = text
				h.SelectionIndex = index
			})
		})
	})
	return w
}
