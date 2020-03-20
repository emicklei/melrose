package main

import "github.com/emicklei/melrose"

func evalFunctions() map[string]interface{} {
	eval := map[string]interface{}{}
	eval["chord"] = func(note string) melrose.Chord {
		n, err := melrose.ParseNote(note)
		if err != nil {
			printError(err)
			return melrose.N("C").Chord()
		}
		return n.Chord()
	}

	eval["join"] = func(playables ...interface{}) interface{} {
		joined := melrose.S("")
		for _, p := range playables {
			if s, ok := p.(melrose.Sequenceable); ok {
				joined = joined.SequenceJoin(s.S())
			}
		}
		return joined
	}

	eval["bpm"] = func(i int) int {
		piano.BeatsPerMinute(float64(i))
		return i
	}

	eval["seq"] = func(s string) melrose.Sequence {
		n, err := melrose.ParseSequence(s)
		if err != nil {
			printError(err)
			return melrose.N("C").S()
		}
		return n
	}

	eval["note"] = func(s string) melrose.Note {
		n, err := melrose.ParseNote(s)
		if err != nil {
			printError(err)
			return melrose.N("C")
		}
		return n
	}

	eval["play"] = func(playables ...interface{}) interface{} {
		for _, p := range playables {
			if s, ok := p.(melrose.Sequenceable); ok {
				piano.Play(s.S())
			}
		}
		return nil
	}
	return eval
}
