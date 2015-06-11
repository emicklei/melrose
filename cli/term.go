package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/GeertJohan/go.linenoise"
	"github.com/robertkrimen/otto"
)

func dispatch(entry string) string {
	if len(entry) == 0 {
		return entry
	}
	if strings.HasSuffix(entry, ";") {
		entry += ";"
	}
	value, err := Otto.Run(entry)
	if err != nil {
		return err.Error()
	}
	if value.IsUndefined() {
		return ""
	}
	s := fmt.Sprintf("%v", value)
	// because of missing IsEmpty check
	if "%!v(PANIC=toString(<nil> <nil>))" == s {
		return ""
	}
	return s
}

var lastHistoryEntry string

func loop() {
	linenoise.LoadHistory(".studio")
	for {
		entered, err := linenoise.Line("> ")
		if err != nil {
			if err == linenoise.KillSignalError {
				os.Exit(0)
			}
			fmt.Println("Unexpected error: %s", err)
			os.Exit(0)
		}
		entry := strings.TrimLeft(entered, "\t ") // without tabs,spaces
		var output string
		switch entry {
		case "?":
			output = help()
		case "q":
			Audio.Close()
			os.Exit(1)
		default:
			if entry != lastHistoryEntry {
				err = linenoise.AddHistory(entry)
				if err != nil {
					fmt.Printf("error: %s\n", entry)
				}
				lastHistoryEntry = entry
				linenoise.SaveHistory(".studio")
			}
			output = dispatch(entry)
		}
		if len(output) > 0 {
			fmt.Println(output)
		}
	}
}

func toValue(value interface{}) otto.Value {
	result, err := otto.ToValue(value)
	if err != nil {
		return toValue(err)
	}
	return result
}
