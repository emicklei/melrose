package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

func main() {
	app := new(views.Application)
	win := &window{app: app}

	inner := views.NewBoxLayout(views.Horizontal)

	input := views.NewText()
	input.SetText("test")
	input.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorYellow).
		Background(tcell.ColorDarkGray))
	inner.AddWidget(input, 1)

	win.AddWidget(inner, 0)

	app.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorYellow).
		Background(tcell.ColorDarkGray))

	app.SetRootWidget(win)
	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}

type window struct {
	app *views.Application
	views.Panel
}

func (w *window) HandleEvent(e tcell.Event) bool {

	switch ev := e.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'Q', 'q':
				w.app.Quit()
				return true
			}
		}
	}

	return w.Panel.HandleEvent(e)
}
