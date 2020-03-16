package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
	
	play("C#5 E_ 2F G A#")
	tempo(150)
	chord("C") -> "C E"
	scale("C") -> "C D E F G A B"
	pitch("C" , -1) -> "B3"
	reverse("C D E") -> "E D C"
	repeat("C",5)
	rotate("C D E",-1) -> "D E C"
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
		for _, n := range functionNames {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
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
		case "?":
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
		line.AppendHistory(entry)
	}
exit:
}
