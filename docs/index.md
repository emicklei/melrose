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
This is an example of a simple major scale C.

```javascript
sequence('c d e f g a b c5')
```

Note sequences in your program can be changed while playing giving you direct audible feedback. 
For the best experience, use the `melrōse` tool together with the Visual Studio Code Plugin for Melrōse.

Basic musical objects in Melrōse are:

- [note](dsl.html#note)
- [sequence](dsl.html#sequence)
- [chord](dsl.html#chord)
- [scale](dsl.html#scale)
- [progression](dsl.html#progression)

Musical object can be composed using many operators such as:

- [pitch](dsl.html#pitch)
- [resequence](dsl.html#resequence)
- [repeat](dsl.html#repeat)
- [dynamic](dsl.html#dynamic)
- [loop](dsl.html#loop)

Audio control objects such as:

- [bpm](dsl.html#bpm)
- [midi note](dsl.html#midi)
- [midi_send](dsl.html#midi_send)

See [Language](dsl.html) for all supported functions.

![screenshot.png](images/screenshot.png)


Software is licensed under [MIT](https://github.com/emicklei/melrose/LICENSE).
&copy; 2014-2021 [ernestmicklei.com](http://ernestmicklei.com)