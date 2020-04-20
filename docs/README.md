# Melrose

[How to install](install.html)
[Using the CLI](cli.html)

The basic musical objects in Melrose are:

- Note
- Sequence
- Chord

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

## Melrose Language

### variables

Variable names must start with a non-digit character and can zero or more characters in `a-z A-Z _ 0-9`.
An assigment `=` is used to create a Variable.
To delete a variable, assign it to the special value `nil`.

### composition functions

Functions create or augment musical objects. 
Objects cannot be changed after creation.
Each function returns a new object or an object wrapped in a function.

### audio functions

These functions control the audio device (playing, changing settings).

### comment

Use `//` at the start of a line to add comment.