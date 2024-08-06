package dsl

import (
	"os"
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
	checkStorex(t, r, "chord('C#/m')")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('(C# E A_)')")
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
	r := eval(t, "chordsequence('c/m (d7 e g) =')")
	checkStorex(t, r, "chordsequence('C/m (D7 E G) =')")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('(C E_ G) (D7 G_7 A7 E A_ B G B D5) =')")
}

func TestChordSequence_Invalid(t *testing.T) {
	mustError(t, "chordsequence('k')", "illegal note")
}

func TestScale(t *testing.T) {
	r := eval(t, "scale('16e2')")
	checkStorex(t, r, "scale('16E2')")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('16E2 16G_2 16A_2 16A2 16B2 16D_3 16E_3')")
}

func TestTranspose_ChordSequence(t *testing.T) {
	r := eval(t, "transpose(1,chordsequence('c/m (d7 e g) ='))")
	checkStorex(t, r, "transpose(1,chordsequence('C/m (D7 E G) ='))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('(D_ E A_) (E_7 G7 B_7 F A C5 A_ C5 E_5) =')")
}

func TestTrack(t *testing.T) {
	r := eval(t, "track('test',1,onbar(1,note('c')))")
	checkStorex(t, r, "track('test',1,onbar(1,note('C')))")
}

func TestChannelSelector(t *testing.T) {
	r := eval(t, "channel(1,note('f'))")
	checkStorex(t, r, "channel(1,note('F'))")
}

func TestDeviceSelector(t *testing.T) {
	r := eval(t, "device(1,note('f'))")
	checkStorex(t, r, "device(1,note('F'))")
}

func TestBars(t *testing.T) {
	r := eval(t, "bars(sequence('a b c d'))")
	if got, want := r.(int), 1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestBars_Arithmetic(t *testing.T) {
	r := eval(t, "1+bars(sequence('a b c d'))")
	if got, want := r.(core.HasValue).Value(), 2; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNotemapNumbers(t *testing.T) {
	r := eval(t, "notemap('2 4 11',note('c'))")
	checkStorex(t, r, "notemap('2 4 11',note('C'))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('= C = C = = = = = = C')")
}

func TestTwoBarsNote(t *testing.T) {
	r := eval(t, "stretch(2,note('1c'))")
	checkStorex(t, r, "stretch(2,note('1C'))")
	s := r.(core.Sequenceable).S()
	n := s.At(0)[0]
	if got, want := n.DurationFactor(), float32(2.0); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestStretchChord(t *testing.T) {
	r := eval(t, "stretch(2,chord('1c'))")
	checkStorex(t, r, "stretch(2,chord('1C'))")
	s := r.(core.Sequenceable).S()
	n := s.At(0)[0]
	if got, want := n.DurationFactor(), float32(2.0); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestChordSequenceWithSustain(t *testing.T) {
	r := eval(t, "chordsequence('> 1g/m/2 ^ 1d5 <')")
	checkStorex(t, r, "chordsequence('> 1G/m/2 ^ 1D5 <')")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('> (1D5 1G5 1B_5) ^ (1D5 1G_5 1A5) <')")
}

func TestNotemapDots(t *testing.T) {
	r := eval(t, "notemap('.!.!',note('16c'))")
	checkStorex(t, r, "notemap('.!.!',note('16C'))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('16= 16C 16= 16C')")
}

func TestNotemapSequence(t *testing.T) {
	r := eval(t, "notemap('.!.!',sequence('16c A'))")
	checkStorex(t, r, "notemap('.!.!',sequence('16C A'))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('16= 16C 16= 16C')")
}

func TestSequenceMapFromStringIterator(t *testing.T) {
	r := eval(t, `ar = iterator('1','2')
sm = resequence(ar,sequence('c d'))`)
	checkStorex(t, r, "resequence(ar,sequence('C D'))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('C')")
}

func TestDeviceOnChannelOnNote(t *testing.T) {
	r := eval(t, `d = device(1,channel(2,note('e')))`)
	checkStorex(t, r.(core.Sequenceable).S(), "sequence('E')")
}

func TestDynamicMapWithTwoSequenceables(t *testing.T) {
	r := eval(t, `s = sequence('c e g')
t = note('b')
dm = dynamicmap(' 1:+ , 2:-- ,3:o,4:o',s,t)`)
	checkStorex(t, r, "dynamicmap('1:+,2:--,3:o,4:o',s,t)")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('C+ E-- G B')")
}

func TestDynamic_String(t *testing.T) {
	r := eval(t, `d = dynamic('+',note('c'))`)
	checkStorex(t, r, "dynamic('+',note('C'))")
}

func TestDynamic_Number(t *testing.T) {
	r := eval(t, `d = dynamic(55,note('c'))`)
	checkStorex(t, r, "dynamic(55,note('C'))")
}

func TestDynamic_Var(t *testing.T) {
	r := eval(t, `v = 55
d = dynamic(v,note('c'))`)
	checkStorex(t, r, "dynamic(v,note('C'))")
}

func TestProcessLanguageTest(t *testing.T) {
	src, _ := os.ReadFile("language_test.mel")
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("panic recovered:", err)
		}
	}()
	eval(t, string(src))
}

func TestVelocityMap(t *testing.T) {
	r := eval(t, "velocitymap('1:30,2:60,2:0',sequence('c g'))")
	checkStorex(t, r, "velocitymap('1:30,2:60,2:0',sequence('C G'))")
	checkStorex(t, r.(core.Sequenceable).S(),
		"sequence('C--- G+ G')")
}

func TestValueOfVar(t *testing.T) {
	r := eval(t, `
i = iterator(1)
v = value(i)`)
	checkStorex(t, r, "value(i)")
	if got, want := core.ValueOf(r), 1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestDeviceOnTrack(t *testing.T) {
	t.Skip() // TODO
	r := eval(t, `
s = sequence('a b')
dt = device(1,track('title',4, onbar(1,s)))`)
	checkStorex(t, r, "value(i)")
}

func TestIteratorIndex(t *testing.T) {
	r := eval(t, `it = iterator(1,2,3)
idx = it.Index()`)
	checkStorex(t, r, "it.Index()")
}
