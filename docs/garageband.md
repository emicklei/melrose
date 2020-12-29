---
title: Using Melrōse with Garageband
---

# Using Melrōse with Garageband

[Home](https://emicklei.github.io/melrose)

![garageband](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/garageband.png)

###### GarageBand is a registered trademark of Apple Inc. 

## Goal

Garageband is a simple Digital Audio Workstation (DAW) that comes with the standard installation of Apple Mac OSX on devices such as MacBook or MacPro.
Melrōse can communicate with GarageaBand by exchanging MIDI messages. GarageBand has a rich set of sounds and instruments to play MIDI notes. 

This article describes the steps to get Melrōse working with GarageBand such that you can play your melodies.

### Install melrōse

Install the latest packaged (unsigned) release of Melrōse for Mac OSX 10+.

#### versions

- [Melrose-v0.36.0.pkg](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/versions/Melrose-v0.36.0.pkg)
- [Melrose-v0.35.0.pkg](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/versions/Melrose-v0.35.0.pkg)
- [Melrose-v0.34.0.pkg](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/versions/Melrose-v0.34.0.pkg)
- [Melrose-v0.33.0.pkg](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/versions/Melrose-v0.33.0.pkg)
- [Melrose-v0.32.0.pkg](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/versions/Melrose-v0.32.0.pkg)

### Configure Audio MIDI Setup

In order to make Melrōse communicate with GarageBand, you need to enable an IAC Driver which is part of standard Mac OSX.
Open the separate Audio MIDI Setup program, next to GarageBand.
Use Spotlight Search to find this program.
Once started, then from the menu choose `Window -> Show MIDI Studio`.
Make sure you have an IAC Driver listed.

![iac](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/iacdriver.png)

Open the settings of this IAC Driver and make sure the device is online.

![iac online](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/iac_online.png)

### Install extension for Visual Studio Code

[Visual Studio Code](https://code.visualstudio.com/download) is a popular free open-source file editor, sponsored by Microsoft, also available for Mac OSX.
The Melrōse extension adds keyboard combinations to play and validate musical objects.
You can install the extension directly from your running Visual Studio Code Editor or by going to the [Marketplace published package](https://marketplace.visualstudio.com/items?itemName=EMicklei.melrose-for-vscode)

### Create a file

Open a new file `demo.mel` in the Visual Studio Code Editor.
The name suffix `.mel` tells the melrōse extension to activate the keyboard combinations. 
You can verify this by looking at the bottom right of the editor window where the file type is recognised as `Melrose`.

### Create GarageBand project

Open the GarageBand application and create a new Empty Project.
Then choose Software Instrument.

![instrument](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/gb_software_instrument.png)

This opens a window with a Classic Electic Piano instrument.


### Using melrose


#### Security warnings

If you have installed melrōse from the public package then you are installing software that is not verified by any party other than the developer who published the release. With newer versions of operating systems, both Apple and Microsoft are more restrictive when it comes to installing software. Currently, it is still allowed to install unregistered software (most open-source packages are) but the you, the user, will be asked to accept the risk.

You have to pass 3 steps of security checks:

1. When downloading the release (.zip archive) from Github, your computer will detect that it contains an application. It will ask you to proceed.

2. When starting the application melrōse, your computer will detect that the developer of the application is not verified. It will ask you to accept.

![lib not verified](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/osx_warning_app.png)

3. When the application melrōse is loading an extra library for MIDI access, your computer will detect that the developer of the application is not verified. It will ask you to accept.

![app not verified](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/osx_warning_portmidi.png)

### Start melrose

The previously installed extension needs to communicate with Melrōse to play your melodies.
Starting melrōse should be done from a Terminal pane of the Visual Studio Code editor.
This way, you can view any messages reported by melrōse and it also give you access to all its commands.
You can open this pane using the menu `Terminal -> New Terminal`.
Within the terminal pane, start the melrōse program, e.g:

```bash
melrose
```

You can verify the connectivity of melrōse and the IAC Driver by executing the command `:m` in the terminal.
It should have entries such as:

```bash
[midi] device id 0: "CoreMIDI/IAC Driver Bus 1", open=false, usage=input
[midi] device id 1: "CoreMIDI/IAC Driver Bus 1", open=true, usage=output
```

Use `:h` to see all available commands and functions.

### Play your first melody

Paste the following program in your demo file

```javascript
bpm(120)
y = sequence('e+ a- c5- b- c5- a- e+ f+ a- c5- b- c5- a- f-')
p = interval(-2,2,1)
ly = loop(pitch(p, fraction(8,y)), next(p)) // quarter->eight, semitones interval
```

![demo in vsc](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/vsc_melrose_demo.png)

Place the cursor anywhere on line 2 containing the sequence and press `cmd+3`.
This will both evaluate the expression and play the result.
You should hear notes being played using the instrument selected (most likely a piano) in GarageBand.

Now select the complete script `cmd+A` and press `cmd+3` and you will hear a loop with changing pitch.

### What's next

Visit the [Melrōse documentation](https://emicklei.github.io/melrose/) to find information about the programming language and the program itself. It also offers tutorials, examples and recorded demo videos.

Happy music coding!
