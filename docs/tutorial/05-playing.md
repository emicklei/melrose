---
title: Melrōse Tutorial 05 Playing
---

[Home](https://emicklei.github.io/melrose)

## play a musical object

There are several ways to start listening to your (composed) musical object.
Note that the editor plugin has keyboard shortcuts to play musical objects. 
### play

```javascript
s = sequence('C E G')
play(s)
```

Using the function `play` will schedule all the notes of a musical object.
The first note(s) starts at the moment of evaluation.
In this example the `C` note from the sequence `s` will be played immediately, the `E` and `G` are scheduled to play one and two quarter durations later.
The actual time is dependent on the value of BPM, beats-per-minute, at the time of evaluation.

### sync

```javascript
s1 = sequence('C E G')
s2 = sequence('C5 E5 G5')
play(s1,s2)
sync(s1,s2)
```

Using the function `sync` will schedule all musical objects at the same time.
The first notes of each musical object will start at the moment of evaluation.
In this example the `C` and `C5` notes from the sequences `s1` and `s2` are played immediately.
Using the `play` function, all musical object will be played after each other.

### begin, end

```javascript
s1 = sequence('C E G')
lp_s1 = loop(s1)
begin(lp_s1)
end(lp_s1)
```

The function `begin` and `end` apply to `Loop` and `Listen` objects only.

### track

### multitrack

### listen

### stop