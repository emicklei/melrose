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

	note := m.ParseNote("C1")
	for i := 0; i < 20; i++ {
		Audio.PlayNote(note.Major(i), 150*time.Millisecond)
	}
}
