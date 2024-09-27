package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
	"github.com/peterh/liner"
)

var (
	history = ".melrose.history"
)

func StartREPL(ctx core.Context) {
	notify.PrintWelcome(core.BuildTag)
	// start REPL
	line := liner.NewLiner()
	defer line.Close()
	defer tearDown(line, ctx)
	// TODO liner catches control+c
	//setupCloseHandler(line)
	setup(line)
	repl(line, ctx)
}

func tearDown(line *liner.State, ctx core.Context) {
	ctx.Control().Reset()
	ctx.Device().Reset()
	if f, err := os.Create(history); err != nil {
		notify.Print(notify.NewErrorf("error writing history file:%v", err))
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func setup(line *liner.State) {
	line.SetCtrlCAborts(true)
	line.SetWordCompleter(completeMe)
	if f, err := os.Open(history); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

func repl(line *liner.State, ctx core.Context) {
	eval := dsl.NewEvaluator(ctx)
	ctx.Control().Start()
	ctx.Device().Command([]string{"i"}) // bit of a hack
	for {
		entry, err := line.Prompt(notify.Prompt())
		if err != nil {
			notify.Print(notify.NewError(err))
			tearDown(line, ctx)
			goto exit
		}
		entry = strings.TrimSpace(entry)
		if strings.HasPrefix(entry, ":") {
			// special case
			if entry == ":q" || entry == ":Q" {
				goto exit
			}
			args := strings.Split(entry, " ")
			if cmd, ok := lookupCommand(args[0]); ok {
				if msg := cmd.Func(ctx, args[1:]); msg != nil {
					notify.Print(msg)
				}
				line.AppendHistory(entry)
				continue
			}
		}
		if strings.HasSuffix(entry, "!") {

			if len(entry) == 1 {
				notify.Errorf("missing expression before '!'")
				continue
			}
			if result, err := eval.RecoveringEvaluateStatement(entry[:len(entry)-1]); err != nil {
				notify.Print(notify.NewError(err))
				continue
			} else {
				if result != nil {
					// create hidden variable
					// assign it the value of the expression before !
					// open the browser on it
					ctx.Variables().Put("_", result)
					open("http://localhost:8118/v1/notes?var=_")
					continue
				}
			}
		}
		if result, err := eval.RecoveringEvaluateStatement(entry); err != nil {
			notify.Print(notify.NewError(err))
			// even on error, add entry to history so we can edit/fix it
		} else {
			//log.Println("write inspection")
			core.InspectValue(ctx, result)
		}
		line.AppendHistory(entry)
	}
exit:
}

// setupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func setupCloseHandler(line *liner.State, ctx core.Context) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		tearDown(line, ctx)
		os.Exit(0)
	}()
}

// Open calls the OS default program for uri
func open(uri string) error {
	switch {
	case "windows" == runtime.GOOS:
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Start()
	case "darwin" == runtime.GOOS:
		return exec.Command("open", uri).Start()
	case "linux" == runtime.GOOS:
		return exec.Command("xdg-open", uri).Start()
	default:
		return fmt.Errorf("unable to open uri:%v on:%v", uri, runtime.GOOS)
	}
}
