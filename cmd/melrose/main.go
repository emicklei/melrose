package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/emicklei/melrose/server"
	"github.com/emicklei/melrose/system"
	"github.com/emicklei/melrose/ui/cli"
)

var (
	BuildTag     = "dev"
	playFilename = flag.String("play", "", "script file to play")
	loadFilename = flag.String("load", "", "script file to load")
)

func main() {
	ctx, err := system.Setup(BuildTag)
	if err != nil {
		log.Fatalln(err)
	}
	defer system.TearDown(ctx)
	// if a file is specified, execute it a
	if *playFilename != "" {
		if err := cli.ExecuteFile(ctx, *playFilename); err != nil {
			log.Fatalln(err)
		}
		// wait for the user to press enter before exiting
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}
	if *loadFilename != "" {
		if err := cli.LoadFile(ctx, *loadFilename); err != nil {
			log.Fatalln(err)
		}
	}
	// start the server and the REPL
	server.Start(ctx)
	cli.StartREPL(ctx)
}
