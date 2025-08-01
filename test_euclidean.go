package main

import (
	"fmt"
)

// https://github.com/computermusicdesign/euclidean-rhythm/blob/master/max-example/euclidSimple.js
func euclideanRhythm(steps, pulses, rotation int) []bool {
	rhythm := make([]bool, 0, steps)
	bucket := 0

	for i := 0; i < steps; i++ {
		bucket += pulses
		if bucket >= steps {
			bucket -= steps
			rhythm = append(rhythm, true)
		} else {
			rhythm = append(rhythm, false)
		}
	}

	return rotatedRhythm(rhythm, rotation+1)
}

func rotatedRhythm(input []bool, rotate int) []bool {
	output := make([]bool, len(input))
	val := len(input) - rotate
	for i := 0; i < len(input); i++ {
		j := (i + val) % len(input)
		if j < 0 {
			j *= -1
		}
		output[i] = input[j]
	}
	return output
}

func printRhythm(name string, steps, pulses, rotation int) {
	result := euclideanRhythm(steps, pulses, rotation)
	fmt.Printf("%s: %v\n", name, result)

	// Also print as pattern
	pattern := ""
	for _, beat := range result {
		if beat {
			pattern += "!"
		} else {
			pattern += "."
		}
	}
	fmt.Printf("%s pattern: %s\n", name, pattern)
	fmt.Println()
}

func main() {
	printRhythm("4 steps 1 pulse", 4, 1, 0)
	printRhythm("8 steps 3 pulses", 8, 3, 0)
	printRhythm("16 steps 5 pulses", 16, 5, 0)
	printRhythm("8 steps 5 pulses", 8, 5, 0)
	printRhythm("4 steps 2 pulses rotation 1", 4, 2, 1)
	printRhythm("8 steps 3 pulses rotation 2", 8, 3, 2)
}
