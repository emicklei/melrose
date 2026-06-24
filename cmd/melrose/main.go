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
	BuildTag = "dev"
	fileName = flag.String("file", "", "script file to execute")
)

func main() {
	ctx, err := system.Setup(BuildTag)
	if err != nil {
		log.Fatalln(err)
	}
	defer system.TearDown(ctx)
	// if a file is specified, execute it a
	if *fileName != "" {
		if err := cli.ExecuteFile(ctx, *fileName); err != nil {
			log.Fatalln(err)
		}
		// wait for the user to press enter before exiting
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}
	// start the server and the REPL
	server.Start(ctx)
	cli.StartREPL(ctx)
}
