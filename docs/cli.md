---
title: Melrōse Command Line Interface (CLI)
---

[Home](index.html)
[Usage](cli.html)
[Language](dsl.html)
[DAW](daw.html)
[Install](install.html)

# Melrōse program

The program `melrōse` is a Read–Eval–Print Loop (REPL) that produces or consumes MIDI. 
By entering statements using the [language](dsl.html), `melrōse` will send out MIDI messages to any connected [DAW](daw.html).
Although it is possible to program directly using the command line interface of `melrōse`, it is much more convenient to use the Visual Studio Code editor with the [Melrose Plugin](vsc.html).

### control
Commands to control the program itself are prefix with a colon `:`.
With `:h` you get list of known functions and commands.

### line editing

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

Melrōse starts a HTTP server on port 8118 and evaluates programs on `POST /v1/statements` providing the source as the payload (HTTP Body).
This server is used by the [Melrōse Plugin for Visual Studio Code](https://github.com/emicklei/melrōse-for-vscode).

### HTTP response

#### 200 OK

If the request was successful processed then the response looks like:

  {
    "type": "melrose.Sequence",
    "object: { ... }
  }

#### 500 Internal Server Error

If the request could not be processed then the response looks like:

  {
    "type": "errors.Error",
    "message": "unknown function",
    "line": 1,
    "column": 1
  }

#### 400 Bad Request

If the request is malformed then the response will have the error message.

### HTTP port

The port can be changed to e.g. 8000 with the program option `-http :8000`.

### tracing

If the HTTP URL has the query parameter `trace=true` then `melrōse` will produce extra logging.

### play

If the HTTP URL has the query parameter `action=play` then `melrōse` will try to play the result of the selected expression(s).

### begin

If the HTTP URL has the query parameter `action=begin` then `melrōse` will try to `begin` the loop of the selected expression.

### end

If the HTTP URL has the query parameter `action=end` then `melrōse` will try to `end` the loop of the selected expression.

### inspecting

If the HTTP URL has the query parameter `action=inspect` then `melrōse` will print inspection details.