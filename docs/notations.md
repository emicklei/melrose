---
title: Melrōse musical object notation
---

[Home](https://emicklei.github.io/melrose)

## Note

Format: `(fraction)(dot)(name|=)(accidental)(dynamic)`

| Notation | Alternative | Description
|----------|-------|-------------
| C4       | ¼C,C,c     | quarter C, octave 4
| 2E5      | ½E5,½e5    | Halftone (2 x ¼), E octave 5
| 1C       | 1c         | Full tone C, octave 4
| F#       | F♯,f♯,f#   | F sharp
| G_       | G♭,g♭,g_   | G flat
| .G       | .g         | duration fraction x 1.5 = 3/8
| =        | =          | quarter rest
| 2=       | ½=         | half rest
| 1=       | 1=         | full rest
| D+       | d+         | quarter D, octave 4, MezzoForte
| 16.E#--  | 16.e♯--    | sixteenth, E sharp, fraction x 1.5, Piano

```javascript
n = note('c#5')
```

## Note dynamics<a name="note-not"></a>

| Notation    | Description
|-------------|---
| \-\-\-\-    |Pianissississimo (pppp)
| \-\-\-      |Pianississimo (ppp)
| \-\-        |Pianissimo (pp)
| \-          |Piano (p)
| o-          |MezzoPiano (mp)
| o (not 0)   |Normal (character is optional)
| o+          |MezzoForte (mf)
| +           |Forte (f)
| ++          |Fortissimo (ff)
| +++         |Fortississimo (fff)
| ++++        |Fortissississimo (ffff)

```javascript
n = note('e++')
```

## Sequence<a name="sequence-not"></a>

| Notation    | Description
|-------------|---
| C D E F       | 4 quarter tones
| (8C E) (d5 f5) | 2 doublets; first doublet has an eight length, second is a quarter
| (1C E G)    | C Chord; whole length

```javascript
doremi = sequence('c d e')
```

## Pedal control

Usable in `sequence` or `note`.

| Notation | Description
|----------|-------------
| >        | sustain pedal `down`
| <        | sustain pedal `up`
| ^        | sustain pedal `up` and immediately `down`


## Chord<a name="chord-not"></a>

| Notation    | Description
|-------------|---
| C#5/m/2     | C sharp triad, Octave 5, Minor, 2nd inversion
| A/7         | A Dominant seventh chord
| E/M7        | E Major seventh chord
| G/m7        | G minor seventh chord
| 1=          | No chord, a whole rest note
| D/dim       | D diminished triad
| D/o         | D diminished triad
| F/dim7/1    | F diminished seventh, 1st inversion
| C/aug       | C augmented triad
| E/+         | E augmented triad
| B_/+7       | B flat augmented seventh

```javascript
b7 = chord('b/7')
```

## Scale<a name="scale-not"></a>

| Notation    | Description
|-------------|---
| C5          | C major scale, Octave 5
| E/m         | E natural minor scale, Octave 4
| G/M7        | G major 7 scale, Octave 4

```javascript
sf = scale(2,'f')
```

## Progression<a name="progression-not"></a>

| Notation    | Description
|-------------|---
| C/m D/m     | C minor followed by a D minor
| (C3 C5)     | C major, Octave 3 together with a C major, Octave 5
| E =         | E major followed by a quarter rest note

```javascript
p = progression('c f g')
```