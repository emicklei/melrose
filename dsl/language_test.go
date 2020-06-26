package dsl

import (
	"testing"

	"github.com/emicklei/melrose"
)

func TestNote(t *testing.T) {
	r := eval(t, "note('c')")
	checkStorex(t, r, "note('C')")
	checkStorex(t, r.(melrose.Sequenceable).S(),
		"sequence('C')")
}

func TestChord(t *testing.T) {
	r := eval(t, "chord('C#/m')")
	checkStorex(t, r, "chord('C♯/m')")
	checkStorex(t, r.(melrose.Sequenceable).S(),
		"sequence('(C♯ E A♭)')")
}

func TestSequence(t *testing.T) {
	r := eval(t, "sequence('c (d e g) =')")
	checkStorex(t, r, "sequence('C (D E G) =')")
}

func TestProgression(t *testing.T) {
	r := eval(t, "progression('c/m (d7 e g) =')")
	checkStorex(t, r, "progression('C/m (D7 E G) =')")
	checkStorex(t, r.(melrose.Sequenceable).S(),
		"sequence('(C E♭ G) (D7 G♭7 A7 E A♭ B G B D5) =')")
}

func TestScale(t *testing.T) {
	r := eval(t, "scale(2,'16e2')")
	checkStorex(t, r, "scale(2,'16E2')")
	checkStorex(t, r.(melrose.Sequenceable).S(),
		"sequence('16E2 16G♭2 16A♭2 16A2 16B2 16D♭3 16E♭3 16E3 16G♭3 16A♭3 16A3 16B3 16D♭ 16E♭')")
}

func TestPitch_Scale(t *testing.T) {
	r := eval(t, "pitch(1,scale(2,'16e2'))")
	checkStorex(t, r, "pitch(1,scale(2,'16E2'))")
	checkStorex(t, r.(melrose.Sequenceable).S(),
		"sequence('16F2 16G2 16A2 16B♭2 16C3 16D3 16E3 16F3 16G3 16A3 16B♭3 16C 16D 16E')")
}

func TestPitch_Progression(t *testing.T) {
	r := eval(t, "pitch(1,progression('c/m (d7 e g) ='))")
	checkStorex(t, r, "pitch(1,progression('C/m (D7 E G) ='))")
	checkStorex(t, r.(melrose.Sequenceable).S(),
		"sequence('(D♭ E A♭) (E♭7 G7 B♭7 F A C5 A♭ C5 E♭5) =')")
}
