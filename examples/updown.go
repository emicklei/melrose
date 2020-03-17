package main

import (
	m "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/audio"
)

var Audio *audio.Device

func main() {
	Audio = new(audio.Device)
	Audio.Open()
	Audio.LoadSounds()
	Audio.BeatsPerMinute(180)
	defer Audio.Close()

	left := newStreamer()
	right := newStreamer()

	note, _ := m.ParseNote("C1")
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
			Audio.Play((<-s.notes).S())
		}
	}()
	return s
}
