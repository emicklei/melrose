package main

import (
	"time"

	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audio"
)

// go run scale.go
func main() {
	a := new(audio.Device)
	a.Open()
	a.LoadSounds()
	defer a.Close()

	a.BeatsPerMinute(180)

	a.Play(C().Scale(Minor).S())
	a.Play(C().Scale(Major).S().Reversed())

	time.Sleep(4 * time.Second)
}
