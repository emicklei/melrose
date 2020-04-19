package main

import (
	"os"
)

// see Makefile how to run this

func main() {
	switch os.Args[1] {
	case "grammar":
		grammar()
	case "snippets":
		snippets()
	}
}
