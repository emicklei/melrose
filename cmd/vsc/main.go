package main

import (
	"fmt"
	"os"
)

// see Makefile how to run this

func main() {
	switch os.Args[1] {
	case "grammar":
		grammar()
	case "snippets":
		snippets()
	case "dslmd":
		dslmarkdown()
	default:
		fmt.Println("unknown cmd")
	}
}
