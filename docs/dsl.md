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
- <a href="#duration">duration</a>
- <a href="#dynamic">dynamic</a>
- <a href="#dynamicmap">dynamicmap</a>
- <a href="#export">export</a>
- <a href="#fraction">fraction</a>
- <a href="#group">group</a>
- <a href="#import">import</a>
- <a href="#interval">interval</a>
- <a href="#iterator">iterator</a>
- <a href="#join">join</a>
- <a href="#joinmap">joinmap</a>
- <a href="#listen">listen</a>
- <a href="#merge">merge</a>
- <a href="#midi_send">midi_send</a>
- <a href="#next">next</a>
- <a href="#notemap">notemap</a>
- <a href="#octave">octave</a>
- <a href="#octavemap">octavemap</a>
- <a href="#onbar">onbar</a>
- <a href="#pitch">pitch</a>
- <a href="#pitchmap">pitchmap</a>
- <a href="#print">print</a>
- <a href="#random">random</a>
- <a href="#repeat">repeat</a>
- <a href="#replace">replace</a>
- <a href="#resequence">resequence</a>
- <a href="#reverse">reverse</a>
- <a href="#stretch">stretch</a>
- <a href="#track">track</a>
- <a href="#undynamic">undynamic</a>
- <a href="#ungroup">ungroup</a>

## Audio control functions

- <a href="#begin">begin</a>
- <a href="#biab">biab</a>
- <a href="#bpm">bpm</a>
- <a href="#channel">channel</a>
- <a href="#device">device</a>
- <a href="#end">end</a>
- <a href="#loop">loop</a>
- <a href="#multitrack">multitrack</a>
- <a href="#play">play</a>
- <a href="#sync">sync</a>


### at<a name="at"></a>
Create an index getter (1-based) to select a musical object.

> at(index,object)

#### examples	
```javascript
at(1,scale('e/m')) // => E
```

### begin<a name="begin"></a>
Begin loop(s). Ignore if it was running.

> begin(loop)

#### examples	
```javascript
lp_cb = loop(sequence('C D E F G A B'))

begin(lp_cb) // end(lp_cb)
```

### biab<a name="biab"></a>
Set the Beats in a Bar; default is 4.

> biab(beats-in-a-bar)

#### examples	
```javascript
biab(4)
```

### bpm<a name="bpm"></a>
Set the Beats Per Minute (BPM) [1..300]; default is 120.

> bpm(beats-per-minute)

#### examples	
```javascript
bpm(90)

speedup = iterator(80,100,120,140)

l = loop(bpm(speedup),sequence('c e g'),next(speedup))
```

### channel<a name="channel"></a>
Select a MIDI channel, must be in [1..16]; must be a top-level operator.

> channel(number,sequenceable)

#### examples	
```javascript
channel(2,note('g3'), sequence('c2 e3')) // plays on instrument connected to MIDI channel 2
```

### chord<a name="chord"></a>
Create a Chord from its string <a href="/melrose/notations.html#chord-not">notation</a>.

> chord('note')

#### examples	
```javascript
chord('c#5/m/1')

chord('g/M/2') // Major G second inversion
```

### device<a name="device"></a>
Select a MIDI device from the available device IDs; must become before channel.

> device(number,sequenceable)

#### examples	
```javascript
device(1,channel(2,sequence('c2 e3'), note('g3'))) // plays on connected device 1 through MIDI channel 2
```

### duration<a name="duration"></a>
Computes the duration of the object using the current BPM.

> duration(object)

#### examples	
```javascript
duration(note('c'))
```

### dynamic<a name="dynamic"></a>
Creates a new modified musical object for which the dynamics of all notes are changed.
	The first parameter controls the emphasis the note, e.g. + (mezzoforte,mf), -- (piano,p).
	.

> dynamic(emphasis,object)

#### examples	
```javascript
dynamic('++',sequence('e f')) // => E++ F++
```

### dynamicmap<a name="dynamicmap"></a>
Changes the dynamic of notes from a musical object. 1-index-based mapping.

> dynamicmap('mapping',object)

#### examples	
```javascript
dynamicmap('1:++,2:--',sequence('e f')) // => E++ F--

dynamicmap('2:o,1:++,2:--,1:++', sequence('a b') // => B A++ B-- A++
```

### end<a name="end"></a>
End running loop(s) or listener(s). Ignore if it was stopped.

> end(control)

#### examples	
```javascript
l1 = loop(sequence('c e g'))

begin(l1)

end(l1)

end() // stop all playables
```

### export<a name="export"></a>
Writes a multi-track MIDI file.

> export(filename,sequenceable)

#### examples	
```javascript
export('myMelody-v1',myObject)
```

### fraction<a name="fraction"></a>
Creates a new object for which the fraction of duration of all notes are changed.
The first parameter controls the fraction of the note, e.g. 1 = whole, 2 = half, 4 = quarter, 8 = eight, 16 = sixteenth.
Fraction can also be an exact float value between 0 and 1.
.

> fraction(object,object)

#### examples	
```javascript
fraction(8,sequence('e f')) // => ⅛E ⅛F , shorten the notes from quarter to eight
```

### group<a name="group"></a>
Create a new sequence in which all notes of a musical object are grouped.

> group(sequenceable)

#### examples	
```javascript
group(sequence('c d e')) // => (C D E)
```

### import<a name="import"></a>
Evaluate all the statements from another file.

> import(filename)

#### examples	
```javascript
import('drumpatterns.mel')
```

### interval<a name="interval"></a>
Create an integer repeating interval (from,to,by,method). Default method is 'repeat', Use next() to get a new integer.

> interval(from,to,by)

#### examples	
```javascript
int1 = interval(-2,4,1)

lp_cdef = loop(pitch(int1,sequence('c d e f')), next(int1))
```

### iterator<a name="iterator"></a>
Iterator that has an array of constant values and evaluates to one. Use next() to increase and rotate the value.

> iterator(array-element)

#### examples	
```javascript
i = iterator(1,3,5,7,9)

p = pitch(i,note('c'))

lp = loop(p,next(i))

		
```

### join<a name="join"></a>
Joins one or more musical objects as one.

> join(first,second)

#### examples	
```javascript
a = chord('a')

b = sequence('(c e g)')

ab = join(a,b)
```

### joinmap<a name="joinmap"></a>
Creates a new join by mapping elements. 1-index-based mapping.

> joinmap('indices',join)

#### examples	
```javascript
j = join(note('c'), sequence('d e f'))

jm = joinmap('1 (2 3) 4',j)
```

### listen<a name="listen"></a>
Listen for note(s) from a device and call a playable function to handle.

> listen(device-id,variable,function)

#### examples	
```javascript
rec = note('c') // define a variable "rec" with a initial object ; this is a place holder

fun = play(rec) // define the playable function to call when notes are received ; loop and print are also possible

ear = listen(1,rec,fun) // start a listener for notes from device 1, store it "rec" and call "fun"
```

### loop<a name="loop"></a>
Create a new loop from one or more musical objects; must be assigned to a variable.

> lp_object = loop(object)

#### examples	
```javascript
cb = sequence('c d e f g a b')

lp_cb = loop(cb,reverse(cb))
```

### merge<a name="merge"></a>
Merges multiple sequences into one sequence.

> merge(sequenceable)

#### examples	
```javascript
m1 = notemap('..!..!..!', note('c2'))

m2 = notemap('4 7 10', note('d2'))

all = merge(m1,m2) // => = = C2 D2 = C2 D2 = C2 D2 = =
```

### midi<a name="midi"></a>
Create a Note from MIDI information and is typically used for drum sets.
The first parameter is a fraction {1,2,4,8,16} or a duration in milliseconds or a time.Duration.
Second parameter is the MIDI number and must be one of [0..127].
The third parameter is the velocity (~ loudness) and must be one of [0..127].

> midi(numberOrDuration,number,number)

#### examples	
```javascript
midi(500,52,80) // => 500ms E3+

midi(16,36,70) // => 16C2 (kick)
```

### midi_send<a name="midi_send"></a>
Sends a MIDI message with status, channel(ignore if < 1), 2nd byte and 3rd byte to an output device. Can be used as a musical object.

> midi_send(device-id,status,channel,2nd-byte,3rd-byte

#### examples	
```javascript
midi_send(1,0xB0,7,0x7B,0) // to device id 1, control change, all notes off in channel 7

midi_send(1,0xC0,2,1,0) // program change, select program 1 for channel 2

midi_send(2,0xB0,4,0,16) // control change, bank select 16 for channel 4

midi_send(3,0xB0,1,120,0) // control change, all notes off for channel 1
```

### multitrack<a name="multitrack"></a>
Create a multi-track object from zero or more tracks.

> multitrack(track)

#### examples	
```javascript
multitrack(track1,track2,track3) // 3 tracks in one multi-track object
```

### next<a name="next"></a>
Is used to produce the next value in a generator such as random, iterator and interval.

> 

#### examples	
```javascript
i = interval(-4,4,2)

pi = pitch(i,sequence('c d e f g a b')) // current value of "i" is used

lp_pi = loop(pi,next(i)) // "i" will advance to the next value

begin(lp_pi)
```

### note<a name="note"></a>
Create a Note using this <a href="/melrose/notations.html#note-not">format</a>.

> note('letter')

#### examples	
```javascript
note('e')

note('2.e#--')
```

### notemap<a name="notemap"></a>
Creates a mapper of notes by index (1-based) or using dots (.) and bangs (!).

> notemap('space-separated-1-based-indices-or-dots-and-bangs',note)

#### examples	
```javascript
m1 = notemap('..!..!..!', note('c2'))

m2 = notemap('3 6 9', note('d2'))
```

### octave<a name="octave"></a>
Change the pitch of notes by steps of 12 semitones for one or more musical objects.

> octave(offset,sequenceable)

#### examples	
```javascript
octave(1,sequence('c d')) // => C5 D5
```

### octavemap<a name="octavemap"></a>
Create a sequence with notes for which the order and the octaves are changed. 1-based indexing.

> octavemap('int2int',object)

#### examples	
```javascript
octavemap('1:-1,2:0,3:1',chord('c')) // => (C3 E G5)
```

### onbar<a name="onbar"></a>
Puts a musical object on a track to start at a specific bar.

> onbar(bar,object)

#### examples	
```javascript
tr = track("solo",2, onbar(1,soloSequence)) // 2 = channel
```

### pitch<a name="pitch"></a>
Change the pitch with a delta of semitones.

> pitch(semitones,sequenceable)

#### examples	
```javascript
pitch(-1,sequence('c d e'))

p = interval(-4,4,1)

pitch(p,note('c'))
```

### pitchmap<a name="pitchmap"></a>
Create a sequence with notes for which the order and the pitch are changed. 1-based indexing.

> pitchmap('int2int',object)

#### examples	
```javascript
pitchmap('1:-1,1:0,1:1',note('c')) // => B3 C D
```

### play<a name="play"></a>
Play all musical objects.

> play(sequenceable)

#### examples	
```javascript
play(s1,s2,s3) // play s3 after s2 after s1
```

### print<a name="print"></a>
Prints an object when evaluated (play,loop).

> 

#### examples	
```javascript

```

### progression<a name="progression"></a>
Create a Chord progression using this <a href="/melrose/notations.html#progression-not">format</a>.

> progression('chords')

#### examples	
```javascript
progression('e f') // => (E A♭ B) (F A C5)

progression('(c d)') // => (C E G D G♭ A)
```

### random<a name="random"></a>
Create a random integer generator. Use next() to generate a new integer.

> random(from,to)

#### examples	
```javascript
num = random(1,10)

next(num)
```

### repeat<a name="repeat"></a>
Repeat one or more musical objects a number of times.

> repeat(times,sequenceables)

#### examples	
```javascript
repeat(4,sequence('c d e'))
```

### replace<a name="replace"></a>
Replaces all occurrences of one musical object with another object for a given composed musical object.

> replace(target,from,to)

#### examples	
```javascript
c = note('c')

d = note('d')

pitchA = pitch(1,c)

pitchD = replace(pitchA, c, d) // c -> d in pitchA
```

### resequence<a name="resequence"></a>
Creates a modifier of sequence notes by index (1-based).

> resequence('space-separated-1-based-indices',sequenceable)

#### examples	
```javascript
s1 = sequence('C D E F G A B')

i1 = resequence('6 5 4 3 2 1',s1) // => B A G F E D

i2 = resequence('(6 5) 4 3 (2 1)',s1) // => (B A) G F (E D)
```

### reverse<a name="reverse"></a>
Reverse the (groups of) notes in a sequence.

> reverse(sequenceable)

#### examples	
```javascript
reverse(chord('a'))
```

### scale<a name="scale"></a>
Create a Scale using this <a href="/melrose/notations.html#scale-not">format</a>.

> scale(octaves,'scale-syntax')

#### examples	
```javascript
scale(1,'e/m') // => E F G A B C5 D5
```

### sequence<a name="sequence"></a>
Create a Sequence using this <a href="/melrose/notations.html#sequence-not">format</a>.

> sequence('space-separated-notes')

#### examples	
```javascript
sequence('c d e')

sequence('(8c d e)') => (⅛C D E)

sequence('c (d e f) a =')
```

### stretch<a name="stretch"></a>
Stretches the duration of musical object(s) with a factor. If the factor < 1 then duration is shortened.

> stretch(factor,object)

#### examples	
```javascript
stretch(2,note('c'))  // 2C

stretch(0.25,sequence('(c e g)'))  // (16C 16E 16G)

stretch(8,note('c'))  // C with length of 2 bars
```

### sync<a name="sync"></a>
Synchronise playing musical objects. Use play() for serial playing.

> sync(object)

#### examples	
```javascript
sync(s1,s2,s3) // play s1,s2 and s3 at the same time

sync(loop1,loop2) // begin loop2 at the next start of loop1
```

### track<a name="track"></a>
Create a named track for a given MIDI channel with a musical object.

> track('title',midi-channel, onbar(1,object))

#### examples	
```javascript
track("lullaby",1,onbar(2, sequence('c d e'))) // => a new track on MIDI channel 1 with sequence starting at bar
```

### undynamic<a name="undynamic"></a>
Set the dymamic to normal for all notes in a musical object.

> undynamic(sequenceable)

#### examples	
```javascript
undynamic('A+ B++ C-- D-') // =>  A B C D
```

### ungroup<a name="ungroup"></a>
Undo any grouping of notes from one or more musical objects.

> ungroup(sequenceable)

#### examples	
```javascript
ungroup(chord('e')) // => E G B

ungroup(sequence('(c d)'),note('e')) // => C D E
```



##### generated by dsl-md.go
