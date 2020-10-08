package term

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// FocusGroup is for rotating the focus among its members.
type FocusGroup struct {
	members    []tview.Primitive
	app        *tview.Application
	focusIndex int
}

// NewFocusGroup creates a FocusGroup for widgets in an application.
func NewFocusGroup(app *tview.Application) *FocusGroup {
	return &FocusGroup{
		app: app,
	}
}

// GetApplication
func (f *FocusGroup) GetApplication() *tview.Application {
	return f.app
}

// GetFocus returns the widget that currently has or should get the focus.
func (f *FocusGroup) GetFocus() tview.Primitive {
	if len(f.members) == 0 {
		return nil
	}
	return f.members[f.focusIndex]
}

// Add appends the widget to the members.
// The first will get the initial focus.
func (f *FocusGroup) Add(p tview.Primitive) {
	f.members = append(f.members, p)
}

func (f *FocusGroup) HandleDone(w tview.Primitive, k tcell.Key) {
	// find index for widget
	index := -1
	for i := 0; i <= len(f.members); i++ {
		if f.members[i] == w {
			index = i
			break
		}
	}
	if index == -1 {
		// not part of members
		return
	}
	// rotate
	switch k {
	default: // tab, enter, escape
		if index == len(f.members)-1 {
			f.focusIndex = 0
		} else {
			f.focusIndex = index + 1
		}
	case tcell.KeyBacktab:
		if index == 0 {
			f.focusIndex = len(f.members) - 1
		} else {
			f.focusIndex = index - 1
		}
	}
	// move focus
	f.app.SetFocus(f.GetFocus())
}
