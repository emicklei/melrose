# Melrose - programming of music melodies

[![Go Report Card](https://goreportcard.com/badge/github.com/emicklei/melrose)](https://goreportcard.com/report/github.com/emicklei/melrose)
[![GoDoc](https://godoc.org/github.com/emicklei/melrose?status.svg)](https://pkg.go.dev/github.com/emicklei/melrose?tab=doc)

## Note notation

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

### Note dynamics

	Pianissimo // ---
	Piano      // --
	MezzoPiano // -
	MezzoForte // +
	Forte      // ++
	Fortissimo // +++

## Sequence notation

| Notation    | Description 
|-------------|---
| C D E F       | 4 quarter tones
| [C E] [d5 f5] | 2 doublets


## Chord notation

| Notation    | Description 
|-------------|---
| C#5/m/2     | C sharp triad, Octave 5, Minor, 2nd inversion


# Melrose Live Coding

Using the command-line tool called `melrose` and a MIDI controlled synthesizer.

## language

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

## help

    𝄞 :h
      info: 
         chord --- create a Chord
       flatten --- flatten all operations on a musical object to a new sequence
          join --- join two or more musical objects
          note --- create a Note from a string
      parallel --- create a new sequence in which all notes of a musical object will be played in parallel
         pitch --- change the pitch with a delta of semitones
        repeat --- repeat the musical object a number of times
       reverse --- reverse the (groups of) notes in a sequence
      sequence --- create a Sequence from a string of notes
        serial --- serialise any parallelisation of notes in a musical object
     undynamic --- undynamic all the notes in a musical object

         bpm --- get or set the Beats Per Minute value [1..300], default is 120
          go --- play all musical objects in parallel
        play --- play a musical object such as Note,Chord,Sequence,...
      record --- creates a recorded sequence of notes from device ID and stop after T seconds of inactivity

      :h --- show help, optional on a command or function
      :l --- load memory from disk, optional use given filename
      :m --- show MIDI information
      :q --- quit
      :s --- save memory to disk, optional use given filename
      :v --- show variables, optional filter on given prefix


## line editing

The following line editing commands are supported on platforms and terminals
that Melrose supports:

Keystroke    | Action
---------    | ------
Ctrl-A, Home | Move cursor to beginning of line
Ctrl-E, End  | Move cursor to end of line
Ctrl-B, Left | Move cursor one character left
Ctrl-F, Right| Move cursor one character right
Ctrl-Left, Alt-B    | Move cursor to previous word
Ctrl-Right, Alt-F   | Move cursor to next word
Ctrl-D, Del  | (if line is *not* empty) Delete character under cursor
Ctrl-D       | (if line *is* empty) End of File - usually quits application
Ctrl-C       | Reset input (create new empty prompt)
Ctrl-L       | Clear screen (line is unmodified)
Ctrl-T       | Transpose previous character with current character
Ctrl-H, BackSpace | Delete character before cursor
Ctrl-W, Alt-BackSpace | Delete word leading up to cursor
Alt-D        | Delete word following cursor
Ctrl-K       | Delete from cursor to end of line
Ctrl-U       | Delete from start of line to cursor
Ctrl-P, Up   | Previous match from history
Ctrl-N, Down | Next match from history
Ctrl-R       | Reverse Search history (Ctrl-S forward, Ctrl-G cancel)
Ctrl-Y       | Paste from Yank buffer (Alt-Y to paste next yank instead)
Tab          | Next completion
Shift-Tab    | (after Tab) Previous completion


Software is licensed under [Apache 2.0 license](LICENSE).
(c) 2014-2020 http://ernestmicklei.com 