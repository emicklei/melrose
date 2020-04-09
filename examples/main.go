package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/emicklei/melrose"
	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/midi"
)

const echoNotes = true

var audio melrose.AudioDevice

var nr = flag.Int("nr", 0, "number of the example")

func main() {
	flag.Parse()
	var err error
	audio, err = midi.Open()
	check(err)
	defer audio.Close()

	switch *nr {
	case 1:
		example1()
	case 2:
		example2()
	default:
		fmt.Println("run with -nr 1 to hear example1")
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// go run main.go -nr 1
func example1() {
	y := MustParseSequence("E+ A- C5- B- C5- A- E+ F+ A- C5- B- C5- A- F-")

	audio.SetBeatsPerMinute(400)
	audio.Play(y, echoNotes)
	audio.Play(Pitch{Target: y, Semitones: 1}.S(), echoNotes)
	audio.Play(Pitch{Target: y, Semitones: 2}.S(), echoNotes)
	audio.Play(y, echoNotes)
}

// go run main.go -nr 2
func example2() {
	y := MustParseSequence("F♯2 C♯3 F♯3 A3 C♯ F♯")

	// play with Classics>Micro Pulse
	audio.SetBeatsPerMinute(280)
	p := Pattern{Target: y, Indices: []int{3, 4, 2, 5, 1, 6, 2, 5}}
	jp := Join{List: []Sequenceable{
		Repeat{Target: p, Times: 2},
		Repeat{Target: Pitch{Target: p, Semitones: 1}, Times: 2},
		Repeat{Target: Pitch{Target: p, Semitones: -2}, Times: 2},
		Repeat{Target: Pitch{Target: p, Semitones: 3}, Times: 1},
	}}
	audio.Play(Repeat{Target: jp, Times: 4}, echoNotes)
}
