package main

import (
	m "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audiolib"
)

func main() {
	Audio, _ := audiolib.Open()
	Audio.BeatsPerMinute(160)
	Audio.LoadSounds()
	defer Audio.Close()

	s, _ := m.ParseSequence("(C D E)")
	for i := 0; i < 10; i++ {
		t := m.PitchBy{Semitones: i}
		Audio.Play(t.Transform(s))
	}
}
