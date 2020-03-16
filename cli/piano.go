package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/emicklei/melrose"
)

var pianoNotes = createPiano()

func createPiano() map[string]melrose.Note {
	piano := map[string]melrose.Note{}
	for octave := 3; octave < 6; octave++ {
		for _, each := range strings.Fields("C D E F G A B") {
			key := fmt.Sprintf("%s%d", each, octave)
			note, err := melrose.ParseNote(key)
			if err != nil {
				log.Fatal(err)
			} else {
				piano[key] = note
			}
		}
	}
	return piano
}
