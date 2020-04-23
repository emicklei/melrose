---
title: Melrose Command Line Interface (CLI)
---

[How to install](install.html)
[Usage](cli.html)
[Language](dsl.html)

# Melrose

The command-line tool `melrose` is a Read–Eval–Print Loop (REPL) that produces or consumes MIDI. 
By entering statements using the [language](dsl.html), `melrose` will send out MIDI messages to any connected [DAW](daw.html).

Commands to control the program itself are prefix with a colon `:`.
With `:h` you get list of known functions and commands.

## help

    𝄞 :h

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

        :h --- show help, optional on a command or function
        :k --- stop all running Loops
        :l --- load memory from disk, optional use given filename
        :m --- MIDI settings
        :q --- quit
        :s --- save memory to disk, optional use given filename
        :v --- show variables, optional filter on given prefix


## line editing

The following line editing commands are supported on platforms and terminals
that Melrose supports:

Keystroke    | Action
---------    | ------
Tab          | Next completion
Shift-Tab    | (after Tab) Previous completion
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

## API server

Melrose starts a HTTP server on port 8118 and evaluates statements on `POST /v1/statements`.
The port can be changed to e.g. 8000 with the program option `-http :8000`.
This server is used by the `Melrose Plugin for Visual Studio Code`.