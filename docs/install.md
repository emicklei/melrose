---
title: Melrose installation guide
---

[Home](index.html)
[Usage](cli.html)
[Language](dsl.html)
[DAW](daw.html)
[Install](install.html)

# Install

In order to work with `melrōse` on your operating system, the following components need to be installed:

- portmidi library
- melrōse executable program
- (optionally) Melrose plugin for Visual Studio Code

Depending on your operating system, different steps are required.

## Mac OSX

    brew install portmidi

See [Brew](https://brew.sh/) for instructions on how to install `brew` on your Mac.

See the [PortMidi](https://sourceforge.net/p/portmedia/wiki/portmidi/) for alternative installation instructions of `portmidi`.

Proceed with [install melrōse](install.md#install-melrōse)

## Linux

On Ubuntu / Debian

	apt-get install libportmidi-dev
	
Proceed with [install melrōse](install.md#install-melrōse)

## Install Melrose<a name="install-melrōse"></a> 

Currently, `melrōse` can only be installed from source.
You need to install the [Go SDK](https://golang.org/dl/) for compiling the program on your machine.

	go install github.com/emicklei/melrōse/cmd/melrōse
	
After installing both `portmidi` and `melrōse`, you can start the tool in a Terminal using:

	$ melrōse
	
If this command cannot be found then you need to add `$GOPATH/bin` to your `PATH`.

## Melrose plugin for Visual Studio Code<a name="plugin"></a>

Currently, the `melrōse` plugin is not yet published so you need to [download the plugin](https://public.philemonworks.com/melrōse/melrōse-for-vscode-1.0.0.vsix) archive. You install the plugin from the Extensions overview in Visual Studio Code.