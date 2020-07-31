---
title: Melrōse Tutorial 02 Sequence
---

[Home](https://emicklei.github.io/melrose)

## create a sequence

    sequence('c d (e f)')

The expression uses the function `sequence` and takes the string argument `'c d (e f)'`.
This argument uses the notation for a note.
The notes inside brackets, here E and F, will be played as the same time.
Another, more complex example is:

    sequence('(8c 8e 8g) = (2c++ 2e++ 2g++)')

[Next: 03 Chord](03-chord.html)