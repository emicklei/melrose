---
title: Melrōse installation guide
---

[Home](https://emicklei.github.io/melrose)

# Download

Due to increased restrictions and costs of releasing non-commercial packages for the Apple Mac OSX and Microsoft Windows plaform, I can no longer provide pre-build binaries. Contact me if you want a commercial, supported version.

# Install from source

In order to work with `melrōse` on your operating system, the following components need to be installed:

- portmidi library
- melrōse executable program
- (optionally) Melrōse plugin for Visual Studio Code

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







## Install Melrōse<a name="install-melrōse"></a> 

Currently, `melrōse` can only be installed from source.
You need to install the [Go SDK](https://golang.org/dl/) for compiling the program on your machine.

	go install github.com/emicklei/melrose/cmd/melrose
	
After installing both `portmidi` and `melrōse`, you can start the tool in a Terminal using:

	$ melrose
	
If this command cannot be found then you need to add `$GOPATH/bin` to your `PATH`.

## Melrōse plugin for Visual Studio Code<a name="plugin"></a>

Currently, the `melrōse` plugin is not yet published on the Marketplace.
You install the plugin from the Extensions overview in Visual Studio Code and open the `.vsix` file from your extracted download.