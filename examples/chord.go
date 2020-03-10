package main

import (
	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audio"
)

// go run chord.go
func main() {
	Audio := new(audio.Device)
	Audio.Open()
	Audio.LoadSounds()
	defer Audio.Close()

	csm := C(Sharp).Chord(Minor).Octaved(-1)
	b1 := B().Chord(Major).Octaved(-1)

	Audio.Play(csm.Join(b1))
}
