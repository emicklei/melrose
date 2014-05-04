package main

import (
	"fmt"
	"sync"
	"time"

	m "github.com/emicklei/melrose"
	"github.com/robertkrimen/otto"
)

var noteLength = time.Duration(500) * time.Millisecond

func playSequence(call otto.FunctionCall) otto.Value {
	arg := call.Argument(0)
	if arg.IsString() {
		input := arg.String()
		seq := m.ParseSequence(input)
		for _, eachGroup := range seq.Notes {
			wg := new(sync.WaitGroup)
			for _, eachNote := range eachGroup {
				wg.Add(1)
				go func(n m.Note) {
					playNote(n, noteLength)
					wg.Done()
				}(eachNote)
			}
			wg.Wait()
		}
	}
	return toValue("")
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
	start := m.ParseNote(starts)
	seq := m.Chord(start, mm)
	return toValue(seq.String())
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
	start := m.ParseNote(starts)
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
	seq := m.ParseSequence(notes)
	pitched := m.PitchBy{int(offset)}.Transform(seq)
	return toValue(pitched.String())
}
