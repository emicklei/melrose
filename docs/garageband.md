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

[Download](https://github.com/emicklei/melrose/releases) the latest release of melrōse and unzip it into your preferred folder. Your Finder should look like this:

![finder](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/melrose_finder.png)

### Configure Audio MIDI Setup

In order to make melrōse communicate with GarageBand, you need to enable an IAC Driver which is part of standard Mac OSX.
Open the separate Audio MIDI Setup program, next to GarageBand.
Use Spotlight Search to find this program.
Once started, then from the menu choose `Window -> Show MIDI Studio`.
Make sure you have an IAC Driver listed.

![iac](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/iacdriver.png)

Open the settings of this IAC Driver and make sure the device is online.

![iac online](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/iac_online.png)

### Install extension for Visual Studio Code

[Visual Studio Code](https://code.visualstudio.com/download) is a popular free open-source file editor, sponsored by Microsoft, also available for Mac OSX. 
The Melrōse extension adds a few keyboard combinations to play and validate musical objects.

The folder to which you unzipped the downloaded archive, contains a file that ends with the name `.vsix`. 
This file contains the source code of the extension.
To install this extension, you need to open the Visual Studio Code Editor and go to `Code -> Preferences -> Extensions`.
This opens the install extensions for your editor.
Clicking on the dotted menu on the top left will open more actions. 
Choose `Install from VSIX...` and select the extension file from the melrose folder.

### Create a file

Open a new file `demo.mel` in the Visual Studio Code Editor.
The name suffix `.mel` tells the melrōse extension to activate the keyboard combinations. 
You can verify this by looking at the bottom right of the editor window where the file type is recognised as `Melrose`.

### Create GarageBand project

Open the GarageBand application and create a new Empty Project.
Then choose Software Instrument.

![instrument](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/gb_software_instrument.png)

This opens a window with a Classic Electic Piano instrument.


### Start melrose

The previously installed extension needs to communicate with melrōse to play your melodies.
Starting melrōse should be done from a Terminal pane of the Visual Studio Code editor.
This way, you can view any messages reported by melrōse and it also give you access to all its commands.
You can open this pane using the menu `Terminal -> New Terminal`.
Within the terminal pane, change to the directory that contains the downloaded melrōse program, e.g:

```bash
cd Melrose
```

Start the program using the following command:

```bash
./run.sh
```

You can verify the connectivity of melrōse and the IAC Driver by executing the command `:m` in the terminal.
It should have entries such as:

```bash
[midi] device id 0: "CoreMIDI/IAC Driver Bus 1", open=false, usage=input
[midi] device id 1: "CoreMIDI/IAC Driver Bus 1", open=true, usage=output
```

Use `:h` to see all available commands and functions.

### Play your first melody

Paste the following expression in your demo file

```javascript
sequence('c e d f e g f a g b a c5 b d5 c5')
```

![demo in vsc](https://storage.googleapis.com/downloads.ernestmicklei.com/melrose/vsc_melrose_demo.png)

Place the cursor anywhere on line 3 containing the expression and press `cmd+3`.
This will both evaluate the expression and play the result.
You should hear notes being played using the instrument selected (most likely a piano) in GarageBand.

### What's next

Visit the [Melrōse documentation](https://emicklei.github.io/melrose/) to find information about the programming language and the program itself. It also offers tutorials, examples and recorded demo videos.

Happy music coding!
