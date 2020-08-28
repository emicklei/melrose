---
title: Melrōse Plugin for Visual Studio Code
---

[Home](https://emicklei.github.io/melrose)

## Plugin for Visual Studio Code

### install

See [instructions](install.html#plugin) on how to install this plugin after you have installed Visual Studio Code.

### edit

This editor extension works with `.mel` and `.melrose` files.
The syntax of the program uses the [Melrōse Lanugage](dsl.html).

### evaluate

To evaluate a single line statement or expression, the cursor must be on that line and then use `cmd+e`.
You can also evaluate source you have selected using the same shortcut `cmd+e`.
To evaluate a program, you need to select all the source and use `cmd+e`.

### play

To play a single line statement or expression, use `cmd+3`.
You can also evaluate the function `play(...)`.

### begin

To begin a loop, use `cmd+4`.
You can also evaluate the function `begin(...)`.

### end

To end a running loop, use `cmd+5`.
You can also evaluate the function `end(...)`.

### inspect

To inspect a variable or a function, just hover with your mouse pointer above its name. To explicitly inspect the value of an expression, use `cmd+2`.

### stop all sound

To stop sounds being played, including loops, use `cmd+k`.

### comment

Lines that start with `//` are not evaluated ; these are commment lines.

	// this is comment

Lines can have inline comment at the end.

	s = note('C#') // C sharp, Octave 4

### multiline

A statement can span multiple lines, each line after the first must be indented by either a TAB or 4 spaces.

	  y = sequence('F#2 
	  [TAB]C#3 F#3 A3 C# F#')
	  x = sequence('A 
	  [SPACE][SPACE][SPACE][SPACE]A5 A6')