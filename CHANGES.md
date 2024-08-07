### v0.53.0:

    - the scale function has changed (the breaking kind)
      instead of: 

        scale(2,'16a2')

      you write:

        join(scale('16a2'),scale('16a3'))   

      and now can create scales with a type:

        scale('major C') // Ionian mode

### v0.52.0:

    - see git log for all past changes