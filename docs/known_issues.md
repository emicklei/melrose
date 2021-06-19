# Known Issues

- cannot use variable for indices in 
    - dynamicmap
    - transposemap
    - octavemap
    - notemap
    - sequencemap
- group cannot have multiple sequences 
- rotated not used 
- not supported
    - legato
    - staccato 
- decreasing velocity , each note(group) gets a lower velocity, linear of fixed, or log. increasing too?
    crescendo(10,127,seq)
    decrescendo(127,10,seq)
        dynamicrange('++++','----',seq)
        fader(algo,sequence('c d e'))
- group only takes one sequenceable
- volume for offsetting the velocity