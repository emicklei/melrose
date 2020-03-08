package main

import (
	"time"

	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audio"
)

var Audio *audio.Device

// go run chord.go
func main() {
	Audio = new(audio.Device)
	Audio.Open()
	Audio.LoadSounds()
	defer Audio.Close()

	csm := Chord(C(Sharp), Minor).Octaved(-1)
	b1 := Chord(B(), Major).Octaved(-1)

	Audio.Play(csm.Join(b1), 6*time.Second)
}
