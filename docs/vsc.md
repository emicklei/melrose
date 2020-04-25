---
title: Plugin for Visual Studio Code
---

[Home](index.html)
[Usage](cli.html)
[Language](dsl.html)
[DAW](daw.html)
[Install](install.html)

## Plugin for Visual Studio Code

### install

See [instructions](install.html#plugin) on how to install this plugin after you have installed Visual Studio Code.

### edit

This editor extension works with `.mel` and `.melrose` files.
The syntax of the program uses the [Melrose Lanugage](dsl.html).


### evaluate

To evaluate a single line statement or expression, the cursor must be on that line and then use `cmd+e`.
You can also evaluate source you have selected using the same shortcut `cmd+e`.
To evaluate a program, you need to select all the source and use `cmd+e`.

### comment

Lines that start with `//` are not evaluated ; these commment lines.

	// this is comment

Lines can have inline comment at the end.

	s = note('C#') // C sharp, Octave 4

### multiline

An statement can span multiple lines, each line after the first must be indented by either a TAB or 4 spaces.

	  y = sequence('F#2 
	  [TAB]C#3 F#3 A3 C# F#')
	  x = sequence('A 
	  [SPACE][SPACE][SPACE][SPACE]A5 A6')