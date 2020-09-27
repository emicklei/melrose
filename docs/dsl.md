---
title: Melrōse Language
---

[Home](https://emicklei.github.io/melrose)

# Language

### expressions

Musical objects are created, composed and played using the <strong>melrõse</strong> tool by evaluating expressions.
Expressions use any of the predefined functions (creation,composition,audio control).
By assigning an expression to a variable name, you can use that expression by its name to compose other objects.

### variables

Variable names must start with a non-digit character and can have zero or more characters in [a-z A-Z _ 0-9].
An assignment "=" is used to create a variable.
To delete a variable, assign it to the special value "nil".

### comment

Use "//" to add comment, either on a new line or and the end of an expression.

## Creation functions

- <a href="#chord">chord</a>
- <a href="#midi">midi</a>
- <a href="#note">note</a>
- <a href="#progression">progression</a>
- <a href="#scale">scale</a>
- <a href="#sequence">sequence</a>

## Composition functions

- <a href="#at">at</a>
- <a href="#dynamic">dynamic</a>
- <a href="#export">export</a>
- <a href="#fraction">fraction</a>
- <a href="#group">group</a>
- <a href="#import">import</a>
- <a href="#interval">interval</a>
- <a href="#iterator">iterator</a>
- <a href="#join">join</a>
- <a href="#joinmap">joinmap</a>
- <a href="#merge">merge</a>
- <a href="#next">next</a>
- <a href="#notemap">notemap</a>
- <a href="#octave">octave</a>
- <a href="#octavemap">octavemap</a>
- <a href="#onbar">onbar</a>
- <a href="#pitch">pitch</a>
- <a href="#print">print</a>
- <a href="#random">random</a>
- <a href="#repeat">repeat</a>
- <a href="#replace">replace</a>
- <a href="#reverse">reverse</a>
- <a href="#sequencemap">sequencemap</a>
- <a href="#track">track</a>
- <a href="#undynamic">undynamic</a>
- <a href="#ungroup">ungroup</a>

## Audio control functions

- <a href="#begin">begin</a>
- <a href="#biab">biab</a>
- <a href="#bpm">bpm</a>
- <a href="#channel">channel</a>
- <a href="#end">end</a>
- <a href="#go">go</a>
- <a href="#loop">loop</a>
- <a href="#multi">multi</a>
- <a href="#play">play</a>
- <a href="#record">record</a>


### at<a name="at"></a>
Create an index getter (1-based) to select a musical object.

```javascript
at(1,scale('E/m')) // => E
```

### begin<a name="begin"></a>
Begin loop(s). Ignore if it was running.

```javascript
lp_cb = loop(sequence('C D E F G A B'))

begin(lp_cb) // end(lp_cb)
```

### biab<a name="biab"></a>
Set the Beats in a Bar [1..6]; default is 4.

```javascript
biab(4)
```

### bpm<a name="bpm"></a>
Set the Beats Per Minute (BPM) [1..300]; default is 120.

```javascript
bpm(90)
```

### channel<a name="channel"></a>
Select a MIDI channel, must be in [1..16]; must be a top-level operator.

```javascript
channel(2,sequence('C2 E3')) // plays on instrument connected to MIDI channel 2
```

### chord<a name="chord"></a>
Create a Chord from its string <a href="/melrose/notations.html#chord-not">notation</a>.

```javascript
chord('C#5/m/1')

chord('G/M/2')
```

### dynamic<a name="dynamic"></a>
Creates a new modified musical object for which the dynamics of all notes are changed.
	The first parameter controls the emphasis the note, e.g. + (mezzoforte,mf), -- (piano,p).
	.

```javascript
dynamic('++',sequence('E F')) // => E++ F++
```

### end<a name="end"></a>
End running loop(s). Ignore if it was stopped.

```javascript
l1 = loop(sequence('C E G'))

begin(l1) // end(l1)
```

### export<a name="export"></a>
Writes a multi-track MIDI file.

```javascript
export('myMelody-v1',myObject)
```

### fraction<a name="fraction"></a>
Creates a new object for which the fraction of duration of all notes are changed.
The first parameter controls the fraction of the note, e.g. 1=whole, 0.5 or 2 = half, 0.25 or 4 = quarter, 0.125 or 8 = eight, 0.0625 or 16 = sixteenth.
.

```javascript
fraction(8,sequence('e f')) // => ⅛E ⅛F , shorten the notes from quarter to eigth
```

### go<a name="go"></a>
Play all musical objects together in the background (do not wait for completion).

```javascript
go(s1,s1,s3) // play s1 and s2 and s3 simultaneously
```

### group<a name="group"></a>
Create a new sequence in which all notes of a musical object are grouped.

```javascript
group(sequence('C D E')) // => (C D E)
```

### import<a name="import"></a>
Evaluate all the statements from another file.

```javascript
import('drumpatterns.mel')
```

### interval<a name="interval"></a>
Create an integer repeating interval (from,to,by,method). Default method is 'repeat', Use next() to get a new integer.

```javascript
int1 = interval(-2,4,1)

lp_cdef = loop(pitch(int1,sequence('C D E F')), next(int1))
```

### iterator<a name="iterator"></a>
Iterator that has an array of constant values and evaluates to one. Use next() to increase and rotate the value.

```javascript
i = iterator(1,3,5,7,9)

		p = pitch(i,note('c'))

		lp = loop(p,next(i))

		
```

### join<a name="join"></a>
Joins two or more musical objects as one.

```javascript
a = chord('A')

b = sequence('(C E G)')

ab = join(a,b)
```

### joinmap<a name="joinmap"></a>
Creates a new join by mapping elements based on an index (1-based).

```javascript

```

### loop<a name="loop"></a>
Create a new loop from one or more musical objects; must be assigned to a variable.

```javascript
cb = sequence('C D E F G A B')

lp_cb = loop(cb,reverse(cb))
```

### merge<a name="merge"></a>
Merges multiple sequences into one sequence.

```javascript
m1 = notemap('..!..!..!', note('c2'))

m2 = notemap('4 7 10', note('d2'))

all = merge(m1,m2) // => = = C2 D2 = C2 D2 = C2 D2 = =
```

### midi<a name="midi"></a>
Create a Note from MIDI information and is typically used for drum sets.
The first parameter is the duration and must be one of {0.0625,0.125,0.25,0.5,1,2,4,8,16}.
A duration of 0.25 or 4 means create a quarter note.
Second parameter is the MIDI number and must be one of [0..127].
The third parameter is the velocity (~ loudness) and must be one of [0..127].

```javascript
midi(0.25,52,80) // => E3+

midi(16,36,70) // => 16C2 (kick)
```

### multi<a name="multi"></a>
Create a multi-track object from zero or more tracks.

```javascript
multi(track1,track2,track3) // one or more tracks in one multi-track object
```

### next<a name="next"></a>
Is used to produce the next value in a generator such as random and interval.

```javascript
i = interval(-4,4,2)

pi = pitch(i,sequence('C D E F G A B'))

lp_pi = loop(pi,next(i))

begin(lp_pi)
```

### note<a name="note"></a>
Create a Note using this <a href="/melrose/notations.html#note-not">format</a>.

```javascript
note('e')

note('2.e#--')
```

### notemap<a name="notemap"></a>
Creates a mapper of notes by index (1-based) or using dots (.) and bangs (!).

```javascript
m1 = notemap('..!..!..!', note('c2'))

m2 = notemap('3 6 9', note('d2'))
```

### octave<a name="octave"></a>
Change the pitch of notes by steps of 12 semitones for one or more musical objects.

```javascript
octave(1,sequence('C D')) // => C5 D5
```

### octavemap<a name="octavemap"></a>
Create a sequence with notes for which the order and the octaves are changed.

```javascript
octavemap('1:-1,2:0,3:1',chord('C')) // => (C3 E G5)
```

### onbar<a name="onbar"></a>
Puts a musical object on a track to start at a specific bar.

```javascript
tr = track("solo",2, onbar(1,soloSequence)) // 2 = channel
```

### pitch<a name="pitch"></a>
Change the pitch with a delta of semitones.

```javascript
pitch(-1,sequence('c d e'))

p = interval(-4,4,1)

pitch(p,note('c'))
```

### play<a name="play"></a>
Play all musical objects.

```javascript
play(s1,s2,s3) // play s3 after s2 after s1
```

### print<a name="print"></a>
Prints the musical object when evaluated (play,go,loop).

```javascript

```

### progression<a name="progression"></a>
Create a Chord progression using this <a href="/melrose/notations.html#progression-not">format</a>.

```javascript
progression('E F') // => (E A♭ B) (F A C5)

progression('(C D)') // => (C E G D G♭ A)
```

### random<a name="random"></a>
Create a random integer generator. Use next() to generate a new integer.

```javascript
num = random(1,10)

next(num)
```

### record<a name="record"></a>
Create a recorded sequence of notes from the current MIDI input device.

```javascript
r = record() // record notes played on the current input device and stop recording after 5 seconds

s = r.S() // returns the sequence of notes from the recording
```

### repeat<a name="repeat"></a>
Repeat the musical object a number of times.

```javascript
repeat(4,sequence('C D E'))
```

### replace<a name="replace"></a>
Replaces all occurrences of one musical object with another object for a given composed musical object.

```javascript
c = note('c')

d = note('d')

pitchA = pitch(1,c)

pitchD = replace(pitchA, c, d) // c -> d in pitchA
```

### reverse<a name="reverse"></a>
Reverse the (groups of) notes in a sequence.

```javascript
reverse(chord('A'))
```

### scale<a name="scale"></a>
Create a Scale using this <a href="/melrose/notations.html#scale-not">format</a>.

```javascript
scale(1,'E/m') // => E F G A B C5 D5
```

### sequence<a name="sequence"></a>
Create a Sequence using this <a href="/melrose/notations.html#sequence-not">format</a>.

```javascript
sequence('C D E')

sequence('(8C D E)') => (⅛C ⅛D ⅛E)

sequence('c (d e f) a =')
```

### sequencemap<a name="sequencemap"></a>
Creates a mapper of sequence notes by index (1-based).

```javascript
s1 = sequence('C D E F G A B')

i1 = sequencemap('6 5 4 3 2 1',s1) // => B A G F E D

i2 = sequencemap('(6 5) 4 3 (2 1)',s1) // => (B A) G F (E D)
```

### track<a name="track"></a>
Create a named track for a given MIDI channel with a musical object.

```javascript
track("lullaby",1,sequence('c d e')) // => a new track on MIDI channel 1
```

### undynamic<a name="undynamic"></a>
Set the dymamic to normal for all notes in a musical object.

```javascript
undynamic('A+ B++ C-- D-') // =>  A B C D
```

### ungroup<a name="ungroup"></a>
Undo any grouping of notes from one or more musical objects.

```javascript
ungroup(chord('E')) // => E G B

ungroup(sequence('(C D)'),note('E')) // => C D E
```



##### generated by dsl-md.go
