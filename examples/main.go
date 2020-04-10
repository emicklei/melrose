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
	audio.SetBeatsPerMinute(200)
	defer audio.Close()

	switch *nr {
	case 1:
		example1()
	case 2:
		example2()
	case 3:
		example3()
	case 4:
		example4()
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
	p := Fluent{}.IndexMap(y, "3 4 2 5 1 6 2 5")
	jp := Join{List: []Sequenceable{
		Repeat{Target: p, Times: 2},
		Repeat{Target: Pitch{Target: p, Semitones: 1}, Times: 2},
		Repeat{Target: Pitch{Target: p, Semitones: -2}, Times: 2},
		Repeat{Target: Pitch{Target: p, Semitones: 3}, Times: 2},
	}}
	audio.Play(Repeat{Target: jp, Times: 4}, echoNotes)
}

// go run main.go -nr 3
func example3() {
	f := Fluent{} // use CLI-like syntax

	f1 := f.Sequence("C D E C")
	f2 := f.Sequence("E F ½G")
	f3 := f.Sequence("⅛G ⅛A ⅛G ⅛F E C")
	f4 := f.Sequence("½C ½G3 ½C 1=")
	r8 := f.Repeat(8, f.Rest())

	v1 := f.Join(f1, f1, f2, f2, f3, f3, f4)
	v2 := f.Join(r8, v1)
	v3 := f.Join(r8, v2)

	f.Go(audio, v1, v2, v3)
}

var m Fluent

// go run main.go -nr 4
func example4() {
	y := m.Sequence("C♯2 G♯2 C♯3 E3 G♯3 E")

	audio.SetBeatsPerMinute(280)
	mapped := m.IndexMap(y, "1 (4 5 6) 2 (4 5 6) 3 (4 5 6) 2 (4 5 6)")

	boomTikTik := m.Sequence("2C2 F#2 F#2")
	left := m.Repeat(8, boomTikTik)
	right := m.Repeat(4, mapped)

	m.Go(
		audio,
		midi.Channel(left, 1),
		midi.Channel(right, 2))
}
