---
title: Melrōse Tutorial 01 Note
---

## create a note

    n = note('c')

This is a statement with a variable `n` and expression `note('c')`.
The expression uses the function `note` and takes the string argument `'c'`.
The argument `c` represents the quarter note middle C, octave 4.

## sharp or flat

    n = note('c#')

Using `#` or `♯` makes the note sharp. Using `_` or `♭` makes the note flat.

## duration

    n = note('2c#')

Change the duration of the note by prefixing a number.
The number `2` or `½` means set the duration to 0.5.
No number, or `4` or `¼` means set the duration to 0.25.
Valid numbers are 1,2,4,8,16.

## dynamic

    n = note('2c#++')

By changing the dynamic of a note can make it sound softer,quieter or harder,louder.
The symbol `-` is used to silence the note.
The symbol `+` is used to emphasize the note.
You can use up to 4 such symbols.

See [Notation](notations.html) for a complete description of the syntax to create notes.





[Next: 02 Sequence](02-sequence.html)