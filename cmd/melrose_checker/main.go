package main

import (
	"io"
	"os"

	"github.com/emicklei/melrose/dsl"
)

func main() {
	src, _ := io.ReadAll(os.Stdin)
	err := dsl.Validate(string(src))
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}
}
