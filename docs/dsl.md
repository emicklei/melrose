---
title: Melrōse Language
---

[Home](index.html)
[Usage](cli.html)
[Language](dsl.html)
[DAW](daw.html)
[Install](install.html)

### variables

Variable names must start with a non-digit character and can have zero or more characters in [a-z A-Z _ 0-9].
An assigment "=" is used to create a variable.
To delete a variable, assign it to the special value "nil".

### comment

Use "//" to add comment, either on a new line or and the end of a statement.

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
- <a href="#flatten">flatten</a>
- <a href="#interval">interval</a>
- <a href="#join">join</a>
- <a href="#octave">octave</a>
- <a href="#octavemap">octavemap</a>
- <a href="#parallel">parallel</a>
- <a href="#pitch">pitch</a>
- <a href="#repeat">repeat</a>
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
- <a href="#echo">echo</a>
- <a href="#end">end</a>
- <a href="#go">go</a>
- <a href="#loop">loop</a>
- <a href="#onbar">onbar</a>
- <a href="#play">play</a>
- <a href="#record">record</a>


### at<a name="at"></a>
Create an index getter (1-based) to select a musical object.

	at(1,scale('E/m')) // => E

### begin<a name="begin"></a>
Begin loop(s). Ignore if it was running.

	l1 = loop(sequence('C D E F G A B'))

	end(l1)

	begin(l1)

### biab<a name="biab"></a>
Set the Beats in a Bar [1..6]; default is 4.

	

### bpm<a name="bpm"></a>
Set the Beats Per Minute [1..300]; default is 120.

	

### channel<a name="channel"></a>
Select a MIDI channel, must be in [0..16].

	channel(2,sequence('C2 E3') // plays on instrument connected to MIDI channel 2'

### chord<a name="chord"></a>
Create a Chord from its string <a href="/melrose/index.html#chord-not">notation</a>.

	chord('C#5/m/1')

	chord('G/M/2')

### duration<a name="duration"></a>
Creates a new modified musical object for which the duration of all notes are changed.
The first parameter controls the length (duration) of the note.
If the parameter is greater than 0 then the note duration is set to a fixed value, e.g. 4=quarter,1=whole.
If the parameter is less than 1 then the note duration is scaled with a value, e.g. 0.5 will make a quarter ¼ into an eight ⅛
.

	duration(8,sequence('E F')) // => ⅛E ⅛F , absolute change

	duration(0.5,sequence('8C 8G')) // => C G , factor change

### echo<a name="echo"></a>
Echo the notes being played; default is true.

	echo(false)

### end<a name="end"></a>
End running loop(s). Ignore if it was stopped.

	l1 = loop(sequence('C E G))

	end(l1)

### flatten<a name="flatten"></a>
Flatten all operations on a musical object to a new sequence.

	flatten(sequence('(C E G) B')) // => C E G B

### go<a name="go"></a>
Play all musical objects in parallel.

	go(s1,s1,s3) // play s1 and s2 and s3 simultaneously

### interval<a name="interval"></a>
Create an integer repeating interval (from,to,by,method). Default method is 'repeat', Use next() to get a new integer.

	i1 = interval(-2,4,1)

	l1 = loop(pitch(i1,sequence('C D E F')), next(i1))

### join<a name="join"></a>
Join two or more musical objects.

	

### loop<a name="loop"></a>
Create a new loop from one or more objects.

	cb = sequence('C D E F G A B')

	lp_cb = loop(cb,reverse(cb))

### midi<a name="midi"></a>
Create a Note.

	midi(52,80) // => E3+

### note<a name="note"></a>
Create a Note  from its string <a href="/index.html#note-not">notation</a>.

	note('E')

	note('2E#.--')

### octave<a name="octave"></a>
Changes the pitch of notes by steps of 12 semitones.

	octave(1,sequence('C D')) // => C5 D5

### octavemap<a name="octavemap"></a>
Create a sequence with notes for which order and the octaves are changed.

	octavemap('1:-1,2:0,3:1',chord('C')) // => (C3 E G5)

### onbar<a name="onbar"></a>
.

	onbar(1,sequence('C D E')) // => immediately play C D E

### parallel<a name="parallel"></a>
Create a new sequence in which all notes of a musical object are synched in time.

	parallel(sequence('C D E')) // => (C D E)

### pitch<a name="pitch"></a>
Change the pitch with a delta of semitones.

	pitch(-1,sequence('C D E'))

	p = interval(-4,4,1)

	pitch(p,note('C'))

### play<a name="play"></a>
Play musical objects such as Note,Chord,Sequence,...

	play(s1,s2,s3) // play s3 after s2 after s1

### progression<a name="progression"></a>
.

	progression('E F') // => (E A♭ B) (F A C5)

	progression('(C D)') // => (C E G D G♭ A)

### record<a name="record"></a>
Creates a recorded sequence of notes from a MIDI device.

	r = record(1,5) // record notes played on device ID=1 and stop recording after 5 seconds

	s = r.Sequence()

### repeat<a name="repeat"></a>
Repeat the musical object a number of times.

	repeat(4,sequence('C D E'))

### reverse<a name="reverse"></a>
Reverse the (groups of) notes in a sequence.

	reverse(chord('A'))

### scale<a name="scale"></a>
Create a Scale using a starting Note and type indicator (Major,minor).

	scale(1,'E/m') // => E F G A B C5 D5

### sequence<a name="sequence"></a>
Create a Sequence from (space separated) notes.

	sequence('C D E')

	sequence('(C D E)')

### sequencemap<a name="sequencemap"></a>
Create a Mapper of sequence notes by index (1-based).

	s1 = sequence('C D E F G A B')

	i1 = sequencemap('6 5 4 3 2 1',s1) // => B A G F E D

### serial<a name="serial"></a>
Serialise any grouping of notes in one or more musical objects.

	serial(chord('E')) // => E G B

	serial(sequence('(C D)'),note('E')) // => C D E

### undynamic<a name="undynamic"></a>
Undynamic all the notes in a musical object.

	undynamic('A+ B++ C-- D-') // =>  A B C D

### watch<a name="watch"></a>
Create a Note.

	



##### generated by dsl-md.go
