---
title: Melrōse Tutorial 01 Notes
---

## create a note

    n = note('c')

This is a statement with a variable `n` and expression `note('c')`.
The expression uses the function `note` and takes the string argument `'c'`.
The argument `c` represents the quarter note middle C, octave 4.

## sharp or flat

    n = note('c#')

Using `#` or `♯` make the note shape. Using `_` or `♭` makes the note flat.

## duration

    n = note('2c#')

Change the duration of the note by adding a number.
The number `2` or `½` means set the duration to 0.5.
No number, or `4` or `¼` means set the duration to 0.25.
Valid numbers are 1,2,4,8,16.
