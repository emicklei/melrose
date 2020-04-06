package main

import (
	"log"

	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/midi"
)

const echoNotes = true

func main() {
	audio, err := midi.Open()
	check(err)
	defer audio.Close()

	y := MustParseSequence("E+ A- C5- B- C5- A- E+ F+ A- C5- B- C5- A- F-")

	audio.SetBeatsPerMinute(400)
	audio.Play(y, echoNotes)
	audio.Play(Pitch{Target: y, Semitones: 1}.S(), echoNotes)
	audio.Play(Pitch{Target: y, Semitones: 2}.S(), echoNotes)
	audio.Play(y, echoNotes)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
