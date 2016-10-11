Melrose - interactive programming of music melodies

### Note notation

| Notation | Alternative | Description 
|----------|-------|-------------
| C4       | C     | quarter C octave 4 
| 2E5      | ½C5   | Halftone (2 x ¼) C octave 5
| F#       | F♯    | F sharp
| G_       | G♭    | G flat
| G.       | G.    | duration x 1.5 
| =        | =     | quarter rest

### Sequence notation

| Notation    | Description 
|-------------|---
| C D E F     | 4 quarter tones
| (C E) (D F) | 2 doublets


### Sequence functions

| Expression       | Result
|------------------|----
| chord("C")       | (C E G)
| reverse("C D E") | E D C
| pitch("C D",1)   | D♭ E♭
| pitch("C D",-2)  | B♭3 C
| scale("C")       | C D E F G A B C5


Software is licensed under [Apache 2.0 license](LICENSE).
(c) 2014-2015 http://ernestmicklei.com 


Similar projects
- http://daveyarwood.github.io/alda/2015/09/05/alda-a-manifesto-and-gentle-introduction/
- http://www.lilypond.org/