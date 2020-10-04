package main

import (
	"log"

	"github.com/emicklei/melrose/system"
)

func main() {
	ctx, err := system.Setup()
	if err != nil {
		log.Fatalln(err)
	}
	defer system.TearDown(ctx)
	mon := NewMonitor()
	startUI(mon)
}
