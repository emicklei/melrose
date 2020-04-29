---
title: Melrose installation guide
---

[Home](index.html)
[Usage](cli.html)
[Language](dsl.html)
[DAW](daw.html)
[Install](install.html)

# Install

In order to work with `melrose` on your operating system, the following components need to be installed:

- portmidi library
- melrose executable program
- (optionally) Melrose plugin for Visual Studio Code

Depending on your operating system, different steps are required.

## Mac OSX

    brew install portmidi

See [Brew](https://brew.sh/) for instructions on how to install `brew` on your Mac.

See the [PortMidi](https://sourceforge.net/p/portmedia/wiki/portmidi/) for alternative installation instructions of `portmidi`.

Proceed with [install melrose](install.md#install-melrose)

## Linux

On Ubuntu / Debian

	apt-get install libportmidi-dev
	
Proceed with [install melrose](install.md#install-melrose)

## Windows

[Download](https://sourceforge.net/projects/portmedia/files/portmidi/217/) and install the pre-compiled portmidi DLL.

Proceed with [install melrose](install.md#install-melrose)

## Install Melrose<a name="install-melrose"></a> 

Currently, `melrose` can only be installed from source.
You need to install the [Go SDK](https://golang.org/dl/) for compiling the program on your machine.

	go install github.com/emicklei/melrose/cmd/melrose
	
After installing both `portmidi` and `melrose`, you can start the tool in a Terminal using:

	$ melrose
	
If this command cannot be found then you need to add `$GOPATH/bin` to your `PATH`.

## Melrose plugin for Visual Studio Code<a name="plugin"></a>

Currently, the `melrose` plugin is not yet published so you need to [download the plugin](https://public.philemonworks.com/melrose/melrose-for-vscode-1.0.0.vsix) archive and put this in the `plugins` folder on your installed `Visual Studio Code` editor.