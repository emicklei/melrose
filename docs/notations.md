---
title: Melrōse musical object notation
---

## Note

Format: `(duration)(pitch)(accidental)(dynamic)`

| Notation | Alternative | Description
|----------|-------|-------------
| C4       | ¼C,C,c  | quarter C octave 4
| 2E5      | ½E5,½e5 | Halftone (2 x ¼) E octave 5
| 1C       |        | Full tone C octave 4
| F#       | F♯,f♯  | F sharp
| G_       | G♭    | G flat
| .G       |       | duration x 1.5 = 3/8
| =        | =     | quarter rest
| 2=       | ½=    | half rest
| 1=       | 1=    | full rest
| D+       | d+    | quarter D octave 4 MezzoForte
| 16.E#--  | 16.e♯-- | sixteenth E sharp duration x 1.5 Piano

    n = note('c#5')

## Note dynamics<a name="note-not"></a>

| Notation    | Description
|-------------|---
| \-\-\-      |Pianissimo (pp)
| \-\-	      |Piano (p)
| \-	      |MezzoPiano (mp)
| 0           |Normal (optional)
| +	          |MezzoForte (mf)
| ++	      |Forte (f)
| +++         |Fortissimo (ff)

    n = note('E++')

## Sequence<a name="sequence-not"></a>

| Notation    | Description
|-------------|---
| C D E F       | 4 quarter tones
| (C E) (d5 f5) | 2 doublets
| (1C 1E 1G)    | C Chord

    doremi = sequence('c d e')

## Chord<a name="chord-not"></a>

| Notation    | Description
|-------------|---
| C#5/m/2     | C sharp triad, Octave 5, Minor, 2nd inversion
| A/7         | A Dominant seventh chord
| E/M7        | E Major seventh chord
| G/m7        | G minor seventh chord
| 1=          | No chord, a whole rest note

    b7 = chord('b/7')

## Scale<a name="scale-not"></a>

| Notation    | Description
|-------------|---
| C5          | C major scale, Octave 5
| E/m         | E natural minor scale, Octave 4
| G/M7        | G major 7 scale, Octave 4

    sf = scale(2,'f')

## Progression<a name="progression-not"></a>

| Notation    | Description
|-------------|---
| C/m D/m     | C minor followed by a D minor 
| (C3 C5)     | C major, Octave 3 together with a C major, Octave 5
| E =         | E major followed by a quarter rest note

    p = progression('C F G')