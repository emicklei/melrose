package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	m "github.com/emicklei/melrose"
	"github.com/robertkrimen/otto"
)

var noteLength = time.Duration(500) * time.Millisecond

func playAllSequences(call otto.FunctionCall) otto.Value {
	for i := 0; i < len(call.ArgumentList); i++ {
		argIndex := i
		arg := call.Argument(argIndex)
		if arg.IsString() {
			playSequence(arg.String())
		}
	}
	return toValue("")
}

func playSequence(input string) {
	seq, err := m.ParseSequence(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse sequence: %s\n", err)
		return
	}

	for _, eachGroup := range seq.Notes {
		wg := new(sync.WaitGroup)
		for _, eachNote := range eachGroup {
			wg.Add(1)
			go func(n m.Note) {
				Audio.PlayNote(n, noteLength)
				wg.Done()
			}(eachNote)
		}
		wg.Wait()
	}
}

func tempo(call otto.FunctionCall) otto.Value {
	bpm, err := call.Argument(0).ToInteger()
	if err != nil {
		return toValue(err)
	}
	noteLength = time.Millisecond * time.Duration(int64(1000*60/bpm))
	return toValue(fmt.Sprintf("%d", noteLength))
}

// ABCMouse.com
func chord(call otto.FunctionCall) otto.Value {
	starts, err := call.Argument(0).ToString()
	if err != nil {
		return toValue(err)
	}
	mms, err := call.Argument(1).ToString()
	if err != nil {
		return toValue(err)
	}
	mm := m.Major
	if mms == "m" {
		mm = m.Minor
	}
	start, err := m.ParseNote(starts)
	if err != nil {
		return toValue(err)
	}
	seq := m.Chord(start, mm)
	return toValue(seq.String())
}

// repeat(sequence, howMany)
func repeat(call otto.FunctionCall) otto.Value {
	seq, err := call.Argument(0).ToString()
	if err != nil {
		return toValue(err)
	}
	howMany, err := call.Argument(1).ToInteger()
	if err != nil {
		return toValue(err)
	}
	for howMany > 0 {
		playSequence(seq)
		howMany -= 1
	}
	return otto.NullValue()
}

func scale(call otto.FunctionCall) otto.Value {
	starts, err := call.Argument(0).ToString()
	if err != nil {
		return toValue(err)
	}
	mms, err := call.Argument(1).ToString()
	if err != nil {
		return toValue(err)
	}
	mm := m.Major
	if mms == "m" {
		mm = m.Minor
	}
	start, err := m.ParseNote(starts)
	if err != nil {
		return toValue(err)
	}
	seq := m.Scale(start, mm)
	return toValue(seq.String())
}

func pitch(call otto.FunctionCall) otto.Value {
	notes, err := call.Argument(0).ToString()
	if err != nil {
		return toValue(err)
	}
	offset, err := call.Argument(1).ToInteger()
	if err != nil {
		return toValue(err)
	}
	seq, err := m.ParseSequence(notes)
	if err != nil {
		return toValue(err)
	}
	pitched := m.PitchBy{int(offset)}.Transform(seq)
	return toValue(pitched.String())
}

// rotate("C D E", -4)
func rotate(call otto.FunctionCall) otto.Value {
	notes, err := call.Argument(0).ToString()
	if err != nil {
		return toValue(err)
	}
	howMany, err := call.Argument(2).ToInteger()
	if err != nil {
		return toValue(err)
	}
	seq, err := m.ParseSequence(notes)
	if err != nil {
		return toValue(err)
	}
	direction := m.Right
	if howMany < 0 {
		direction = m.Left
		howMany *= -1
	}
	t := seq.RotatedBy(direction, int(howMany))
	return toValue(t.String())
}

func reverse(call otto.FunctionCall) otto.Value {
	notes, err := call.Argument(0).ToString()
	if err != nil {
		return toValue(err)
	}
	seq, err := m.ParseSequence(notes)
	if err != nil {
		return toValue(err)
	}
	reversed := seq.Reversed()
	return toValue(reversed.String())
}
