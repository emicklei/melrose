package main

import (
	"flag"
	"log"

	"github.com/emicklei/melrose/system"
	"github.com/emicklei/melrose/ui/cli"
	"github.com/emicklei/melrose/ui/term"
)

var (
	oCLI = flag.Bool("cli", true, "use the command line interface")
)

func main() {
	ctx, err := system.Setup()
	if err != nil {
		log.Fatalln(err)
	}
	defer system.TearDown(ctx)
	if *oCLI {
		cli.StartREPL(ctx)
	} else {
		term.NewMonitor().Open(ctx)
	}
}
