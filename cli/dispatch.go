package main

import (
	"fmt"
	"log"
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
	}
	setupPiano(env)
	for k, v := range memory {
		env[k] = v
	}
	program, err := expr.Compile(entry, expr.Env(env))
	if err != nil {
		return nil, err
	}
	return expr.Run(program, env)
}

func evalNote(s string) melrose.Note {
	n, err := melrose.ParseNote(s)
	if err != nil {
		printError(err)
		return melrose.N("C")
	}
	return n
}

func evalPlay(p interface{}) error {
	if n, ok := p.(melrose.Note); ok {
		piano.Play(n.S())
		return nil
	}
	if s, ok := p.(melrose.Sequence); ok {
		piano.Play(s)
		return nil
	}
	if c, ok := p.(melrose.Chord); ok {
		piano.Play(c.S())
		return nil
	}
	return nil
}

func setupPiano(env map[string]interface{}) {
	piano := map[string]melrose.Note{}
	env["piano"] = piano
	for octave := 3; octave < 6; octave++ {
		for _, each := range strings.Fields("C D E F G A B") {
			key := fmt.Sprintf("%s%d", each, octave)
			note, err := melrose.ParseNote(key)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(key, note)
				piano[key] = note
			}
		}
	}
}
