---
title: Melrōse Design
---

[Home](index.html)
[Usage](cli.html)
[Language](dsl.html)
[DAW](daw.html)
[Install](install.html)

## Desgin

Melrōse is a tool to describe, play and interactively change music by manipulating objects that represent sequences of notes.

The tool by itself does not produce any sound playing these notes.
Instead, it produces MIDI events that can be send to a Digital Work Station (DAW).

Melrōse uses a custom language to describe music in a way that is easy to read and write.

Basic concepts in the language are:

- variables
- functions
    - creation of basic musical objects
    - creation of operations on basic musical objects ; operations themselves are musical objects
    - play, record and export of musical objects
- expressions composed of functions and variables

Melrōse provides audio feedback, syntax error feedback, object inspection feedback, measurements.

A program is an ordered list of statements.
A statement can be an expression or an assigment.
An assignment has a variable and an expressions.
An expression is variable or a function call.
A function call uses constants or variables or an expression.

A progam can have any number of empty lines and comments.
A comment starts with `//` and ends with a newline.

The most basic musical object is a Note.
There are two functions to create a single Note: `note` and `midi`.

    note('c')

This creates an immutable object Note that represent the middle C, has a quarter length, octave 4 and no dynamic.

    note('8e#5++')

