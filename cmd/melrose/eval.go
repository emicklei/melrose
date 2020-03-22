package main

import (
	"fmt"

	"github.com/emicklei/melrose"
)

var evalFuncMap = evalFunctions()

type Function struct {
	Description string
	Sample      string
	Func        interface{}
}

func evalFunctions() map[string]Function {
	eval := map[string]Function{}
	eval["chord"] = Function{
		Description: "create a triad Chord with a Note",
		Sample:      `chord("C4")`,
		Func: func(note string) melrose.Chord {
			n, err := melrose.ParseNote(note)
			if err != nil {
				printError(err)
				return melrose.N("C").Chord()
			}
			return n.Chord()
		}}

	eval["join"] = Function{
		Description: "join two or more Musical objects",
		Sample:      `join()`,
		Func: func(playables ...interface{}) interface{} {
			joined := []melrose.Sequenceable{}
			for _, p := range playables {
				if s, ok := p.(melrose.Sequenceable); ok {
					joined = append(joined, s)
				} else {
					printWarning(fmt.Sprintf("cannot join (%T) %v", p, p))
				}
			}
			return melrose.Join{List: joined}
		}}

	eval["bpm"] = Function{
		Description: "set the Beats Per Minute value [1..300], default is 120",
		Sample:      `bpm(180)`,
		Func: func(i int) int {
			currentDevice.BeatsPerMinute(float64(i))
			return i
		}}

	eval["seq"] = Function{
		Description: "create a Sequence from a string of notes",
		Sample:      `seq("C C5")`,
		Func: func(s string) melrose.Sequence {
			n, err := melrose.ParseSequence(s)
			if err != nil {
				printError(err)
				return melrose.N("C").S()
			}
			return n
		}}

	eval["note"] = Function{
		Description: "create a Note from a string",
		Sample:      `note("C#3")`,
		Func: func(s string) melrose.Note {
			n, err := melrose.ParseNote(s)
			if err != nil {
				printError(err)
				return melrose.N("C")
			}
			return n
		}}

	eval["play"] = Function{
		Description: "play a Musical object such as Note,Chord,Sequence,...",
		Sample:      `play()`,
		Func: func(playables ...interface{}) interface{} {
			for _, p := range playables {
				if s, ok := p.(melrose.Sequenceable); ok {
					currentDevice.Play(s.S())
				} else {
					printWarning(fmt.Sprintf("cannot play (%T) %v", p, p))
				}
			}
			return nil
		}}

	eval["var"] = Function{
		Description: "create a reference to a known variable",
		Sample:      `var(v1)`,
		Func: func(value interface{}) Variable {
			varName := varStore.NameFor(value)
			if len(varName) == 0 {
				printWarning("no variable found with this Musical Object")
				return Variable{Name: "?", store: varStore}
			}
			return Variable{Name: varName, store: varStore}
		}}
	return eval
}
