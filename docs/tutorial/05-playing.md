---
title: Melrōse Tutorial 05 Playing
---

[Home](https://emicklei.github.io/melrose)

## play a musical object

There are several ways to start listening to your (composed) musical object.
Note that the editor plugin has keyboard shortcuts to play musical objects instead of using functions such as `play`, `begin` and `end`.
### play

```javascript
s = sequence('C E G')
play(s)
```

The function `play` will schedule all the notes of a musical object.
The first note(s) start at the moment of evaluation.
In this example the `C` note from the sequence `s` will be played immediately, the `E` and `G` are scheduled to play one and two quarter durations later.
The actual time is dependent on the value of BPM, beats-per-minute, at the time of evaluation. You can change the tempo by using the `bpm` function.

### sync

```javascript
s1 = sequence('C E G')
s2 = sequence('C5 E5 G5')
sync(s1,s2)
play(s1,s2)
```

The function `sync` will schedule all musical objects at the same time.
The first notes of each musical object will start at the moment of evaluation.
In this example the `C` and `C5` notes from the sequences `s1` and `s2` are played immediately.
Using the `play` function, all musical object will be played after each other.

### begin, end loop

```javascript
s1 = sequence('C E G')
lp_s1 = loop(s1)
begin(lp_s1)
end(lp_s1)
```

The function `begin` and `end` apply to `Loop` and `Listen` objects only.
Basicly, it will play one or more musical objects, each after the other, repeatedly.
Both functions require a variable, here `lp_s1` to which a loop is assigned and is needed for you to stop playing.
The function `begin` will start playing the loop immediately.
The function `end` will stop this loop.

### track

```javascript
bpm(120) 

f1 = sequence('C D E C')
f2 = sequence('E F 2G')
f3 = sequence('8G 8A 8G 8F E C')
f4 = sequence('2C 2G3 2C 2=') 

v1 = join(f1,f1,f2,f2,f3,f3,f4) 

t = track('frere',1,
    onbar(1,v1),
    onbar(3,v1),
    onbar(5,v1))
play(t)
```

To schedule multiple sequences in a timeline, you can create a `track`.
A track uses a bar count to set the start time of playing a musical ohject.
The example is a well-know melody, "Frere Jacques", in which the same melody is repeated with an offset of 2 bars.
The variable `v1` is the combination of several sequences such as `f1` and `f2`.
The track specifies the MIDI channel, here `1`, and multiple entries on the timeline using the `onbar` function.

### multitrack

```javascript
all = multitrack(drumTrack,pianoTrack,bassTrack)
play(all)
export("my-project",all)
```

A multitrack is simply a collection of tracks.
It can be used to play all tracks at once or to export to a MIDI file.