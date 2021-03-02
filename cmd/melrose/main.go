package main

import (
	"log"

	"github.com/emicklei/melrose/system"
	"github.com/emicklei/melrose/ui/cli"
)

var BuildTag = "dev"

func main() {
	ctx, err := system.Setup(BuildTag)
	if err != nil {
		log.Fatalln(err)
	}
	defer system.TearDown(ctx)
	cli.StartREPL(ctx)
}
