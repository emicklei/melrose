- cannot use variable for indices in 
    - dynamicmap
    - pitchmap
    - octavemap
    - notemap
    - sequencemap
- group cannot have multiple sequences
- ? no longer use special unicode character, because of copy from output resuability ? 
- rotated not used
- recording needs rewrite
- terminal ui drop down for echoing
- interval storex
- stretch storex
- not supported
    - legato
    - staccato
- fractionmap missing (half way)
- somehow show current BPM and BIAB in plugins
- the :m command must become function to be useable in a script
- don't like the dynamic using o+ and o-
- rename watch to print

Roger Linn on MPC60 swing
"Swing – applied to quantized 16th-note beats – is a big part of it. My implementation of swing has always been very simple: I merely delay the second 16th note within each 8th note. In other words, I delay all the even-numbered 16th notes within the beat (2, 4, 6, 8, etc.) In my products I describe the swing amount in terms of the ratio of time duration between the first and second 16th notes within each 8th note. For example, 50% is no swing, meaning that both 16th notes within each 8th note are given equal timing. And 66% means perfect triplet swing, meaning that the first 16th note of each pair gets 2/3 of the time, and the second 16th note gets 1/3, so the second 16th note falls on a perfect 8th note triplet. "


LSP design
----------
https://blog.logrocket.com/how-to-use-the-language-server-protocol-to-extending-a-client-764da0e7863c/
https://github.com/rcjsuen/dockerfile-language-server-nodejs