package main

import (
	"time"

	m "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audio"
)

var Audio *audio.Device

func main() {
	Audio = new(audio.Device)
	Audio.Open()
	Audio.LoadSounds()
	defer Audio.Close()

	s, _ := m.ParseSequence("(C D E)")
	for i := 0; i < 10; i++ {
		t := m.PitchBy{Semitones: i}
		Audio.Play(t.Transform(s), 1*time.Second)
	}
}
