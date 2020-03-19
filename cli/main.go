package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audio"
	"github.com/peterh/liner"
)

var piano *audio.Device

var history = ".melrose.history"

var verbose = flag.Bool("v", false, "verbose logging")

func main() {
	flag.Parse()
	fmt.Println(help())
	piano = new(audio.Device)
	piano.Open()
	defer piano.Close()
	line := liner.NewLiner()
	defer line.Close()
	defer tearDown(line)
	setup(line)
	loop(line)
}

var functionNames = []string{"play"}

func help() string {
	return `
	melrose
	
	v = seq("C D E F A B C5")
	bpm(180)
	c = note("C4").Chord()
	play(v,v)
	
	:q
	:v
	:h
`
}

func tearDown(line *liner.State) {
	fmt.Println("closing melrose ...")
	if f, err := os.Create(history); err != nil {
		printError("error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func setup(line *liner.State) {
	fmt.Println("melrose - v0.0.1")
	line.SetCtrlCAborts(true)
	line.SetCompleter(func(line string) (c []string) {
		// if line ends with dot then lookup methods for the value before the dot
		if strings.HasSuffix(line, ".") {
			n := melrose.C()
			for _, each := range availableMethodNames(n, line) {
				c = append(c, line+each)
			}
			log.Println(c)
		} else {
			for _, n := range functionNames {
				if strings.HasPrefix(n, strings.ToLower(line)) {
					c = append(c, n)
				}
			}
		}
		return
	})
	if f, err := os.Open(history); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	piano.LoadSounds()
}

func loop(line *liner.State) {
	for {
		entry, err := line.Prompt("𝄞 ")
		if err != nil {
			printError(err)
			continue
		}
		switch entry {
		case "?", ":h":
			help()
		// commands starting with : control the program itself
		case ":q":
			goto exit
		case ":v":
			for k, v := range memory {
				fmt.Printf("%s = (%T) %v\n", k, v, v)
			}
		default:
			if err := dispatch(entry); err != nil {
				printError(err)
				// even on error, add entry to history so we can edit/fix it
			}
		}
		// do not remember commands
		if !strings.HasPrefix(entry, ":") {
			line.AppendHistory(entry)
		}
	}
exit:
}
