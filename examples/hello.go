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

	note, _ := m.ParseNote("C5")
	Audio.Play(note.S(), 2*time.Second)
}
