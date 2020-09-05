---
title: Melrōse Tutorial 02 Sequence
---

[Home](https://emicklei.github.io/melrose)

## create a sequence

```javascript
sequence('c d (e f)')
```

The expression uses the function `sequence` and takes the string argument `'c d (e f)'`.
This argument uses the notation for a note.
The notes inside brackets, here E and F, will be played as the same time.
Another, more complex example is:

```javascript
sequence('(8c e g) = (2c++ e++ g++)')
```

The first group `(8c e g)` has duration 1/8 the second 1/2. 
Only the first note should set the duration of the group.

### sustain pedal

In a sequence, you can control the postion of the sustain pedal.
To press the pedal down, you use the character `>`.
To release the pedal up, you use the character `<`.
The character `^` is short for up and down.

```javascript
sequence('> c d e ^ e d c <')
```

[Next: 03 Chord](03-chord.html)