# Install from source

In order to build and work with `melrose` on your operating system, the following components need to be installed:

- RtMidi library
- melrose executable program
- (optionally) Melrōse plugin for Visual Studio Code

Depending on your operating system, different steps are required.

## Mac OSX

You need to install the [Go SDK](https://golang.org/dl/) for compiling the program on your machine.

	go install github.com/emicklei/melrose/cmd/melrose

After installing both `portmidi` and `melrōse`, you can start the tool in a Terminal using:

	$ melrose

If this command cannot be found then you need to add `$GOPATH/bin` to your `PATH`.

## Linux

On Ubuntu / Debian

	sudo apt-get install libasound2-dev

You need to install the [Go SDK](https://golang.org/dl/) for compiling the program on your machine.

	go install github.com/emicklei/melrose/cmd/melrose

After installing both `libasound2` and `melrōse`, you can start the tool in a Terminal using:

	$ melrose

If this command cannot be found then you need to add `$GOPATH/bin` to your `PATH`.

## Windows

Follow the instructions for installing `portmidi` as documented by [Zack](https://schollz.com/blog/portmidi/).
Then compile the program on your machine.

	go install github.com/emicklei/melrose/cmd/melrose

After installing both `portmidi` and `melrōse`, you can start the tool in a Terminal using:

	$ melrose
