---
title: Melrōse Language
---

$$menu

# Language

### expressions

Musical objects are created, composed and played using the <strong>melrõse</strong> tool by evaluating expressions.
Expressions use any of the predefined functions (creation,composition,audio control).
By assigning an expression to a variable name, you can use that expression by its name to compose other objects.

###variables

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
- <a href="#duration">duration</a>
- <a href="#export">export</a>
- <a href="#flatten">flatten</a>
- <a href="#interval">interval</a>
- <a href="#join">join</a>
- <a href="#next">next</a>
- <a href="#notemap">notemap</a>
- <a href="#notemerge">notemerge</a>
- <a href="#octave">octave</a>
- <a href="#octavemap">octavemap</a>
- <a href="#on">on</a>
- <a href="#parallel">parallel</a>
- <a href="#pitch">pitch</a>
- <a href="#random">random</a>
- <a href="#repeat">repeat</a>
- <a href="#replace">replace</a>
- <a href="#reverse">reverse</a>
- <a href="#sequencemap">sequencemap</a>
- <a href="#serial">serial</a>
- <a href="#undynamic">undynamic</a>
- <a href="#watch">watch</a>

## Audio control functions

- <a href="#begin">begin</a>
- <a href="#biab">biab</a>
- <a href="#bpm">bpm</a>
- <a href="#channel">channel</a>
- <a href="#end">end</a>
- <a href="#go">go</a>
- <a href="#loop">loop</a>
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

```

### bpm<a name="bpm"></a>
Set the Beats Per Minute [1..300]; default is 120.

```javascript

```

### channel<a name="channel"></a>
Select a MIDI channel, must be in [1..16]; must be a top-level operator.

```javascript
channel(2,sequence('C2 E3')) // plays on instrument connected to MIDI channel 2
```

### chord<a name="chord"></a>
Create a Chord from its string <a href="/melrose/melrose/notations.html#chord-not">notation</a>.

```javascript
chord('C#5/m/1')

chord('G/M/2')
```

### duration<a name="duration"></a>
Creates a new modified musical object for which the duration of all notes are changed.
The first parameter controls the length (duration) of the note.
If the parameter is greater than 0 then the note duration is set to a fixed value, e.g. 4=quarter,1=whole.
If the parameter is less than 1 then the note duration is scaled with a value, e.g. 0.5 will make a quarter ¼ into an eight ⅛
.

```javascript
duration(8,sequence('E F')) // => ⅛E ⅛F , absolute change

duration(0.5,sequence('8C 8G')) // => C G , factor change
```

### end<a name="end"></a>
End running loop(s). Ignore if it was stopped.

```javascript
l1 = loop(sequence('C E G))

begin(l1) // end(l1)
```

### export<a name="export"></a>
Writes a multi-track MIDI file.

```javascript
export("myMelody-v1",myObject)
```

### flatten<a name="flatten"></a>
Flatten (ungroup) all operations on a musical object to a new sequence.

```javascript
flatten(sequence('(C E G) B')) // => C E G B
```

### go<a name="go"></a>
Play all musical objects together in the background (do not wait for completion).

```javascript
go(s1,s1,s3) // play s1 and s2 and s3 simultaneously
```

### interval<a name="interval"></a>
Create an integer repeating interval (from,to,by,method). Default method is 'repeat', Use next() to get a new integer.

```javascript
int1 = interval(-2,4,1)

lp_cdef = loop(pitch(int1,sequence('C D E F')), next(int1))
```

### join<a name="join"></a>
When played, each musical object is played in sequence.

```javascript
a = chord('A')

b = sequence('(C E G))

ab = join(a,b)
```

### loop<a name="loop"></a>
Create a new loop from one or more musical objects; must be assigned to a variable.

```javascript
cb = sequence('C D E F G A B')

lp_cb = loop(cb,reverse(cb))
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

### next<a name="next"></a>
Is used to seed the next value in a generator such as random and interval.

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

note('2e#.--')
```

### notemap<a name="notemap"></a>
Creates a mapper of notes by index (1-based) or using dots (.) and bangs (!).

```javascript

```

### notemerge<a name="notemerge"></a>
Merges multiple notemaps into one sequenc.

```javascript

```

### octave<a name="octave"></a>
Changes the pitch of notes by steps of 12 semitones for one or more musical objects.

```javascript
octave(1,sequence('C D')) // => C5 D5
```

### octavemap<a name="octavemap"></a>
Create a sequence with notes for which order and the octaves are changed.

```javascript
octavemap('1:-1,2:0,3:1',chord('C')) // => (C3 E G5)
```

### on<a name="on"></a>
Puts a musical object on a track to start at a specific bar.

```javascript

```

### parallel<a name="parallel"></a>
Create a new sequence in which all notes of a musical object are grouped.

```javascript
parallel(sequence('C D E')) // => (C D E)
```

### pitch<a name="pitch"></a>
Change the pitch with a delta of semitones.

```javascript
pitch(-1,sequence('C D E'))

p = interval(-4,4,1)

pitch(p,note('C'))
```

### play<a name="play"></a>
Play all musical objects.

```javascript
play(s1,s2,s3) // play s3 after s2 after s1
```

### progression<a name="progression"></a>
Create a Chord progression using this <a href="/melrose/notations.html#progression-not">format</a>.

```javascript
progression('E F') // => (E A♭ B) (F A C5)

progression('(C D)') // => (C E G D G♭ A)
```

### random<a name="random"></a>
Create a random integer generator. Use next() to seed a new integer.

```javascript
num = random(1,10)

next(num)
```

### record<a name="record"></a>
Creates a recorded sequence of notes from the current MIDI input device.

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

pitchA = pitch(1,c)

pitchD = replace(pitchA, c, note('d'))
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

sequence('(C D E)')
```

### sequencemap<a name="sequencemap"></a>
Creates a mapper of sequence notes by index (1-based).

```javascript
s1 = sequence('C D E F G A B')

i1 = sequencemap('6 5 4 3 2 1',s1) // => B A G F E D
```

### serial<a name="serial"></a>
Serialise any grouping of notes from one or more musical objects.

```javascript
serial(chord('E')) // => E G B

serial(sequence('(C D)'),note('E')) // => C D E
```

### undynamic<a name="undynamic"></a>
Set the dymamic to normal for all notes in a musical object.

```javascript
undynamic('A+ B++ C-- D-') // =>  A B C D
```

### watch<a name="watch"></a>
Create a Note.

```javascript

```



##### generated by dsl-md.go
