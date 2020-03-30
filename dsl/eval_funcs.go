package dsl

import (
	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

type Function struct {
	Description string
	Sample      string
	Func        interface{}
}

func EvalFunctions(varStore *VariableStore) map[string]Function {
	eval := map[string]Function{}
	eval["chord"] = Function{
		Description: "create a triad Chord with a Note",
		Sample:      `chord("C4")`,
		Func: func(note string) FunctionResult {
			n, err := melrose.ParseNote(note)
			if err != nil {
				return result(nil, notify.Errorf("%v", err))
			}
			return result(n.Chord(), nil)
		}}

	eval["pitch"] = Function{
		Description: "change the pitch with a delta of semitones",
		Sample:      `pitch(1,?)`,
		Func: func(semitones int, m interface{}) FunctionResult {
			s, ok := getSequenceable(m)
			if !ok {
				return result(nil, notify.Warningf("cannot pitch (%T) %v", m, m))
			}
			return result(melrose.Pitch{Target: s, Semitones: semitones}, nil)
		}}

	eval["reverse"] = Function{
		Description: "reverse the (groups of) notes in a sequence",
		Sample:      `reverse(?)`,
		Func: func(m interface{}) FunctionResult {
			s, ok := getSequenceable(m)
			if !ok {
				return result(nil, notify.Warningf("cannot reverse (%T) %v", m, m))
			}
			return result(melrose.Reverse{Target: s}, nil)
		}}

	eval["repeat"] = Function{
		Description: "repeat the musical object a number of times",
		Sample:      `repeat(2,?)`,
		Func: func(howMany int, m interface{}) FunctionResult {
			s, ok := getSequenceable(m)
			if !ok {
				return result(nil, notify.Warningf("cannot repeat (%T) %v", m, m))
			}
			return result(melrose.Repeat{Target: s, Times: howMany}, nil)
		}}

	eval["join"] = Function{
		Description: "join two or more musical objects",
		Sample:      `join(?,?)`,
		Func: func(playables ...interface{}) interface{} { // Note: return type cannot be EvaluationResult
			joined := []melrose.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					joined = append(joined, s)
				} else {
					return result(nil, notify.Warningf("cannot join (%T) %v", p, p))
				}
			}
			return result(melrose.Join{List: joined}, nil)
		}}

	eval["bpm"] = Function{
		Description: "get or set the Beats Per Minute value [1..300], default is 120",
		Sample:      `bpm(180)`,
		Func: func(f ...float64) FunctionResult {
			if len(f) == 0 {
				return result(melrose.CurrentDevice().BeatsPerMinute(), nil)
			}
			melrose.CurrentDevice().SetBeatsPerMinute(f[0])
			return result(f[0], nil)
		}}

	eval["seq"] = Function{
		Description: "create a Sequence from a string of notes",
		Sample:      `seq("C C5")`,
		Func: func(s string) FunctionResult {
			n, err := melrose.ParseSequence(s)
			if err != nil {
				return result(nil, notify.Error(err))
			}
			return result(n, nil)
		}}

	eval["note"] = Function{
		Description: "create a Note from a string",
		Sample:      `note("C#3")`,
		Func: func(s string) FunctionResult {
			n, err := melrose.ParseNote(s)
			if err != nil {
				return result(nil, notify.Error(err))
			}
			return result(n, nil)
		}}

	eval["play"] = Function{
		Description: "play a musical object such as Note,Chord,Sequence,...",
		Sample:      `play()`,
		Func: func(playables ...interface{}) interface{} { // Note: return type cannot be EvaluationResult
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					melrose.CurrentDevice().Play(s.S(), true)
				} else {
					return result(nil, notify.Warningf("cannot play (%T) %v", p, p))
				}
			}
			return result(nil, nil)
		}}

	eval["go"] = Function{
		Description: "play all musical objects in parallel",
		Sample:      `go()`,
		Func: func(playables ...interface{}) interface{} { // Note: return type cannot be EvaluationResult
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					go melrose.CurrentDevice().Play(s.S(), false)
				}
			}
			return result(nil, nil)
		}}

	eval["var"] = Function{
		Description: "create a reference to a known variable",
		Sample:      `var(v1)`,
		Func: func(value interface{}) FunctionResult {
			varName := varStore.NameFor(value)
			if len(varName) == 0 {
				return result(nil, notify.Warningf("no variable found with this Musical Object"))
			}
			return result(variable{Name: varName, store: varStore}, nil)
		}}

	eval["del"] = Function{
		Description: "delete a variable",
		Sample:      `del(v1)`,
		Func: func(value interface{}) FunctionResult {
			varName := varStore.NameFor(value)
			varStore.Delete(varName)
			return result(value, notify.Infof("deleted %s", varName))
		}}

	eval["flat"] = Function{
		Description: "flat (ungroup) the groups of a variable",
		Sample:      `flat(v1)`,
		Func: func(value interface{}) FunctionResult {
			if s, ok := getSequenceable(value); ok {
				return result(melrose.Ungroup{Target: s}, nil)
			} else {
				return result(nil, notify.Warningf("cannot flat (%T) %v", value, value))
			}
		}}
	return eval
}

func getSequenceable(v interface{}) (melrose.Sequenceable, bool) {
	if s, ok := v.(melrose.Sequenceable); ok {
		return s, ok
	}
	if f, ok := v.(FunctionResult); ok {
		if f.Notification != nil {
			notify.Print(f.Notification)
		}
		return getSequenceable(f.Result)
	}
	return nil, false
}
