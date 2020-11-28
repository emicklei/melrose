package dsl

import (
	"io/ioutil"
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestNote(t *testing.T) {
	r := eval(t, "note('c')")
	checkStorex(t, r, "note('C')")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('C')")
}

func TestNote_Invalid(t *testing.T) {
	mustError(t, "note('k')", "illegal note")
}

func TestChord(t *testing.T) {
	r := eval(t, "chord('C#/m')")
	checkStorex(t, r, "chord('C♯/m')")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('(C♯ E A♭)')")
}

func TestChord_Invalid(t *testing.T) {
	mustError(t, "chord('k')", "illegal note")
}

func TestSequence(t *testing.T) {
	r := eval(t, "sequence('c (d e g) =')")
	checkStorex(t, r, "sequence('C (D E G) =')")
}

func TestSequence_Invalid(t *testing.T) {
	mustError(t, "sequence('k')", "illegal note")
}

func TestProgression(t *testing.T) {
	r := eval(t, "progression('c/m (d7 e g) =')")
	checkStorex(t, r, "progression('C/m (D7 E G) =')")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('(C E♭ G) (D7 G♭7 A7 E A♭ B G B D5) =')")
}

func TestProgression_Invalid(t *testing.T) {
	mustError(t, "progression('k')", "illegal note")
}

func TestScale(t *testing.T) {
	r := eval(t, "scale(2,'16e2')")
	checkStorex(t, r, "scale(2,'16E2')")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('16E2 16G♭2 16A♭2 16A2 16B2 16D♭3 16E♭3 16E3 16G♭3 16A♭3 16A3 16B3 16D♭ 16E♭')")
}

func TestPitch_Scale(t *testing.T) {
	r := eval(t, "pitch(1,scale(2,'16e2'))")
	checkStorex(t, r, "pitch(1,scale(2,'16E2'))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('16F2 16G2 16A2 16B♭2 16C3 16D3 16E3 16F3 16G3 16A3 16B♭3 16C 16D 16E')")
}

func TestPitch_Progression(t *testing.T) {
	r := eval(t, "pitch(1,progression('c/m (d7 e g) ='))")
	checkStorex(t, r, "pitch(1,progression('C/m (D7 E G) ='))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('(D♭ E A♭) (E♭7 G7 B♭7 F A C5 A♭ C5 E♭5) =')")
}

func TestTrack(t *testing.T) {
	r := eval(t, "track('test',1,onbar(1,note('c')))")
	checkStorex(t, r, "track('test',1,onbar(1,note('C')))")
}

func TestChannelSelector(t *testing.T) {
	r := eval(t, "channel(1,note('e'),note('f'))")
	checkStorex(t, r, "channel(1,note('E'),note('F'))")
}

func TestDeviceSelector(t *testing.T) {
	r := eval(t, "device(1,note('e'),note('f'))")
	checkStorex(t, r, "device(1,note('E'),note('F'))")
}

func TestBars(t *testing.T) {
	r := eval(t, "bars(sequence('a b c d'))")
	if got, want := r.(int), 1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestBars_Arithmetic(t *testing.T) {
	r := eval(t, "1+bars(sequence('a b c d'))")
	if got, want := r.(int), 2; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNotemapNumbers(t *testing.T) {
	r := eval(t, "notemap('2 4 11',note('c'))")
	checkStorex(t, r, "notemap('2 4 11',note('C'))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('= C = C = = = = = = C')")
}

func TestNotemapDots(t *testing.T) {
	r := eval(t, "notemap('.!.!',note('16c'))")
	checkStorex(t, r, "notemap('.!.!',note('16C'))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('16= 16C 16= 16C')")
}

func TestSequenceMapFromStringIterator(t *testing.T) {
	r := eval(t, `ar = iterator('1','2')
sm = sequencemap(ar,sequence('c d'))`)
	checkStorex(t, r, "sequencemap(ar,sequence('C D'))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('C')")
}

func TestDynamicMapWithTwoSequenceables(t *testing.T) {
	r := eval(t, `s = sequence('c e g')
t = note('b')
dm = dynamicmap(' 1:+ , 2:-- ,3:o,4:o',s,t)`)
	checkStorex(t, r, "dynamicmap('1:+,2:--,3:o,4:o',s,t)")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('C+ E-- G B')")
}

func TestProcessLanguageTest(t *testing.T) {
	src, _ := ioutil.ReadFile("language_test.mel")
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("panic recovered:", err)
		}
	}()
	eval(t, string(src))
}
