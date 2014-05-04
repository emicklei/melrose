package main

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

var Otto = otto.New()

func main() {
	setup()
	fmt.Println(help())
	openDevice()
	defer closeDevice()
	loop()
}

func help() string {
	return `
	melrose
`
}

func setup() {
	loadSounds()
	Otto.Set("play", playSequence)
	Otto.Set("tempo", tempo)
	Otto.Set("chord", chord)
	Otto.Set("scale", scale)
	Otto.Set("pitch", pitch)
	Otto.Set("reverse", reverse)
}
