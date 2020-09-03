---
title: Melrōse
---

[Home](https://emicklei.github.io/melrose)
[Videos](https://emicklei.github.io/melrose/videos)
[Tutorial](https://emicklei.github.io/melrose/tutorial)
[Examples](https://emicklei.github.io/melrose/examples)
[Language](dsl.html)
[Notations](notations.html)
[Tool](cli.html)
[DAW](daw.html)
[Install](install.html)


## Introduction

`melrōse` is a tool to create and play music by programming melodies.
It uses a custom language to compose notes and create loops and tracks to play.
This is an example of a simple scale.

```javascript
sequence('C D E F G A B C5')
```

Note sequences in your program can be changed while playing giving you direct audible feedback. 
For the best experience, use the `melrōse` tool together with the Visual Studio Code Plugin for Melrōse.

Basic musical objects in Melrōse are:

- [note](dsl.html#note)
- [sequence](dsl.html#sequence)
- [chord](dsl.html#chord)
- [scale](dsl.html#scale)
- [progression](dsl.html#progression)