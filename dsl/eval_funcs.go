package dsl

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/emicklei/melrose/control"
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"

	"github.com/emicklei/melrose/midi/file"

	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/melrose/op"
)

// SyntaxVersion tells what language version this package is supporting.
const SyntaxVersion = "0.36" // major,minor

func IsCompatibleSyntax(s string) bool {
	if len(s) == 0 {
		// ignore syntax ; you are on your own
		return true
	}
	mm := strings.Split(SyntaxVersion, ".")
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

func (f Function) HumanizedTemplate() string {
	r := strings.NewReplacer(
		"${1:", "",
		"${2:", "",
		"${3:", "",
		"${4:", "",
		"}", "")
	return r.Replace(f.Template)
}

func EvalFunctions(ctx core.Context) map[string]Function {
	eval := map[string]Function{}

	// TODO allow fractions:  0.5, 0.25, 0.0125
	eval["fraction"] = Function{
		Title: "Duration fraction operator",
		Description: `Creates a new object for which the fraction of duration of all notes are changed.
The first parameter controls the fraction of the note, e.g. 1 = whole, 2 = half, 4 = quarter, 8 = eight, 16 = sixteenth.
Fraction can also be an exact float value between 0 and 1.
`,
		Prefix:     "fra",
		IsComposer: true,
		Template:   `fraction(${1:object},${2:object})`,
		Samples:    `fraction(8,sequence('e f')) // => ⅛E ⅛F , shorten the notes from quarter to eight`,
		Func: func(param float64, playables ...interface{}) interface{} {
			// if err := op.CheckFraction(param); err != nil {
			// 	notify.Print(notify.Error(err))
			// 	return nil
			// }
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Warnf("cannot fraction (%T) %v", p, p)
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
		Samples:    `dynamic('++',sequence('e f')) // => E++ F++`,
		Func: func(emphasis string, playables ...interface{}) interface{} {
			if err := op.CheckDynamic(emphasis); err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Warnf("cannot dynamic (%T) %v", p, p)
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return op.Dynamic{Target: joined, Emphasis: emphasis}
		}}

	eval["dynamicmap"] = Function{
		Title:       "Dynamic Map creator",
		Description: `changes the dynamic of notes from a musical object. 1-index-based mapping`,
		Prefix:      "dyna",
		IsComposer:  true,
		Template:    `dynamicmap('${1:mapping}',${2:object})`,
		Samples: `dynamicmap('1:++,2:--',sequence('e f')) // => E++ F--
dynamicmap('2:o,1:++,2:--,1:++', sequence('a b') // => B A++ B-- A++`,
		Func: func(mapping string, playables ...interface{}) interface{} {
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					return notify.Panic(fmt.Errorf("cannot dynamicmap (%T) %v", p, p))
				} else {
					joined = append(joined, s)
				}
			}
			mapper, err := op.NewDynamicMap(joined, mapping)
			if err != nil {
				notify.Warningf("cannot create dynamic mapping %v", err)
				return nil
			}
			return mapper
		}}

	eval["progression"] = Function{
		Title: "Chord progression creator",
		//Description: `create a Chord progression using this <a href="/melrose/notations.html#progression-not">format</a>`,
		Prefix:   "pro",
		IsCore:   true,
		Template: `progression('${1:scale}','${2:space-separated-roman-chords}')`,
		Samples:  `progression('C','II V I') // => (D F A) (G B D5) (C E G)`,
		Func: func(scale, chords interface{}) interface{} {
			return core.NewChordProgression(getValueable(scale), getValueable(chords))
		}}

	eval["chordsequence"] = Function{
		Title:       "Sequence of chords creator",
		Description: `create a Chord sequence using this <a href="/melrose/notations.html#chordsequence-not">format</a>`,
		Prefix:      "pro",
		IsCore:      true,
		Template:    `chordsequence('${1:chords}')`,
		Samples: `chordsequence('e f') // => (E A♭ B) (F A C5)
		chordsequence('(c d)') // => (C E G D G♭ A)`,
		Func: func(chords string) interface{} {
			p, err := core.ParseChordSequence(chords)
			if err != nil {
				return notify.Panic(err)
			}
			return p
		}}

	eval["joinmap"] = Function{
		Title:       "Join Map creator",
		Description: "creates a new join by mapping elements. 1-index-based mapping",
		Prefix:      "joinm",
		IsComposer:  true,
		Template:    `joinmap('${1:indices}',${2:join})`,
		Samples: `j = join(note('c'), sequence('d e f'))
jm = joinmap('1 (2 3) 4',j) // => C = D =`,
		Func: func(indices interface{}, join interface{}) interface{} { // allow multiple seq?
			v := getValueable(join)
			vNow := v.Value()
			if _, ok := vNow.(op.Join); !ok {
				return notify.Panic(fmt.Errorf("cannot joinmap (%T) %v, must be a join", join, join))
			}
			p := getValueable(indices)
			return op.NewJoinMap(v, p)
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
		Template:    `track('${1:title}',${2:midi-channel}, onbar(1,${3:object}))`,
		Samples:     `track("lullaby",1,onbar(2, sequence('c d e'))) // => a new track on MIDI channel 1 with sequence starting at bar`,
		Func: func(title string, channel int, onbars ...core.SequenceOnTrack) interface{} {
			if len(title) == 0 {
				return notify.Panic(fmt.Errorf("cannot have a track without title"))
			}
			if channel < 1 || channel > 15 {
				return notify.Panic(fmt.Errorf("MIDI channel must be in [1..15]"))
			}
			tr := core.NewTrack(title, channel)
			for _, each := range onbars {
				tr.Add(each)
			}
			return tr
		}}

	eval["multitrack"] = Function{
		Title:         "Multi track creator",
		Description:   "create a multi-track object from zero or more tracks",
		Prefix:        "mtr",
		Template:      `multitrack(${1:track})`,
		Samples:       `multitrack(track1,track2,track3) // 3 tracks in one multi-track object`,
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
The first parameter is a fraction {1,2,4,8,16} or a duration in milliseconds or a time.Duration.
Second parameter is the MIDI number and must be one of [0..127].
The third parameter is the velocity (~ loudness) and must be one of [0..127]`,
		Prefix:   "mid",
		Template: `midi(${1:numberOrDuration},${2:number},${3:number})`,
		Samples: `midi(500,52,80) // => 500ms E3+
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
		Description: "prints an object when evaluated (play,loop)",
		Func: func(m interface{}) interface{} {
			return core.Watch{Context: ctx, Target: m}
		}}

	eval["chord"] = Function{
		Title:       "Chord creator",
		Description: `create a Chord from its string <a href="/melrose/notations.html#chord-not">notation</a>`,
		Prefix:      "cho",
		Template:    `chord('${1:note}')`,
		Samples: `chord('c#5/m/1')
chord('g/M/2') // Major G second inversion`,
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
		Description: "create a sequence with notes for which the order and the octaves are changed. 1-based indexing",
		Prefix:      "octavem",
		Template:    `octavemap('${1:int2int}',${2:object})`,
		IsComposer:  true,
		Samples:     `octavemap('1:-1,2:0,3:1',chord('c')) // => (C3 E G5)`,
		Func: func(indices string, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot octavemap (%T) %v", m, m))
			}
			return op.NewOctaveMap(s, indices)
		}}

	eval["pitchmap"] = Function{
		Title:       "Pitch Map operator",
		Description: "create a sequence with notes for which the order and the pitch are changed. 1-based indexing",
		Prefix:      "pitchm",
		Template:    `pitchmap('${1:int2int}',${2:object})`,
		IsComposer:  true,
		Samples:     `pitchmap('1:-1,1:0,1:1',note('c')) // => B3 C D`,
		Func: func(indices string, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot pitchmap (%T) %v", m, m))
			}
			return op.NewPitchMap(s, indices)
		}}

	eval["pitch"] = Function{
		Title:       "Pitch operator",
		Description: "change the pitch with a delta of semitones",
		Prefix:      "pit",
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
		Template:    `reverse(${1:sequenceable})`,
		Samples:     `reverse(chord('a'))`,
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
		Description: "repeat one or more musical objects a number of times",
		Prefix:      "rep",
		Template:    `repeat(${1:times},${2:sequenceables})`,
		Samples:     `repeat(4,sequence('c d e'))`,
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
		Description: "joins one or more musical objects as one",
		Prefix:      "joi",
		Template:    `join(${1:first},${2:second})`,
		Samples: `a = chord('a')
b = sequence('(c e g)')
ab = join(a,b) // => (A D♭5 E5) (C E G)`,
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
		Samples: `bpm(90)
speedup = iterator(80,100,120,140)
l = loop(bpm(speedup),sequence('c e g'),next(speedup))`,
		Func: func(v interface{}) interface{} {
			return control.NewBPM(core.On(v), ctx)
		}}

	eval["duration"] = Function{
		Title:       "Duration calculator",
		Description: "computes the duration of the object using the current BPM",
		Prefix:      "dur",
		Template:    `duration(${1:object})`,
		Samples:     `duration(note('c')) // => 375ms`,
		Func: func(m interface{}) time.Duration {
			if s, ok := getSequenceable(m); ok {
				return s.S().Duration(ctx.Control().BPM())
			}
			return time.Duration(0)
		}}

	eval["biab"] = Function{
		Title:         "Beats in a Bar",
		Description:   "set the Beats in a Bar; default is 4",
		ControlsAudio: true,
		Prefix:        "biab",
		Template:      `biab(${1:beats-in-a-bar})`,
		Samples:       `biab(4)`,
		Func: func(i int) interface{} {
			if i < 1 {
				return notify.Panic(fmt.Errorf("invalid beats-in-a-bar, must be positive, %d = ", i))
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
		Template:    `sequence('${1:space-separated-notes}')`,
		Samples: `sequence('c d e')
sequence('(8c d e)') // => (⅛C D E)
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
		Template:    `scale(${1:octaves},'${2:scale-syntax}')`,
		IsCore:      true,
		Samples:     `scale(1,'e/m') // => E F G A B C5 D5`,
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
		Samples:     `at(1,scale('e/m')) // => E`,
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
			return core.NewSequenceOnTrack(getValueable(bar), s)
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
		Title:         "Play musical objects in order. Use sync() for parallel playing",
		Description:   "play all musical objects",
		ControlsAudio: true,
		Prefix:        "pla",
		Template:      `play(${1:sequenceable})`,
		Samples:       `play(s1,s2,s3) // play s3 after s2 after s1`,
		Func: func(playables ...interface{}) interface{} {
			list := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok { // unwrap var
					list = append(list, s)
				} else {
					notify.Warnf("cannot play (%T) %v", p, p)
				}
			}
			return control.NewPlay(ctx, list, false)
		}}

	eval["sync"] = Function{
		Title:         "Synchroniser creator",
		Description:   "Synchronise playing musical objects. Use play() for serial playing",
		ControlsAudio: true,
		Prefix:        "syn",
		Template:      `sync(${1:object})`,
		Samples: `sync(s1,s2,s3) // play s1,s2 and s3 at the same time
sync(loop1,loop2) // begin loop2 at the next start of loop1`,
		Func: func(playables ...interface{}) interface{} {
			vals := []core.Valueable{}
			for _, p := range playables {
				vals = append(vals, getValueable(p))
			}
			return control.NewSyncPlay(vals)
		}}

	eval["ungroup"] = Function{
		Title:       "Ungroup operator",
		Description: "undo any grouping of notes from one or more musical objects",
		Prefix:      "ung",
		Template:    `ungroup(${1:sequenceable})`,
		IsComposer:  true,
		Samples: `ungroup(chord('e')) // => E G B
ungroup(sequence('(c d)'),note('e')) // => C D E`,
		Func: func(playables ...interface{}) interface{} {
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Warningf("cannot ungroup (%T) %v", p, p)
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
		Samples:     `octave(1,sequence('c d')) // => C5 D5`,
		Func: func(scalarOrVar interface{}, playables ...interface{}) interface{} {
			list := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Warningf("cannot octave (%T) %v", p, p)
					return nil
				} else {
					list = append(list, s)
				}
			}
			return op.Octave{Target: list, Offset: core.ToValueable(scalarOrVar)}
		}}

	// 	eval["record"] = Function{
	// 		Title:         "Recording creator",
	// 		Description:   "create a recorded sequence of notes from the current MIDI input device",
	// 		ControlsAudio: true,
	// 		Prefix:        "rec",
	// 		Template:      `record()`,
	// 		Samples: `r = record() // record notes played on the current input device and stop recording after 5 seconds
	// s = r.S() // returns the sequence of notes from the recording`,
	// 		Func: func() interface{} {
	// 			seq, err := ctx.Device().Record(ctx)
	// 			if err != nil {
	// 				return notify.Panic(err)
	// 			}
	// 			return seq
	// 		}}

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

	registerFunction(eval, "stretch", Function{
		Title:       "Stretch operator",
		Description: "stretches the duration of musical object(s) with a factor. If the factor < 1 then duration is shortened",
		Prefix:      "st",
		Template:    `stretch(${1:factor},${2:object})`,
		Samples: `stretch(2,note('c'))  // 2C
stretch(0.25,sequence('(c e g)'))  // (16C 16E 16G)
stretch(8,note('c'))  // C with length of 2 bars`,
		Func: func(factor interface{}, m ...interface{}) interface{} {
			list, ok := getSequenceableList(m...)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot stretch (%T) %v", m, m))
			}
			return op.NewStretch(getValueable(factor), list)
		}})

	eval["group"] = Function{
		Title:       "Group operator",
		Description: "create a new sequence in which all notes of a musical object are grouped",
		Prefix:      "par",
		Template:    `group(${1:sequenceable})`,
		Samples:     `group(sequence('c d e')) // => (C D E)`,
		IsComposer:  true,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				return notify.Panic(fmt.Errorf("cannot group (%T) %v", value, value))
			} else {
				return op.Group{Target: s}
			}
		}}
	// BEGIN Loop and control
	eval["loop"] = Function{
		Title:         "Loop creator",
		Description:   "create a new loop from one or more musical objects; must be assigned to a variable",
		ControlsAudio: true,
		Prefix:        "loo",
		Template:      `lp_${1:object} = loop(${1:object})`,
		Samples: `cb = sequence('c d e f g a b')
lp_cb = loop(cb,reverse(cb))`,
		Func: func(playables ...interface{}) interface{} {
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Warnf("cannot loop (%T) %v", p, p)
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return core.NewLoop(ctx, joined)
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
					notify.Warnf("cannot begin (%T) %v", l, l)
					continue
				}
				_ = l.Play(ctx, time.Now())
				notify.Print(notify.Infof("begin %s", each.Name))
			}
			return nil
		}}

	eval["end"] = Function{
		Title:         "End loop or listen command",
		Description:   "end running loop(s) or listener(s). Ignore if it was stopped.",
		ControlsAudio: true,
		Template:      `end(${1:control})`,
		Samples: `l1 = loop(sequence('c e g'))
begin(l1)
end(l1)
end() // stop all playables`,
		Func: func(vars ...variable) interface{} {
			if len(vars) == 0 {
				StopAllPlayables(ctx)
				return nil
			}
			for _, each := range vars {
				if l, ok := each.Value().(core.Playable); ok {
					notify.Print(notify.Infof("ending %s", each.Name))
					_ = l.Stop(ctx)
				} else {
					notify.Warnf("cannot end (%T) %v", each.Value(), each.Value())
				}
			}
			return nil
		}}
	// END Loop and control
	eval["channel"] = Function{
		Title:         "MIDI channel selector",
		Description:   "select a MIDI channel, must be in [1..16]; must be a top-level operator",
		ControlsAudio: true,
		Prefix:        "chan",
		Template:      `channel(${1:number},${2:sequenceable})`,
		Samples:       `channel(2,note('g3'), sequence('c2 e3')) // plays on instrument connected to MIDI channel 2`,
		Func: func(midiChannel interface{}, m ...interface{}) interface{} {
			list, ok := getSequenceableList(m...)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot decorate with channel (%T) %v", m, m))
			}
			return core.NewChannelSelector(list, getValueable(midiChannel))
		}}

	eval["device"] = Function{
		Title:         "MIDI device selector",
		Description:   "select a MIDI device from the available device IDs; must become before channel",
		ControlsAudio: true,
		Prefix:        "dev",
		Template:      `device(${1:number},${2:sequenceable})`,
		Samples:       `device(1,channel(2,sequence('c2 e3'), note('g3'))) // plays on connected device 1 through MIDI channel 2`,
		Func: func(deviceID interface{}, m ...interface{}) interface{} {
			list, ok := getSequenceableList(m...)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot decorate with device (%T) %v", m, m))
			}
			return core.NewDeviceSelector(list, getValueable(deviceID))
		}}

	eval["interval"] = Function{
		Title:       "Interval creator",
		Description: "create an integer repeating interval (from,to,by,method). Default method is 'repeat', Use next() to get a new integer",
		Prefix:      "int",
		Template:    `interval(${1:from},${2:to},${3:by})`,
		Samples: `int1 = interval(-2,4,1)
lp_cdef = loop(pitch(int1,sequence('c d e f')), next(int1))`,
		IsComposer: true,
		Func: func(from, to, by interface{}) *core.Interval {
			return core.NewInterval(core.ToValueable(from), core.ToValueable(to), core.ToValueable(by), core.RepeatFromTo)
		}}

	eval["resequence"] = Function{
		Title:       "Sequence modifier",
		Description: "creates a modifier of sequence notes by index (1-based)",
		Prefix:      "resq",
		Template:    `resequence('${1:space-separated-1-based-indices}',${2:sequenceable})`,
		Samples: `s1 = sequence('C D E F G A B')
i1 = resequence('6 5 4 3 2 1',s1) // => B A G F E D
i2 = resequence('(6 5) 4 3 (2 1)',s1) // => (B A) G F (E D)`,
		IsComposer: true,
		Func: func(pattern, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot create resequencer on (%T) %v", m, m))
			}
			return op.NewResequencer(s, core.ToValueable(pattern))
		}}

	eval["notemap"] = Function{
		Title:       "Note Map creator",
		Description: "creates a mapper of notes by index (1-based) or using dots (.) and bangs (!)",
		Template:    `notemap('${1:space-separated-1-based-indices-or-dots-and-bangs}',${2:has-note})`,
		IsComposer:  true,
		Samples: `m1 = notemap('..!..!..!', note('c2'))
m2 = notemap('3 6 9', octave(-1,note('d2')))`,
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
		Description: `is used to produce the next value in a generator such as random, iterator and interval`,
		Samples: `i = interval(-4,4,2)
pi = pitch(i,sequence('c d e f g a b')) // current value of "i" is used
lp_pi = loop(pi,next(i)) // "i" will advance to the next value
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

	eval["midi_send"] = Function{
		Title:       "Send MIDI message",
		Description: "Sends a MIDI message with status, channel(ignore if < 1), 2nd byte and 3rd byte to an output device. Can be used as a musical object",
		Template:    "midi_send(${1:device-id},${1:status},${2:channel},${3:2nd-byte},${4:3rd-byte}",
		Samples: `midi_send(1,0xB0,7,0x7B,0) // to device id 1, control change, all notes off in channel 7
midi_send(1,0xC0,2,1,0) // program change, select program 1 for channel 2
midi_send(2,0xB0,4,0,16) // control change, bank select 16 for channel 4
midi_send(3,0xB0,1,120,0) // control change, all notes off for channel 1`,
		Func: func(deviceID int, status int, channel, data1, data2 interface{}) interface{} {
			return midi.NewMessage(ctx.Device(), core.On(deviceID), status, core.On(channel), core.On(data1), core.On(data2))
		}}

	registerFunction(eval, "listen", Function{
		Title:       "Start a MIDI listener",
		Description: "Listen for note(s) from a device and call a playable function to handle",
		Template:    "listen(${1:variable-or-device-selector},${2:function})",
		Samples: `rec = note('c') // define a variable "rec" with a initial object ; this is a place holder
fun = play(rec) // define the playable function to call when notes are received ; loop and print are also possible
ear = listen(rec,fun) // start a listener for notes from default input device, store it in "rec" and call "fun"
alt = listen(device(1,rec),fun) // start a listener for notes from input device 1`,
		Func: func(varOrDeviceSelector interface{}, function interface{}) interface{} {
			_, ok := getValue(function).(core.Evaluatable)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot listen and call (%T) %s", function, core.Storex(function)))
			}
			var injectable variable
			deviceID, _ := ctx.Device().DefaultDeviceIDs()
			if ds, ok := varOrDeviceSelector.(core.DeviceSelector); ok {
				deviceID = ds.DeviceID()
				if len(ds.Target) == 0 {
					return notify.Panic(fmt.Errorf("missing variable parameter"))
				}
				first := ds.Target[0]
				if v, ok := first.(variable); ok {
					injectable = v
				} else {
					return notify.Panic(fmt.Errorf("missing variable parameter"))
				}
			} else {
				// must be variable
				if v, ok := varOrDeviceSelector.(variable); ok {
					injectable = v
				} else {
					return notify.Panic(fmt.Errorf("missing variable parameter"))
				}
			}
			// use function as Valueable and not the Evaluatable to allow redefinition of the callback function in the script
			return control.NewListen(ctx, deviceID, injectable.Name, getValueable(function))
		},
	})

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

func getSequenceableList(m ...interface{}) (list []core.Sequenceable, ok bool) {
	ok = true
	for _, each := range m {
		if s, ok := getSequenceable(each); ok {
			list = append(list, s)
		} else {
			return list, false
		}
	}
	return
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

func getLoop(v interface{}) (core.Valueable, bool) {
	if val, ok := v.(core.Valueable); ok {
		if _, ok := val.Value().(*core.Loop); ok {
			return val, true
		}
		return val, false
	}
	if l, ok := v.(*core.Loop); ok {
		return core.On(l), true
	}
	return core.On(v), false
}
