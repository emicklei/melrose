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

	left := newStreamer()
	right := newStreamer()

	note := m.ParseNote("C1")
	for i := 0; i < 40; i++ {
		left.put(note.Major(i))
		right.put(note.Major(39 - i))
	}
}

type streamer struct {
	notes chan m.Note
}

func (s *streamer) put(n m.Note) {
	s.notes <- n
}

func newStreamer() *streamer {
	s := new(streamer)
	s.notes = make(chan m.Note)
	go func() {
		for {
			Audio.PlayNote(<-s.notes, 300*time.Millisecond)
		}
	}()
	return s
}
