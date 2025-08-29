# Melrōse program

The program `melrōse` is a Read–Eval–Print Loop (REPL) that produces or consumes MIDI. 
By entering statements using the language, `melrōse` will send out MIDI messages to any connected DAW.

Although it is possible to program directly using the command line interface of `melrōse`, it is much more convenient to use the Visual Studio Code editor with the Melrose Plugin which uses the [HTTP API](http.md) of the same running program.

### program flags

You can start the program `melrōse` without any flags. 
You can use the following flags to change its behavior.

    -http <address>
        address on which to listen for HTTP requests (default ":8118")
    -d
        debug logging
    -log
        log file location

### CLI control

Commands to control the program itself are prefix with a colon `:`.
With `:h` you get the list of known commands.

    :h                    show help, optional on a command or function
    :v [prefix]           show variables, optional filter on given prefix
    :k                    stop all sound and loops
    :b                    beat settings
    -m                    MIDI settings
    -q                    quit
    -d                    toggle debug lines
    -p                    list all running
    -e                    echo MIDI

### CLI line editing

The following line editing commands are supported on platforms and terminals
that `melrōse` supports:

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


### special commands

    !<object> 

play the object

    <object>!

browse the notes of the object