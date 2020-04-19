package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/melrose/server"
	"github.com/peterh/liner"
)

var (
	version   = "v0.1"
	verbose   = flag.Bool("v", false, "verbose logging")
	inputFile = flag.String("i", "", "read expressions from a file")
	httpPort  = flag.String("http", ":8118", "address on which to listen for HTTP requests")

	history                      = ".melrose.history"
	varStore dsl.VariableStorage = dsl.NewVariableStore()
)

func main() {
	welcome()
	flag.Parse()

	// set audio
	currentDevice := setupAudio("midi")
	defer currentDevice.Close()
	melrose.SetCurrentDevice(currentDevice)

	// process file if given
	if len(*inputFile) > 0 {
		processInputFile(*inputFile)
	}

	if len(*httpPort) > 0 {
		// start DSL server
		go server.NewLanguageServer(varStore, *httpPort).Start()
	}

	// start REPL
	line := liner.NewLiner()
	defer line.Close()
	defer tearDown(line)
	// TODO liner catches control+c
	//setupCloseHandler(line)
	setup(line)
	loop(line)
}

func welcome() {
	fmt.Println("\033[1;34mmelrose\033[0m" + " - program your melody")
}

var functionNames = []string{"play"}

func tearDown(line *liner.State) {
	dsl.StopAllLoops(varStore)
	if f, err := os.Create(history); err != nil {
		notify.Print(notify.Errorf("error writing history file:%v", err))
	} else {
		line.WriteHistory(f)
		f.Close()
	}
	fmt.Println("\033[1;34mmelrose\033[0m" + " sings bye!")
}

func setup(line *liner.State) {
	line.SetCtrlCAborts(true)
	line.SetWordCompleter(completeMe)
	if f, err := os.Open(history); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

func loop(line *liner.State) {
	eval := dsl.NewEvaluator(varStore)
	for {
		entry, err := line.Prompt("𝄞 ")
		if err != nil {
			notify.Print(notify.Error(err))
			continue
		}
		entry = strings.TrimSpace(entry)
		if strings.HasPrefix(entry, ":") {
			// special case
			if entry == ":q" || entry == ":Q" {
				goto exit
			}
			args := strings.Split(entry, " ")
			if cmd, ok := lookupCommand(args[0]); ok {
				if msg := cmd.Func(args[1:]); msg != nil {
					notify.Print(msg)
				}
				continue
			}
		}
		if result, err := eval.Dispatch(entry); err != nil {
			notify.Print(notify.Error(err))
			// even on error, add entry to history so we can edit/fix it
		} else {
			if result != nil {
				printValue(result)
			}
		}
		line.AppendHistory(entry)
	}
exit:
}

func processInputFile(fileName string) {
	data, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		notify.Print(notify.Errorf("unable to read file:%v", err))
		return
	}
	eval := dsl.NewEvaluator(varStore)
	for line, each := range strings.Split(string(data), "\n") {
		entry := strings.TrimSpace(each)
		if _, err := eval.Dispatch(entry); err != nil {
			notify.Print(notify.Errorf("line %d:%v", line, err))
		}
	}
}

func printValue(v interface{}) {
	if v == nil {
		return
	}
	if s, ok := v.(melrose.Storable); ok {
		fmt.Printf("\033[94m(%T)\033[0m %s\n", v, s.Storex())
	} else {
		fmt.Printf("\033[94m(%T)\033[0m %v\n", v, v)
	}
}

// setupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func setupCloseHandler(line *liner.State) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		tearDown(line)
		os.Exit(0)
	}()
}
