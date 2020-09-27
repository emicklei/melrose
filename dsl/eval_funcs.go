package dsl

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/midi/file"

	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/melrose/op"
)

// Syntax tells what language version this package is supporting.
const Syntax = "1.0-rc" // major,minor

func IsCompatibleSyntax(s string) bool {
	if len(s) == 0 {
		// ignore syntax ; you are on your own
		return true
	}
	mm := strings.Split(Syntax, ".")
	ss := strings.Split(s, ".")
	return mm[0] == ss[0] && ss[1] <= mm[1]
}

type Function struct {
	Title         string
	Description   string
	Prefix        string // for autocomplete
	Alias         string // short notation
	Template      string // for autocomplete in VSC
	Samples       string // for doc generation
	ControlsAudio bool
	Tags          string // space separated
	IsCore        bool   // creates a core musical object
	IsComposer    bool   // can decorate a musical object or other decorations
	Func          interface{}
}

func EvalFunctions(ctx core.Context) map[string]Function {
	eval := map[string]Function{}

	eval["fraction"] = Function{
		Title: "Duration fraction operator",
		Description: `Creates a new object for which the fraction of duration of all notes are changed.
The first parameter controls the fraction of the note, e.g. 1=whole, 0.5 or 2 = half, 0.25 or 4 = quarter, 0.125 or 8 = eight, 0.0625 or 16 = sixteenth.
`,
		Prefix:     "fra",
		IsComposer: true,
		Template:   `fraction(${1:object},${2:object})`,
		Samples:    `fraction(8,sequence('e f')) // => ⅛E ⅛F , shorten the notes from quarter to eigth`,
		Func: func(param float64, playables ...interface{}) interface{} {
			if err := op.CheckFraction(param); err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Print(notify.Warningf("cannot fraction (%T) %v", p, p))
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return op.NewFraction(param, joined)
		}}

	eval["dynamic"] = Function{
		Title: "Dynamic operator",
		Description: `Creates a new modified musical object for which the dynamics of all notes are changed.
	The first parameter controls the emphasis the note, e.g. + (mezzoforte,mf), -- (piano,p).
	`,
		Prefix:     "dy",
		IsComposer: true,
		Template:   `dynamic(${1:emphasis},${2:object})`,
		Samples:    `dynamic('++',sequence('E F')) // => E++ F++`,
		Func: func(emphasis string, playables ...interface{}) interface{} {
			if err := op.CheckDynamic(emphasis); err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Print(notify.Warningf("cannot dynamic (%T) %v", p, p))
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return op.Dynamic{Target: joined, Emphasis: emphasis}
		}}

	eval["progression"] = Function{
		Title:       "Progress creator",
		Description: `create a Chord progression using this <a href="/melrose/notations.html#progression-not">format</a>`,
		Prefix:      "pro",
		IsCore:      true,
		Template:    `progression('${1:chords}')`,
		Samples: `progression('E F') // => (E A♭ B) (F A C5)
progression('(C D)') // => (C E G D G♭ A)`,
		Func: func(chords string) interface{} {
			p, err := core.ParseProgression(chords)
			if err != nil {
				return notify.Panic(err)
			}
			return p
		}}

	eval["joinmap"] = Function{
		Title:       "Join Map creator",
		Description: "creates a new join by mapping elements based on an index (1-based)",
		Prefix:      "joinm",
		IsComposer:  true,
		Template:    `joinmap('${1:indices}',${2:join})`,
		Func: func(indices string, join interface{}) interface{} { // allow multiple seq?
			v := getValueable(join)
			vNow := v.Value()
			if _, ok := vNow.(op.Join); !ok {
				return notify.Panic(fmt.Errorf("cannot joinmap (%T) %v", join, join))
			}
			return op.NewJoinMapper(v, indices)
		}}

	eval["bars"] = Function{
		Prefix:      "ba",
		Description: "compute the number of bars that is taken when playing a musical object",
		IsComposer:  true,
		Template:    `bars(${1:object})`,
		Func: func(seq interface{}) interface{} {
			s, ok := getSequenceable(seq)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot compute how many bars for (%T) %v", seq, seq))
			}
			// TODO handle loop
			biab := ctx.Control().BIAB()
			return int(math.Round((s.S().DurationFactor() * 4) / float64(biab)))
		}}

	eval["beats"] = Function{
		Prefix:      "be",
		Description: "compute the number of beats that is taken when playing a musical object",
		IsComposer:  true,
		Template:    `beats(${1:object})`,
		Func: func(seq interface{}) interface{} {
			s, ok := getSequenceable(seq)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot compute how many beats for (%T) %v", seq, seq))
			}
			return len(s.S().Notes)
		}}

	eval["track"] = Function{
		Title:       "Track creator",
		Description: "create a named track for a given MIDI channel with a musical object",
		Prefix:      "tr",
		Template:    `track('${1:title}',${2:channel})`,
		Samples:     `track("lullaby",1,sequence('c d e')) // => a new track on MIDI channel 1`,
		Func: func(title string, channel int, playables ...interface{}) interface{} {
			if len(title) == 0 {
				return notify.Panic(fmt.Errorf("cannot have a track without title"))
			}
			if channel < 1 || channel > 15 {
				return notify.Panic(fmt.Errorf("MIDI channel must be in [1..15]"))
			}
			tr := core.NewTrack(title, channel)
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					return notify.Panic(fmt.Errorf("cannot compose track with (%T) %v", p, p))
				} else {
					tr.Add(s)
				}
			}
			return tr
		}}

	eval["multi"] = Function{
		Title:         "Multi track creator",
		Description:   "create a multi-track object from zero or more tracks",
		Prefix:        "mtr",
		Template:      `multi(${1:track})`,
		Samples:       `multi(track1,track2,track3) // one or more tracks in one multi-track object`,
		ControlsAudio: true,
		Func: func(varOrTrack ...interface{}) interface{} {
			tracks := []core.Valueable{}
			for _, each := range varOrTrack {
				tracks = append(tracks, getValueable(each))
			}
			return core.MultiTrack{Tracks: tracks}
		}}

	eval["midi"] = Function{
		Title: "Note creator",
		Description: `create a Note from MIDI information and is typically used for drum sets.
The first parameter is the duration and must be one of {0.0625,0.125,0.25,0.5,1,2,4,8,16}.
A duration of 0.25 or 4 means create a quarter note.
Second parameter is the MIDI number and must be one of [0..127].
The third parameter is the velocity (~ loudness) and must be one of [0..127]`,
		Prefix:   "mid",
		Alias:    "M",
		Template: `midi(${1:number},${2:number},${3:number})`,
		Samples: `midi(0.25,52,80) // => E3+
midi(16,36,70) // => 16C2 (kick)`,
		IsCore: true,
		Func: func(dur, nr, velocity interface{}) interface{} {
			durVal := getValueable(dur)
			nrVal := getValueable(nr)
			velVal := getValueable(velocity)
			return core.NewMIDI(durVal, nrVal, velVal)
		}}

	eval["print"] = Function{
		Title:       "Printer creator",
		Description: "prints the musical object when evaluated (play,go,loop)",
		Func: func(m interface{}) interface{} {
			s, ok := getSequenceable(getValue(m))
			if !ok {
				return notify.Panic(fmt.Errorf("cannot print (%T) %v", m, m))
			}
			return core.Watch{Context: ctx, Target: s}
		}}

	eval["chord"] = Function{
		Title:       "Chord creator",
		Description: `create a Chord from its string <a href="/melrose/notations.html#chord-not">notation</a>`,
		Prefix:      "cho",
		Alias:       "C",
		Template:    `chord('${1:note}')`,
		Samples: `chord('C#5/m/1')
chord('G/M/2')`,
		IsCore: true,
		Func: func(chord string) interface{} {
			c, err := core.ParseChord(chord)
			if err != nil {
				return notify.Panic(err)
			}
			return c
		}}

	eval["octavemap"] = Function{
		Title:       "Octave Map operator",
		Description: "create a sequence with notes for which the order and the octaves are changed",
		Prefix:      "octavem",
		Template:    `octavemap('${1:int2int}',${2:object})`,
		IsComposer:  true,
		Samples:     `octavemap('1:-1,2:0,3:1',chord('C')) // => (C3 E G5)`,
		Func: func(indices string, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot octavemap (%T) %v", m, m))
			}
			return op.NewOctaveMapper(s, indices)
		}}

	eval["pitch"] = Function{
		Title:       "Pitch operator",
		Description: "change the pitch with a delta of semitones",
		Prefix:      "pit",
		Alias:       "Pi",
		Template:    `pitch(${1:semitones},${2:sequenceable})`,
		Samples: `pitch(-1,sequence('c d e'))
p = interval(-4,4,1)
pitch(p,note('c'))`,
		IsComposer: true,
		Func: func(semitones, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot pitch (%T) %v", m, m))
			}
			return op.Pitch{Target: s, Semitones: getValueable(semitones)}
		}}

	eval["reverse"] = Function{
		Title:       "Reverse operator",
		Description: "reverse the (groups of) notes in a sequence",
		Prefix:      "rev",
		Alias:       "Rv",
		Template:    `reverse(${1:sequenceable})`,
		Samples:     `reverse(chord('A'))`,
		IsComposer:  true,
		Func: func(m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot reverse (%T) %v", m, m))
			}
			return op.Reverse{Target: s}
		}}

	eval["repeat"] = Function{
		Title:       "Repeat operator",
		Description: "repeat the musical object a number of times",
		Prefix:      "rep",
		Alias:       "Rp",
		Template:    `repeat(${1:times},${2:sequenceable})`,
		Samples:     `repeat(4,sequence('C D E'))`,
		IsComposer:  true,
		Func: func(howMany interface{}, playables ...interface{}) interface{} {
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					return notify.Panic(fmt.Errorf("cannot repeat (%T) %v", p, p))
				} else {
					joined = append(joined, s)
				}
			}
			return op.Repeat{Target: joined, Times: getValueable(howMany)}
		}}

	registerFunction(eval, "join", Function{
		Title:       "Join operator",
		Description: "joins two or more musical objects as one",
		Prefix:      "joi",
		Alias:       "J",
		Template:    `join(${1:first},${2:second})`,
		Samples: `a = chord('A')
b = sequence('(C E G)')
ab = join(a,b)`,
		IsComposer: true,
		Func: func(playables ...interface{}) interface{} {
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					return notify.Panic(fmt.Errorf("cannot join (%T) %v", p, p))
				} else {
					joined = append(joined, s)
				}
			}
			return op.Join{Target: joined}
		}})

	eval["bpm"] = Function{
		Title:         "Beats Per Minute",
		Description:   "set the Beats Per Minute (BPM) [1..300]; default is 120",
		ControlsAudio: true,
		Prefix:        "bpm",
		Template:      `bpm(${1:beats-per-minute})`,
		Samples:       `bpm(90)`,
		Func: func(f float64) interface{} {
			if f < 1 || f > 300 {
				return notify.Panic(fmt.Errorf("invalid beats-per-minute [1..300], %f = ", f))
			}
			ctx.Control().SetBPM(f)
			return nil
		}}

	eval["biab"] = Function{
		Title:         "Beats in a Bar",
		Description:   "set the Beats in a Bar [1..6]; default is 4",
		ControlsAudio: true,
		Prefix:        "biab",
		Template:      `biab(${1:beats-in-a-bar})`,
		Samples:       `biab(4)`,
		Func: func(i int) interface{} {
			if i < 1 || i > 6 {
				return notify.Panic(fmt.Errorf("invalid beats-in-a-bar [1..6], %d = ", i))
			}
			ctx.Control().SetBIAB(i)
			return nil
		}}

	registerFunction(eval, "import", Function{
		Title:         "Import script",
		Description:   "evaluate all the statements from another file",
		ControlsAudio: false,
		Template:      `import(${1:filename})`,
		Samples:       `import('drumpatterns.mel')`,
		Func: func(f string) interface{} {
			err := ImportProgram(ctx, f)
			if err != nil {
				return notify.Panic(fmt.Errorf("failed to import [%s], %v", f, err))
			}
			return nil
		},
	})

	eval["sequence"] = Function{
		Title:       "Sequence creator",
		Description: `create a Sequence using this <a href="/melrose/notations.html#sequence-not">format</a>`,
		Prefix:      "seq",
		Alias:       "S",
		Template:    `sequence('${1:space-separated-notes}')`,
		Samples: `sequence('C D E')
sequence('(8C D E)') => (⅛C ⅛D ⅛E)
sequence('c (d e f) a =')`,
		IsCore: true,
		Func: func(s string) interface{} {
			sq, err := core.ParseSequence(s)
			if err != nil {
				return notify.Panic(err)
			}
			return sq
		}}

	eval["note"] = Function{
		Title:       "Note creator",
		Description: `create a Note using this <a href="/melrose/notations.html#note-not">format</a>`,
		Prefix:      "no",
		Alias:       "N",
		Template:    `note('${1:letter}')`,
		Samples: `note('e')
note('2.e#--')`,
		IsCore: true,
		Func: func(s string) interface{} {
			n, err := core.ParseNote(s)
			if err != nil {
				return notify.Panic(err)
			}
			return n
		}}

	eval["scale"] = Function{
		Title:       "Scale creator",
		Description: `create a Scale using this <a href="/melrose/notations.html#scale-not">format</a>`,
		Prefix:      "sc",
		Template:    `scale(${1:octaves},'${2:note}')`,
		IsCore:      true,
		Samples:     `scale(1,'E/m') // => E F G A B C5 D5`,
		Func: func(octaves int, s string) interface{} {
			if octaves < 1 {
				return notify.Panic(fmt.Errorf("octaves must be >= 1%v", octaves))
			}
			sc, err := core.NewScale(octaves, s)
			if err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			return sc
		}}

	eval["at"] = Function{
		Title:       "Index getter",
		Description: "create an index getter (1-based) to select a musical object",
		Prefix:      "at",
		Template:    `at(${1:index},${2:object})`,
		Samples:     `at(1,scale('E/m')) // => E`,
		Func: func(index interface{}, object interface{}) interface{} {
			indexVal := getValueable(index)
			objectSeq, ok := getSequenceable(object)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot index (%T) %v", object, object))
			}
			return op.NewAtIndex(indexVal, objectSeq)
		}}

	eval["onbar"] = Function{
		Title:       "Track modifier",
		Description: "puts a musical object on a track to start at a specific bar",
		Prefix:      "onbar",
		Template:    `onbar(${1:bar},${2:object})`,
		Samples:     `tr = track("solo",2, onbar(1,soloSequence)) // 2 = channel`,
		Func: func(bar interface{}, seq interface{}) interface{} {
			s, ok := getSequenceable(seq)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot put on track (%T) %v", seq, seq))
			}
			return core.NewSequenceOnTrack(getValueable(bar), getValueable(0), s)
		}}

	eval["random"] = Function{
		Title:       "Random generator",
		Description: "create a random integer generator. Use next() to generate a new integer",
		Prefix:      "ra",
		Template:    `random(${1:from},${2:to})`,
		Samples: `num = random(1,10)
next(num)`,
		Func: func(from interface{}, to interface{}) interface{} {
			fromVal := getValueable(from)
			toVal := getValueable(to)
			return op.NewRandomInteger(fromVal, toVal)
		}}

	eval["play"] = Function{
		Title:         "Play musical objects in the foreground",
		Description:   "play all musical objects",
		ControlsAudio: true,
		Prefix:        "pla",
		Template:      `play(${1:sequenceable})`,
		Samples:       `play(s1,s2,s3) // play s3 after s2 after s1`,
		Func: func(playables ...interface{}) interface{} {
			moment := time.Now()
			for _, p := range playables {
				if s, ok := getSequenceable(getValue(p)); ok { // unwrap var
					moment = ctx.Device().Play(s, ctx.Control().BPM(), moment)
				} else {
					notify.Print(notify.Warningf("cannot play (%T) %v", p, p))
				}
			}
			// wait until the play is completed before allowing a new one to.
			// add a bit of time to allow the previous play to finish all its notes.
			// time.Sleep(moment.Sub(time.Now()) + (50 * time.Millisecond))
			return nil
		}}

	eval["go"] = Function{
		Title:         "Play musical objects in the background",
		Description:   "play all musical objects together in the background (do not wait for completion)",
		ControlsAudio: true,
		Prefix:        "go",
		Template:      `go(${1:sequenceable})`,
		Samples:       `go(s1,s1,s3) // play s1 and s2 and s3 simultaneously`,
		Func: func(playables ...interface{}) interface{} {
			moment := time.Now()
			for _, p := range playables {
				if s, ok := getSequenceable(getValue(p)); ok { // unwrap var
					ctx.Device().Play(s, ctx.Control().BPM(), moment)
				} else {
					notify.Print(notify.Warningf("cannot go (%T) %v", p, p))
				}
			}
			return nil
		}}

	eval["ungroup"] = Function{
		Title:       "Ungroup operator",
		Description: "undo any grouping of notes from one or more musical objects",
		Prefix:      "ung",
		Template:    `ungroup(${1:sequenceable})`,
		IsComposer:  true,
		Samples: `ungroup(chord('E')) // => E G B
ungroup(sequence('(C D)'),note('E')) // => C D E`,
		Func: func(playables ...interface{}) interface{} {
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Print(notify.Warningf("cannot ungroup (%T) %v", p, p))
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return op.Serial{Target: joined}
		}}

	eval["octave"] = Function{
		Title:       "Octave operator",
		Description: "change the pitch of notes by steps of 12 semitones for one or more musical objects",
		Prefix:      "oct",
		Template:    `octave(${1:offset},${2:sequenceable})`,
		IsComposer:  true,
		Samples:     `octave(1,sequence('C D')) // => C5 D5`,
		Func: func(scalarOrVar interface{}, playables ...interface{}) interface{} {
			list := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Print(notify.Warningf("cannot octave (%T) %v", p, p))
					return nil
				} else {
					list = append(list, s)
				}
			}
			return op.Octave{Target: list, Offset: core.ToValueable(scalarOrVar)}
		}}

	eval["record"] = Function{
		Title:         "Recording creator",
		Description:   "create a recorded sequence of notes from the current MIDI input device",
		ControlsAudio: true,
		Prefix:        "rec",
		Template:      `record()`,
		Samples: `r = record() // record notes played on the current input device and stop recording after 5 seconds
s = r.S() // returns the sequence of notes from the recording`,
		Func: func() interface{} {
			seq, err := ctx.Device().Record(ctx)
			if err != nil {
				return notify.Panic(err)
			}
			return seq
		}}

	eval["undynamic"] = Function{
		Title:       "Undo dynamic operator",
		Description: "set the dymamic to normal for all notes in a musical object",
		Prefix:      "und",
		Template:    `undynamic(${1:sequenceable})`,
		IsComposer:  true,
		Samples:     `undynamic('A+ B++ C-- D-') // =>  A B C D`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				return notify.Panic(fmt.Errorf("cannot undynamic (%T) %v", value, value))
			} else {
				return op.Undynamic{Target: s}
			}
		}}

	registerFunction(eval, "iterator", Function{
		Title:       "Iterator creator",
		Description: "iterator that has an array of constant values and evaluates to one. Use next() to increase and rotate the value.",
		Prefix:      "it",
		Alias:       "I",
		Template:    `iterator(${1:array-element})`,
		Samples: `i = iterator(1,3,5,7,9)
		p = pitch(i,note('c'))
		lp = loop(p,next(i))
		`,
		Func: func(values ...interface{}) *core.Iterator {
			return &core.Iterator{
				Target: values,
			}
		}})

	eval["group"] = Function{
		Title:       "Group operator",
		Description: "create a new sequence in which all notes of a musical object are grouped",
		Prefix:      "par",
		Alias:       "Pa",
		Template:    `group(${1:sequenceable})`,
		Samples:     `group(sequence('C D E')) // => (C D E)`,
		IsComposer:  true,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				return notify.Panic(fmt.Errorf("cannot group (%T) %v", value, value))
			} else {
				return op.Parallel{Target: s}
			}
		}}
	// BEGIN Loop and control
	eval["loop"] = Function{
		Title:         "Loop creator",
		Description:   "create a new loop from one or more musical objects; must be assigned to a variable",
		ControlsAudio: true,
		Prefix:        "loo",
		Alias:         "L",
		Template:      `lp_${1:object} = loop(${1:object})`,
		Samples: `cb = sequence('C D E F G A B')
lp_cb = loop(cb,reverse(cb))`,
		Func: func(playables ...interface{}) interface{} {
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Print(notify.Warningf("cannot loop (%T) %v", p, p))
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			if len(joined) == 1 {
				return core.NewLoop(ctx, joined[0])
			}
			return core.NewLoop(ctx, op.Join{Target: joined})
		}}

	eval["begin"] = Function{
		Title:         "Begin loop command",
		Description:   "begin loop(s). Ignore if it was running.",
		ControlsAudio: true,
		Prefix:        "beg",
		Template:      `begin(${1:loop})`,
		Samples: `lp_cb = loop(sequence('C D E F G A B'))
begin(lp_cb) // end(lp_cb)`,
		Func: func(vars ...variable) interface{} {
			for _, each := range vars {
				l, ok := each.Value().(*core.Loop)
				if !ok {
					notify.Print(notify.Warningf("cannot begin (%T) %v", l, l))
					continue
				}
				ctx.Control().StartLoop(l)
				notify.Print(notify.Infof("started loop: %s", each.Name))
			}
			return nil
		}}

	eval["end"] = Function{
		Title:         "End loop command",
		Description:   "end running loop(s). Ignore if it was stopped.",
		ControlsAudio: true,
		Template:      `end(${1:loop-or-empty})`,
		Samples: `l1 = loop(sequence('C E G'))
begin(l1) // end(l1)`,
		Func: func(vars ...variable) interface{} {
			if len(vars) == 0 {
				StopAllLoops(ctx)
				return nil
			}
			for _, each := range vars {
				l, ok := each.Value().(*core.Loop)
				if !ok {
					notify.Print(notify.Warningf("cannot end (%T) %v", l, l))
					continue
				}
				notify.Print(notify.Infof("stopping loop: %s", each.Name))
				ctx.Control().EndLoop(l)
			}
			return nil
		}}
	// END Loop and control
	eval["channel"] = Function{
		Title:         "MIDI channel operator",
		Description:   "select a MIDI channel, must be in [1..16]; must be a top-level operator",
		ControlsAudio: true,
		Prefix:        "chan",
		Alias:         "Ch",
		Template:      `channel(${1:number},${2:sequenceable})`,
		Samples:       `channel(2,sequence('C2 E3')) // plays on instrument connected to MIDI channel 2`,
		Func: func(midiChannel, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot decorate with channel (%T) %v", m, m))
			}
			return core.ChannelSelector{Target: s, Number: getValueable(midiChannel)}
		}}
	eval["interval"] = Function{
		Title:       "Interval creator",
		Description: "create an integer repeating interval (from,to,by,method). Default method is 'repeat', Use next() to get a new integer",
		Prefix:      "int",
		Alias:       "I",
		Template:    `interval(${1:from},${2:to},${3:by})`,
		Samples: `int1 = interval(-2,4,1)
lp_cdef = loop(pitch(int1,sequence('C D E F')), next(int1))`,
		IsComposer: true,
		Func: func(from, to, by interface{}) *core.Interval {
			return core.NewInterval(core.ToValueable(from), core.ToValueable(to), core.ToValueable(by), core.RepeatFromTo)
		}}

	eval["sequencemap"] = Function{
		Title:       "Sequence Map creator",
		Description: "creates a mapper of sequence notes by index (1-based)",
		Prefix:      "ind",
		Alias:       "Im",
		Template:    `sequencemap('${1:space-separated-1-based-indices}',${2:sequenceable})`,
		Samples: `s1 = sequence('C D E F G A B')
i1 = sequencemap('6 5 4 3 2 1',s1) // => B A G F E D
i2 = sequencemap('(6 5) 4 3 (2 1)',s1) // => (B A) G F (E D)`,
		IsComposer: true,
		Func: func(pattern, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot create sequence mapper on (%T) %v", m, m))
			}
			return op.NewSequenceMapper(s, core.ToValueable(pattern))
		}}

	eval["notemap"] = Function{
		Title:       "Note Map creator",
		Description: "creates a mapper of notes by index (1-based) or using dots (.) and bangs (!)",
		Template:    `notemap('${1:space-separated-1-based-indices-or-dots-and-bangs}',${2:note})`,
		IsComposer:  true,
		Samples: `m1 = notemap('..!..!..!', note('c2'))
m2 = notemap('3 6 9', note('d2'))`,
		Func: func(indices string, note interface{}) interface{} {
			m, err := op.NewNoteMap(indices, getValueable(note))
			if err != nil {
				return notify.Panic(fmt.Errorf("cannot create notemap, error:%v", err))
			}
			return m
		}}

	eval["merge"] = Function{
		Title:       "Merge creator",
		Description: `merges multiple sequences into one sequence`,
		Template:    `merge(${1:sequenceable})`,
		Samples: `m1 = notemap('..!..!..!', note('c2'))
m2 = notemap('4 7 10', note('d2'))
all = merge(m1,m2) // => = = C2 D2 = C2 D2 = C2 D2 = =`,
		IsComposer: true,
		Func: func(seqs ...interface{}) op.Merge {
			s := []core.Sequenceable{}
			for _, each := range seqs {
				seq, ok := getSequenceable(each)
				if ok {
					s = append(s, seq)
				} else {
					notify.Panic(fmt.Errorf("cannot merge (%T) %v", each, each))
				}
			}
			return op.Merge{Target: s}
		}}

	eval["next"] = Function{
		Title:       "Next operator",
		Description: `is used to produce the next value in a generator such as random and interval`,
		Samples: `i = interval(-4,4,2)
pi = pitch(i,sequence('C D E F G A B'))
lp_pi = loop(pi,next(i))
begin(lp_pi)`,
		Func: func(v interface{}) interface{} {
			return core.Nexter{Target: getValueable(v)}
		}}

	eval["export"] = Function{
		Title:       "Export command",
		Description: `writes a multi-track MIDI file`,
		Template:    `export(${1:filename},${2:sequenceable})`,
		Samples:     `export('myMelody-v1',myObject)`,
		Func: func(filename string, m interface{}) interface{} {
			if len(filename) == 0 {
				return notify.Panic(fmt.Errorf("missing filename to export MIDI %v", m))
			}
			_, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot MIDI export (%T) %v", m, m))
			}
			if !strings.HasSuffix(filename, "mid") {
				filename += ".mid"
			}
			return file.Export(filename, getValue(m), ctx.Control().BPM())
		}}

	eval["replace"] = Function{
		Title:       "Replace operator",
		Description: `replaces all occurrences of one musical object with another object for a given composed musical object`,
		Template:    `replace(${1:target},${2:from},${3:to})`,
		Samples: `c = note('c')
d = note('d')
pitchA = pitch(1,c)
pitchD = replace(pitchA, c, d) // c -> d in pitchA`,
		Func: func(target interface{}, from, to interface{}) interface{} {
			targetS, ok := getSequenceable(target)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot create replace inside (%T) %v", target, target))
			}
			fromS, ok := getSequenceable(from)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot create replace (%T) %v", from, from))
			}
			toS, ok := getSequenceable(to)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot create replace with (%T) %v", to, to))
			}
			return op.Replace{Target: targetS, From: fromS, To: toS}
		}}

	return eval
}

func registerFunction(m map[string]Function, k string, f Function) {
	if dup, ok := m[k]; ok {
		log.Fatal("duplicate function key detected:", dup)
	}
	if len(f.Alias) > 0 {
		if dup, ok := m[f.Alias]; ok {
			log.Fatal("duplicate function alias key detected:", dup)
		}
	}
	m[k] = f
	if len(f.Alias) > 0 {
		// modify title
		f.Title = fmt.Sprintf("%s [%s]", f.Title, f.Alias)
		m[f.Alias] = f
	}
}

func getSequenceable(v interface{}) (core.Sequenceable, bool) {
	if s, ok := v.(core.Sequenceable); ok {
		return s, ok
	}
	return nil, false
}

func getValueable(val interface{}) core.Valueable {
	if v, ok := val.(core.Valueable); ok {
		return v
	}
	return core.On(val)
}

// getValue returns the Value() of val iff val is a Valueable, else returns val
func getValue(val interface{}) interface{} {
	if v, ok := val.(core.Valueable); ok {
		return v.Value()
	}
	return val
}
