Melrose - interactive programming of music melodies

### Note notation

| Notation | Alternative | Description 
|----------|-------|-------------
| 1C4      | C     | C octave 4 
| 2E5      | ½C5   | Halftone C octave 5
| F#       | F♯    | F sharp
| G_       | G♭    | G flat
| G.       | G.    | duration x 1.5 
| r        | r     | quarter rest

### Sequence notation

| Notation    | Description 
|-------------|---
| C D E F     | 4 half tones
| (C E) (D F) | 2 doublets


### Sequence functions

| Expression       | Result
|------------------|----
| chord("C")       | (C E G)
| reverse("C D E") | E D C
| pitch("C D",1)   | D♭ E♭
| pitch("C D",-2)  | B♭3 C
| scale("C")       | C D E F G A B C5

(c) 2014-2015 http://ernestmicklei.com Sofware is licensed under [Apache 2.0 license](LICENSE).