package main

import (

	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audio"
)

var Audio *audio.Device

// go run scale.go
func main() {
	Audio = new(audio.Device)
	Audio.Open()
	Audio.LoadSounds()
	defer Audio.Close()

	Audio.BeatsPerMinute(140)
	Audio.Play(C().Scale(Minor).S())
	Audio.Play(C().Scale(Major).S().Reversed())
}
