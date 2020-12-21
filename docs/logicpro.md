---
title: Using Melrōse with Logic Pro
---

# Using Melrōse with Logic Pro X

[Home](https://emicklei.github.io/melrose)

![logicpro](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/logicpro.png)

###### Logic Pro is a product of Apple Inc.

## Goal

Logic Pro is a professional Digital Audio Workstation (DAW) available for Apple Mac OSX.
Melrōse can communicate with Logic Pro by sending MIDI messages. 
Logic Pro has a rich set of sounds and instruments to play MIDI notes. 
In addition, Logic Pro is able to play multiple MIDI channels simultaneously which can be programmed individualy by `melrōse`.

This article describes the steps to get Melrōse working with Logic Pro such that you can play your melody per-channel.

### Installation

The installation steps are the same as described for [GarageBand](https://emicklei.github.io/melrose/garageband). 
Follow this article and open Logic instead of GarageBand. 

Keep the editor and the `demo.mel` file open.
Enter the following program into the editor, replaceing the previous example.

```javascript
bpm(100)

y1 = fraction(8,sequence('e+ a- c5- b- c5- a- e+ f+ a- c5- b- c5- a- f-'))
y2 = fraction(8,sequence('c+ a3- f3- e3- f3- a3- c+ c+ g3- e3- d3- e3- g3- c+'))

// sync(y1,y2)

p = interval(-2,2,1)
p1 = pitch(p,y1)
p2 = pitch(p,y2)

lp_1 = loop(channel(1,p1))         // Luminous Tines
lp_2 = loop(channel(2,p2),next(p)) // Plectrum Pad

begin(lp_1)
begin(lp_2)

// end(lp_1)
// end(lp_2)
```

### Create project

Open the application Logic Pro and create a new project.
Go to the project settings, Recording, and enable `auto demix by channel if multi-track recording`.


### Add Instruments

For the demo, we need two tracks with 2 different instruments.

![logic](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/logicpro_demo.png)

Make sure track 1 is receiving from MIDI Channel 1, track 2 from Channel 2.


### Play

#### test play one expression
Put the cursor on line 3 and hit `cmd+3` for instant play of that sequence.

### play all

Select all the text, `cmd+A` and hit `cmd+3`.
This will evaluate all statments and expressions cause 2 loops to begin at the same time.

### stop all

Select the text `end(lp_1)` and this `cmd+3` to stop the first loop. 
You must repeat it for the other loop.

### What's next

Visit the [Melrōse documentation](https://emicklei.github.io/melrose/) to find information about the programming language and the program itself. It also offers tutorials, examples and recorded demo videos.

Happy music coding!
