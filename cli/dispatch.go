package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/emicklei/melrose"
)

var assignmentRegex = regexp.MustCompile(`^[a-z]+\[[0-9]+\]=.*$`)

func dispatch(entry string) error {
	if len(entry) == 0 {
		fmt.Println()
		return nil
	}
	if value, ok := memory[entry]; ok {
		printValue(value)
		return nil
	}
	// is assignment?
	// TODO not correct
	if strings.Contains(entry, "=") {
		parts := strings.Split(entry, "=")
		variable := strings.TrimSpace(parts[0])
		expression := parts[1]
		r, err := eval(expression)
		if err != nil {
			return err
		}
		memory[variable] = r
		printValue(r)
		return nil
	}
	// evaluate and print
	r, err := eval(entry)
	if err != nil {
		return err
	}
	printValue(r)
	return nil
}

func printValue(v interface{}) {
	if v == nil {
		fmt.Println()
		return
	}
	fmt.Printf("(%T) %v\n", v, v)
}

func eval(entry string) (interface{}, error) {
	env := map[string]interface{}{
		"note": evalNote,
		"play": evalPlay,
		"seq":  evalSeq,
		"bpm":  evalBPM,
		"join": evalJoin,
	}
	env["piano"] = pianoNotes
	for k, v := range memory {
		env[k] = v
	}
	program, err := expr.Compile(entry, expr.Env(env))
	if err != nil {
		return nil, err
	}
	return expr.Run(program, env)
}

func evalJoin(playables ...interface{}) interface{} {
	joined := melrose.S("")
	for _, p := range playables {
		if s, ok := p.(melrose.Sequenceable); ok {
			joined = joined.SequenceJoin(s.S())
		}
	}
	return joined
}

func evalBPM(i int) int {
	piano.BeatsPerMinute(float64(i))
	return i
}

func evalSeq(s string) melrose.Sequence {
	n, err := melrose.ParseSequence(s)
	if err != nil {
		printError(err)
		return melrose.N("C").S()
	}
	return n
}

func evalNote(s string) melrose.Note {
	n, err := melrose.ParseNote(s)
	if err != nil {
		printError(err)
		return melrose.N("C")
	}
	return n
}

func evalPlay(playables ...interface{}) interface{} {
	for _, p := range playables {
		if s, ok := p.(melrose.Sequenceable); ok {
			piano.Play(s.S())
		}
	}
	return nil
}
