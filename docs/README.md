# Melrose

[How to install](install.html)
[Using the CLI](cli.html)

The basic musical objects in Melrose are:

- Note
- Sequence
- Chord
- Scale

Musical objects can be composed using:

- Repeat
- Pitch
- Reverse
- Rotate
- Join
- Parallel (chord)
- Serial (arpeggio)
- Undynamic
- IndexMapper
- Loop

Parameters of compositions can be:

- Scalar values (integer, float)
- Interval
- Variable to a scalar or interval

## Notations

### Note notation

| Notation | Alternative | Description 
|----------|-------|-------------
| C4       | ¼C,C,c  | quarter C octave 4 
| 2E5      | ½E5,½e5 | Halftone (2 x ¼) E octave 5
| 1C       |        | Full tone C octave 4
| F#       | F♯,f♯  | F sharp
| G_       | G♭    | G flat
| G.       | G.    | duration x 1.5 
| =        | =     | quarter rest
| 2=       | ½=    | half rest
| 1=       | 1=    | full rest
| D+       | d+    | quarter D octave 4 MezzoForte
| 16E#.--  | 16e♯.-- | sixteenth E sharp duration x 1.5 Piano

### Note dynamics

| Notation    | Description
|-------------|---
| --- |Pianissimo
| --	|Piano
| -	  |MezzoPiano
| +	  |MezzoForte
| ++	|Forte
| +++ |Fortissimo

### Sequence notation

| Notation    | Description
|-------------|---
| C D E F       | 4 quarter tones
| [C E] [d5 f5] | 2 doublets
| [1C 1E 1G]    | C Chord

### Chord notation

| Notation    | Description
|-------------|---
| C#5/m/2     | C sharp triad, Octave 5, Minor, 2nd inversion

### Scale notation

| Notation    | Description
|-------------|---
| C5          | C major scale, Octave 5
| E/m         | E natural minor scale, Octave 4

## Melrose Language

### variables

Variable names must start with a non-digit character and can zero or more characters in `a-z A-Z _ 0-9`.
An assigment `=` is used to create a Variable.
To delete a variable, assign it to the special value `nil`.

### composition functions

Functions create or augment musical objects. 
Objects cannot be changed after creation.
Each function returns a new object or an object wrapped in a function.

            chord --- create a Chord
          flatten --- flatten all operations on a musical object to a new sequence
         indexmap --- create a Mapper of Notes by index (1-based)
         interval --- create an integer repeating interval (from,to,by)
             join --- join two or more musical objects
             note --- Note, e.g. C 2G#5. =
         parallel --- create a new sequence in which all notes of a musical object will be played in parallel
            pitch --- change the pitch with a delta of semitones
           repeat --- repeat the musical object a number of times
          reverse --- reverse the (groups of) notes in a sequence
         sequence --- create a Sequence from a string of notes
           serial --- serialise any parallelisation of notes in a musical object
        undynamic --- undynamic all the notes in a musical object

### audio functions

These functions control the audio device (playing, changing settings).

             bpm --- set the Beats Per Minute [1..300], default is 120
         channel --- select a MIDI channel, must be in [0..16]
            echo --- Echo the notes being played (default is true)
              go --- play all musical objects in parallel
            loop --- create a new loop
            play --- play musical objects such as Note,Chord,Sequence,...
          record --- creates a recorded sequence of notes from device ID and stop after T seconds of inactivity
             run --- start loop(s). Ignore if it was running.
            stop --- stop running loop(s). Ignore if it was stopped.
        velocity --- set the base velocity [1..127], default is 70

### comment

Use `//` at the start of a line to add comment.