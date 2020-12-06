package main

import (
	"log"

	"github.com/emicklei/melrose/system"
	"github.com/emicklei/melrose/ui/cli"
)

func main() {
	ctx, err := system.Setup()
	if err != nil {
		log.Fatalln(err)
	}
	defer system.TearDown(ctx)
	cli.StartREPL(ctx)
}
