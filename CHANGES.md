## v1.0.0:

## [v1.0.0-beta] - 

- add `melrose-mcp` server
- add tempo as alias for bpm
- add flag for log location
- allow 'b' for flat
- fix help for :e
- reworked cli midi settings 
- add euclidean rythm
- add map function
- allow replace the replacer
- fix Index(), iternator storex
- add 1/32 notes

## [v0.53.0] - Fri Apr 11 13:57:50 2025 +0200

- the scale function has changed (the breaking kind)
  instead of: 

    scale(2,'16a2')

  you write:

    join(scale('16a2'),scale('16a3'))   

  and now can create scales with a type:

    scale('major C') // Ionian mode

## [v0.52.0] - Thu Jul 25 09:40:54 2024 +0200

    - see git log for all past changes