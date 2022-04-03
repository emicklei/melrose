package dsl

import (
	"errors"
	"fmt"
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
const SyntaxVersion = "0.37" // major,minor

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
		Samples:    `fraction(8,sequence('e f')) // => 8E 8F , shorten the notes from quarter to eight`,
		Func: func(param interface{}, playables ...interface{}) interface{} {
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Warnf("cannot fraction (%T) %v", p, p)
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return op.NewFraction(getHasValue(param), joined)
		}}

	eval["dynamic"] = Function{
		Title: "Dynamic operator",
		Description: `Creates a new modified musical object for which the dynamics of all notes are changed.
	The first parameter controls the emphasis the note, e.g. + (mezzoforte,mf), -- (piano,p) or a velocity [0..127].
	`,
		Prefix:     "dy",
		IsComposer: true,
		Template:   `dynamic(${1:emphasis},${2:object})`,
		Samples: `dynamic('++',sequence('e f')) // => E++ F++
dynamic(112,note('a')) // => A++++`,
		Func: func(emphasis interface{}, playables ...interface{}) interface{} {
			joined := []core.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Warnf("cannot dynamic (%T) %v", p, p)
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return op.Dynamic{Target: joined, Emphasis: getHasValue(emphasis)}
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
				notify.NewWarningf("cannot create dynamic mapping %v", err)
				return nil
			}
			return mapper
		}}

	eval["progression"] = Function{
		Title:       "Chord progression creator",
		Description: `create a Chord progression using this <a href="/docs/reference/notations/#chordprogression">format</a>`,
		Prefix:      "pro",
		IsCore:      true,
		Template:    `progression('${1:scale}','${2:space-separated-roman-chords}')`,
		Samples:     `progression('C','II V I') // => (D F A) (G B D5) (C E G)`,
		Func: func(scale, chords interface{}) interface{} {
			return core.NewChordProgression(getHasValue(scale), getHasValue(chords))
		}}

	eval["chordsequence"] = Function{
		Title:       "Sequence of chords creator",
		Description: `create a Chord sequence using this <a href="/docs/reference/notations/#chordsequence">format</a>`,
		Prefix:      "pro",
		IsCore:      true,
		Template:    `chordsequence('${1:chords}')`,
		Samples: `chordsequence('e f') // => (E A_ B) (F A C5)
chordsequence('(c d)') // => (C E G D G_ A)`,
		Func: func(chords string) interface{} {
			p, err := core.ParseChordSequence(chords)
			if err != nil {
				return notify.Panic(err)
			}
			return p
		}}

	eval["prob"] = Function{
		Title:    "Probabilistic music object.",
		Prefix:   "prob",
		IsCore:   true,
		Template: `prob(${1:perc},${2:note-or-sequenceable})`,
		Samples: `prob(50,note('c')) // 50% chance of playing the note C, otherwise a quarter rest
prob(0.8,sequence('(c e g)')) // 80% chance of playing the chord C, otherwise a quarter rest`,
		Func: func(prec interface{}, noteOrSeq interface{}) interface{} {
			return op.NewProbability(getHasValue(prec), getHasValue(noteOrSeq))
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
			v := getHasValue(join)
			vNow := v.Value()
			if _, ok := vNow.(op.Join); !ok {
				return notify.Panic(fmt.Errorf("cannot joinmap (%T) %v, must be a join", join, join))
			}
			p := getHasValue(indices)
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
		Samples:     `track("lullaby",1,onbar(2, sequence('c d e'))) // => a new track on MIDI channel 1 with sequence starting at bar 2`,
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
			tracks := []core.HasValue{}
			for _, each := range varOrTrack {
				tracks = append(tracks, getHasValue(each))
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
			durVal := getHasValue(dur)
			nrVal := getHasValue(nr)
			velVal := getHasValue(velocity)
			return core.NewMIDI(durVal, nrVal, velVal)
		}}

	eval["print"] = Function{
		Title:       "Printer creator",
		Description: "prints an object when evaluated (play,loop)",
		Func: func(m interface{}) interface{} {
			return core.Print{Context: ctx, Target: m}
		}}

	eval["chord"] = Function{
		Title:       "Chord creator",
		Description: `create a Chord from its string <a href="/docs/reference/notations/#chord">format</a>`,
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

	registerFunction(eval, "transposemap", Function{
		Title:       "Transpose Map operator",
		Description: "create a sequence with notes for which the order and the pitch are changed. 1-based indexing",
		Alias:       "pitchmap",
		Template:    `transposemap('${1:int2int}',${2:object})`,
		IsComposer:  true,
		Samples:     `transposemap('1:-1,1:0,1:1',note('c')) // => B3 C D`,
		Func: func(indices string, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot transposemap (%T) %v", m, m))
			}
			return op.NewTransposeMap(s, indices)
		}})

	eval["octavemap"] = Function{
		Title:       "Octave Map operator",
		Description: "create a sequence with notes for which the order and the octaves are changed",
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

	eval["velocitymap"] = Function{
		Title:       "Velocity Map operator",
		Description: "create a sequence with notes for which the order and the velocities are changed. Velocity 0 means no change.",
		Prefix:      "velocitym",
		Template:    `velocitymap('${1:int2int}',${2:object})`,
		IsComposer:  true,
		Samples:     `velocitymap('1:30,2:0,3:60',chord('c')) // => (C3--- E G5+)`,
		Func: func(indices string, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot velocitymap (%T) %v", m, m))
			}
			return op.NewVelocityMap(s, indices)
		}}

	registerFunction(eval, "transpose", Function{
		Title:       "Transpose operator",
		Description: "change the pitch with a delta of semitones",
		Alias:       "pitch",
		Prefix:      "tran",
		Template:    `transpose(${1:semitones},${2:sequenceable})`,
		Samples: `transpose(-1,sequence('c d e'))
p = interval(-4,4,1)
transpose(p,note('c'))`,
		IsComposer: true,
		Func: func(semitones, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot transpose (%T) %v", m, m))
			}
			return op.Transpose{Target: s, Semitones: getHasValue(semitones)}
		}})

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
			return op.Repeat{Target: joined, Times: getHasValue(howMany)}
		}}

	registerFunction(eval, "join", Function{
		Title:       "Join operator",
		Description: "joins one or more musical objects as one",
		Prefix:      "joi",
		Template:    `join(${1:first},${2:second})`,
		Samples: `a = chord('a')
b = sequence('(c e g)')
ab = join(a,b) // => (A D_5 E5) (C E G)`,
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
			if !ctx.Capabilities().ImportMelrose {
				return notify.NewWarningf("import not available")
			}
			err := ImportProgram(ctx, f)
			if err != nil {
				return notify.Panic(fmt.Errorf("failed to import [%s], %v", f, err))
			}
			return nil
		},
	})

	eval["sequence"] = Function{
		Title:       "Sequence creator",
		Description: `create a Sequence using this <a href="/docs/reference/notations/#sequence">format</a>`,
		Prefix:      "se",
		Alias:       "seq",
		Template:    `sequence('${1:space-separated-notes}')`,
		Samples: `sequence('c d e')
sequence('(8c d e)') // => (8C D E)
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
		Description: `create a Note using this <a href="/docs/reference/notations/#note">format</a>`,
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
		Description: `create a Scale using this <a href="/docs/reference/notations/#scale">format</a>`,
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
				notify.Print(notify.NewError(err))
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
			indexVal := getHasValue(index)
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
			return core.NewSequenceOnTrack(getHasValue(bar), s)
		}}

	eval["random"] = Function{
		Title:       "Random generator",
		Description: "create a random integer generator. Use next() to generate a new integer",
		Prefix:      "ra",
		Template:    `random(${1:from},${2:to})`,
		Samples: `num = random(1,10)
next(num)`,
		Func: func(from interface{}, to interface{}) interface{} {
			fromVal := getHasValue(from)
			toVal := getHasValue(to)
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
				// first check Playable
				if pl, ok := getPlayable(p); ok {
					pl.Play(ctx, time.Now())
					continue
				}
				fmt.Printf("not a playable %T\n", p)
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
			vals := []core.HasValue{}
			for _, p := range playables {
				vals = append(vals, getHasValue(p))
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
					notify.NewWarningf("cannot ungroup (%T) %v", p, p)
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
					notify.NewWarningf("cannot octave (%T) %v", p, p)
					return nil
				} else {
					list = append(list, s)
				}
			}
			return op.Octave{Target: list, Offset: core.ToHasValue(scalarOrVar)}
		}}

	eval["record"] = Function{
		Title:         "Recording creator",
		Description:   "create a recorded sequence of notes from the current MIDI input device using the currrent BPM",
		ControlsAudio: true,
		Template:      `record(rec)`,
		Samples: `rec = sequence('') // variable to store the recorded sequence
record(rec) // record notes played on the current input device`,
		Func: func(varOrDeviceSelector interface{}) interface{} {
			var injectable variable
			deviceID, _ := ctx.Device().DefaultDeviceIDs()
			if ds, ok := varOrDeviceSelector.(core.DeviceSelector); ok {
				deviceID = ds.DeviceID()
				first := ds.Target
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
			return control.NewRecording(deviceID, injectable.Name, ctx.Control().BPM())
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
		Template:    `iterator(${1:array-element})`,
		Samples: `i = iterator(1,3,5,7,9)
p = transpose(i,note('c'))
lp = loop(p,next(i))
		`,
		Func: func(values ...interface{}) *core.Iterator {
			return &core.Iterator{
				Target: values,
			}
		}})

	registerFunction(eval, "rotate", Function{
		Title:       "Rotation modifier",
		Description: "rotates note(groups) in a sequence. count is negative for rotating left",
		Template:    `rotate(${1:count},${2:object})`,
		Samples: `rotate(-1,sequence('C E G')) // E G C
			`,
		Func: func(count interface{}, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot rotate (%T) %v", s, s))
			}
			return op.Rotate{
				Times:  getHasValue(count),
				Target: s,
			}
		}})

	registerFunction(eval, "stretch", Function{
		Title:       "Stretch operator",
		Description: "stretches the duration of musical object(s) with a factor. If the factor < 1 then duration is shortened",
		Prefix:      "st",
		Template:    `stretch(${1:factor},${2:object})`,
		Samples: `stretch(2,note('c'))  // 2C
stretch(0.25,sequence('(c e g)'))  // (16C 16E 16G)
stretch(8,note('c'))  // C with length of 8 x 0.25 (quarter) = 2 bars`,
		Func: func(factor interface{}, m ...interface{}) interface{} {
			list, ok := getSequenceableList(m...)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot stretch (%T) %v", m, m))
			}
			return op.NewStretch(getHasValue(factor), list)
		}})

	eval["group"] = Function{
		Title:       "Group operator",
		Description: "create a new sequence in which all notes of a musical object are grouped",
		Prefix:      "gro",
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
		Description:   "create a new loop from one or more musical objects",
		ControlsAudio: true,
		Prefix:        "loo",
		Template:      `loop(${1:object})`,
		Samples: `cb = sequence('c d e f g a b')
loop(cb,reverse(cb))`,
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

	eval["stop"] = Function{
		Title:         "Stop a loop or listen",
		Description:   "stop running loop(s) or listener(s). Ignore if it was stopped.",
		ControlsAudio: true,
		Template:      `stop(${1:control})`,
		Samples: `l1 = loop(sequence('c e g'))
play(l1)
stop(l1)
stop() // stop all playables`,
		Func: func(vars ...variable) interface{} {
			if len(vars) == 0 {
				StopAllPlayables(ctx)
				return nil
			}
			for _, each := range vars {
				if l, ok := each.Value().(core.Stoppable); ok {
					notify.Infof("stopping %s", each.Name)
					_ = l.Stop(ctx)
				} else {
					notify.Warnf("cannot stop (%T) %v", each.Value(), each.Value())
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
		Samples:       `channel(2,sequence('c2 e3')) // plays on instrument connected to MIDI channel 2`,
		Func: func(midiChannel interface{}, m interface{}) interface{} {
			seq, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot decorate with channel (%T) %v", m, m))
			}
			return core.NewChannelSelector(seq, getHasValue(midiChannel))
		}}

	eval["fractionmap"] = Function{
		Title:       "Fraction Map operator",
		Description: "create a sequence with notes for which the fractions are changed. 1-based indexing. use space or comma as separator",
		Prefix:      "frm",
		Template:    `fractionmap('${1:fraction-mapping}',${2:object})`,
		IsComposer:  true,
		Samples: `fractionmap('3:. 2:4,1:2',sequence('c e g')) // => .G E 2C
fractionmap('. 8 2',sequence('c e g')) // => .C 8E 2G`,
		Func: func(indices interface{}, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot fractionmap (%T) %v", m, m))
			}
			return op.NewFractionMap(getHasValue(indices), s)
		}}

	// eval["input"] = Function{
	// 	Title: "MIDI Input device",
	// 	//Description:   "Look up an input device by name",
	// 	ControlsAudio: true,
	// 	Func: func(deviceName string, optionalChannel ...int) interface{} {
	// 		in, _ := ctx.Device().DefaultDeviceIDs()
	// 		return control.NewChannelOnDevice(true, deviceName, -1, in)
	// 	}}

	// 	eval["onpress"] = Function{
	// 		Title: "Computer keyboard key press",
	// 		Description: `Use the key to trigger playing.
	// If this key is pressed the playable will start.
	// If pressed again, the play will stop.
	// Remove the assignment using the value nil for the playable`,
	// 		Template: `onpress(${1:key},${2:playable-or-evaluatable-or-nil})`,
	// 		Samples: `loopA = loop(scale(2,'c'))
	// onpress('a',loopA)`,
	// 		ControlsAudio: true,
	// 		Func: func(char string, playOrEval interface{}) interface{} {
	// 			if len(char) == 0 {
	// 				return notify.Panic(fmt.Errorf("key cannot be empty"))
	// 			}
	// 			// allow nil, playable and evaluatable
	// 			if playOrEval == nil {
	// 				// uninstall binding
	// 				// TODO
	// 				return nil
	// 			}
	// 			return nil
	// 		}}

	eval["key"] = Function{
		Title:       "MIDI Keyboard key",
		Description: "Use the key to trigger the play of musical object",
		Template:    `key('${2:note}')`,
		Samples: `c2 = key('c2') // C2 key on the default input device and default channel
c2 = key(device(1,note('c2'))) // C2 key on input device 1
c2 = key(device(1,channel(2,note('c2'))) // C2 key on input device 1 and channel 2
c2 = key(channel(3,note('c2')) // C2 key on the default input device and channel 3`,
		ControlsAudio: true,
		Func: func(noteEntry interface{}) interface{} {
			// check string
			if s, ok := noteEntry.(string); ok {
				note, err := core.ParseNote(s)
				if err != nil {
					return notify.Panic(fmt.Errorf("cannot create Note with input %q", note))
				}
				return control.NewKey(1, 1, note)
			}
			deviceID, _ := ctx.Device().DefaultDeviceIDs()
			channel := 1 // TODO
			note := core.Rest4
			// check device
			if d, ok := getValue(noteEntry).(core.DeviceSelector); ok {
				deviceID = d.DeviceID()
				noteEntry = d.Target
			}
			// check channel
			if c, ok := getValue(noteEntry).(core.ChannelSelector); ok {
				channel = c.Channel()
				noteEntry = c.Target
			}
			// check note
			if n, ok := getValue(noteEntry).(core.Note); ok {
				note = n // TODO
			}
			return control.NewKey(deviceID, channel, note)
		}}

	eval["knob"] = Function{
		Title:       "MIDI controller knob",
		Description: "Use the knob as an integer value for a parameter in any object",
		Template:    `knob(${1:device-id},${2:midi-number})`,
		Samples: `axiom = 1 // device ID for my connected M-Audio Axiom 25
B1 = 20 // MIDI number assigned to this knob on the controller
k = knob(axiom,B1)
transpose(k,scale(1,'E')) // when played, use the current value of knob "k"`,
		ControlsAudio: true,
		Func: func(deviceIDOrVar interface{}, numberOrVar interface{}) interface{} {
			deviceID, ok := getValue(deviceIDOrVar).(int)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot create knob with device (%T) %v", deviceIDOrVar, deviceIDOrVar))
			}
			number, ok := getValue(numberOrVar).(int)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot create knob with number (%T) %v", numberOrVar, numberOrVar))
			}
			k := control.NewKnob(deviceID, 0, number)
			ctx.Device().Listen(deviceID, k, true)
			return k
		}}

	eval["onkey"] = Function{
		Title: "Key trigger creator",
		Description: `Assign a playable to a key.
If this key is pressed the playable will start. 
If pressed again, the play will stop.
Remove the assignment using the value nil for the playable`,
		ControlsAudio: true,
		Prefix:        "onk",
		Template:      `onkey(${1:key},${2:playable-or-evaluatable-or-nil})`,
		Samples: `onkey('c',myLoop) // on the default input device, when C4 is pressed then start or stop myLoop

axiom = 1 // device ID for the M-Audio Axiom 25
c2 = key(device(axiom,note('c2')))
fun = play(scale(2,'c')) // what to do when a key is pressed (NoteOn)
onkey(c2, fun) // if C2 is pressed on the axiom device then evaluate the function "fun"`,
		Func: func(keyOrVar interface{}, playOrEval interface{}) interface{} {
			if !ctx.Device().HasInputCapability() {
				return notify.Panic(errors.New("input is not available for this device"))
			}
			var key control.Key
			// key is mandatory
			noteName, ok := getValue(keyOrVar).(string)
			if ok {
				note := core.MustParseNote(noteName)
				// it is a note name on the default input device
				in, _ := ctx.Device().DefaultDeviceIDs()
				key = control.NewKey(in, 1, note) // TODO what is channel on default input dev?
			} else {
				keyVar, ok := getValue(keyOrVar).(control.Key)
				if !ok {
					return notify.Panic(fmt.Errorf("cannot install onkey because parameter is not a key (%T) %v", keyOrVar, keyOrVar))
				}
				key = keyVar
			}
			// allow nil, playable and evaluatable
			if playOrEval == nil {
				// uninstall binding
				ctx.Device().OnKey(ctx, key.DeviceID(), key.Channel(), key.Note(), nil)
				return nil
			}
			_, ok = getValue(playOrEval).(core.Playable)
			if !ok {
				_, ok = getValue(playOrEval).(core.Evaluatable)
				if !ok {
					return notify.Panic(fmt.Errorf("cannot onkey and call (%T) %s", playOrEval, core.Storex(playOrEval)))
				}
			}
			err := ctx.Device().OnKey(ctx, key.DeviceID(), key.Channel(), key.Note(), getHasValue(playOrEval))
			if err != nil {
				return notify.Panic(fmt.Errorf("cannot install onkey because error:%v", err))
			}
			return nil
		}}

	eval["device"] = Function{
		Title:         "MIDI device selector",
		Description:   "select a MIDI device from the available device IDs; must become before channel",
		ControlsAudio: true,
		Prefix:        "dev",
		Template:      `device(${1:number},${2:sequenceable})`,
		Samples:       `device(1,channel(2,sequence('c2 e3'))) // plays on connected device 1 through MIDI channel 2`,
		Func: func(deviceID interface{}, m interface{}) interface{} {
			seq, ok := getSequenceable(m)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot decorate with device (%T) %v", m, m))
			}
			return core.NewDeviceSelector(seq, getHasValue(deviceID))
		}}

	eval["interval"] = Function{
		Title:       "Interval creator",
		Description: "create an integer repeating interval (from,to,by,method). Default method is 'repeat', Use next() to get a new integer",
		Prefix:      "int",
		Template:    `interval(${1:from},${2:to},${3:by})`,
		Samples: `int1 = interval(-2,4,1)
lp_cdef = loop(transpose(int1,sequence('c d e f')), next(int1))`,
		IsComposer: true,
		Func: func(from, to, by interface{}) *core.Interval {
			return core.NewInterval(core.ToHasValue(from), core.ToHasValue(to), core.ToHasValue(by), core.RepeatFromTo)
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
			return op.NewResequencer(s, core.ToHasValue(pattern))
		}}

	eval["notemap"] = Function{
		Title:       "Note Map creator",
		Description: "creates a mapper of notes by index (1-based) or using dots (.) and bangs (!)",
		Template:    `notemap('${1:space-separated-1-based-indices-or-dots-and-bangs}',${2:has-note})`,
		IsComposer:  true,
		Samples: `m1 = notemap('..!..!..!', note('c2'))
m2 = notemap('3 6 9', octave(-1,note('d2')))`,
		Func: func(indices string, note interface{}) interface{} {
			m, err := op.NewNoteMap(indices, getHasValue(note))
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

	eval["if"] = Function{
		Title:       "Conditional operator",
		Description: "Supports conditions with operators on numbers: <,<=,>,>=,!=,==",
		Samples:     ``,
		Func: func(c interface{}, thenelse ...interface{}) interface{} {
			if len(thenelse) == 0 {
				notify.Panic(fmt.Errorf("requires at least a <then>"))
			}
			if len(thenelse) > 2 {
				notify.Panic(fmt.Errorf("requires at most a <then> and an <else>"))
			}
			thenarg, ok := getSequenceable(thenelse[0])
			if !ok {
				notify.Panic(fmt.Errorf("cannot conditional use (%T) %v", thenelse[0], thenelse[0]))
			}
			ifop := op.IfCondition{Condition: getHasValue(c), Then: thenarg, Else: core.EmptySequence}
			if len(thenelse) == 2 {
				elsearg, ok := getSequenceable(thenelse[1])
				if !ok {
					notify.Panic(fmt.Errorf("cannot conditional use (%T) %v", thenelse[1], thenelse[1]))
				}
				ifop.Else = elsearg
			}
			return ifop
		},
	}

	eval["value"] = Function{
		Title:       "Value operator",
		Description: "returns the current value of a variable",
		Func: func(v interface{}) interface{} {
			return core.ValueFunction{
				StoreString: fmt.Sprintf("value(%s)", core.Storex(v)),
				Function: func() interface{} {
					return core.ValueOf(v)
				},
			}
		},
	}

	eval["index"] = Function{
		Title:       "Index operator",
		Description: "returns the current index of an object (e.g. iterator,interval,repeat)",
		Func: func(v interface{}) interface{} {
			return core.ValueFunction{
				StoreString: fmt.Sprintf("index(%s)", core.Storex(v)),
				Function: func() interface{} {
					return core.IndexOf(v)
				},
			}
		},
	}

	eval["next"] = Function{
		Title: "Next operator",
		Description: `is used to produce the next value in a generator such as random, iterator and interval.
The function itself does not return the value; use the generator for that.`,
		Samples: `i = interval(-4,4,2)
pi = transpose(i,sequence('c d e f g a b')) // current value of "i" is used
lp_pi = loop(pi,next(i)) // "i" will advance to the next value
begin(lp_pi)`,
		Func: func(v interface{}) interface{} {
			return core.Nexter{Target: getHasValue(v)}
		}}

	eval["export"] = Function{
		Title:       "Export command",
		Description: `writes a multi-track MIDI file`,
		Template:    `export(${1:filename},${2:sequenceable})`,
		Samples:     `export('myMelody-v1',myObject)`,
		Func: func(filename string, m interface{}) interface{} {
			if !ctx.Capabilities().ExportMIDI {
				return notify.NewWarningf("export MIDI not available")
			}
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
			return file.Export(filename, getValue(m), ctx.Control().BPM(), ctx.Control().BIAB())
		}}

	eval["trim"] = Function{
		Title:       "Trim notes|groups from start or end",
		Description: `create a new sequence object with notes trimmed at the start or/and at the end.`,
		Template:    `trim(${1:remove-from-start},${2:remove-from-end},${3:object})`,
		Samples:     `t = trim(1,2,sequence('c d e f a') // d e`,
		Func: func(skipStart, skipEnd, object interface{}) interface{} {
			s, ok := getSequenceable(object)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot trim non-sequenceable"))
			}
			return op.Trim{
				Start:  getHasValue(skipStart),
				End:    getHasValue(skipEnd),
				Target: s}
		}}

	eval["replace"] = Function{
		Title:       "Replace operator",
		Description: `replaces all occurrences of one musical object with another object for a given composed musical object`,
		Template:    `replace(${1:target},${2:from},${3:to})`,
		Samples: `c = note('c')
d = note('d')
pitchA = transpose(1,c)
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

	registerFunction(eval, "set", Function{
		Title:         "Change a setting",
		Description:   "Generic function to change a default setting",
		ControlsAudio: true,
		Template:      "set(${1:setting-name},${2:setting-value})",
		Samples: `set('midi.in',1) // default MIDI input device is 1
set('midi.in.channel',2,10) // default MIDI channel for device 2 is 10
set('midi.out',3) // default MIDI output device is 3`,
		Func: func(settingName string, settingValues ...interface{}) interface{} {
			if err := ctx.Device().HandleSetting(settingName, settingValues); err != nil {
				notify.Errorf("%v", err)
			}
			return nil
		},
	})

	registerFunction(eval, "listen", Function{
		Title:       "Start a MIDI listener",
		Description: "Listen for note(s) from a device and call a playable function to handle",
		Template:    "listen(${1:variable-or-device-selector},${2:function})",
		Samples: `rec = note('c') // define a variable "rec" with a initial object ; this is a place holder
fun = play(rec) // define the playable function to call when notes are received ; loop and print are also possible
listen(rec,fun) // start a listener for notes from default input device, store it in "rec" and call "fun"
listen(device(1,rec),fun) // start a listener for notes from input device 1`,
		Func: func(varOrDeviceSelector interface{}, function interface{}) interface{} {
			_, ok := getValue(function).(core.Evaluatable)
			if !ok {
				return notify.Panic(fmt.Errorf("cannot listen and call (%T) %s", function, core.Storex(function)))
			}
			var injectable variable
			deviceID, _ := ctx.Device().DefaultDeviceIDs()
			if ds, ok := varOrDeviceSelector.(core.DeviceSelector); ok {
				deviceID = ds.DeviceID()
				first := ds.Target
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
			// use function as HasValue and not the Evaluatable to allow redefinition of the callback function in the script
			return control.NewListen(ctx, deviceID, injectable.Name, getHasValue(function))
		},
	})

	registerFunction(eval, "onoff", Function{
		Title:         "Note ON/OFF switch",
		Description:   "play will send MIDI Note On, stop will send MIDI Note Off",
		Template:      "onoff(${2:note})",
		ControlsAudio: true,
		Samples: `// latch example
// if C4 is hit on input device 1 
// then play (sustain) key E on the default output device. 
// A second hit of C4 will stop it

onkey('c4',onoff('e')) // uses default input and default output MIDI device`,
		Func: func(noteSource string) interface{} {
			// Simple first
			_, deviceID := ctx.Device().DefaultDeviceIDs()
			note, err := core.ParseNote(noteSource)
			if err != nil {
				notify.Panic(err)
			}
			return control.NewOnOff(deviceID, 1, note)
		},
	})

	return eval
}
