# Install 

In order to build and work with `melrose` on your operating system, the following components need to be installed:

- melrose executable program
- (optionally) Melrōse plugin for Visual Studio Code

Depending on your operating system, different steps are required.


## [all platforms] Melrōse plugin for Visual Studio Code<a name="plugin"></a>

Search and install the extension from the editor or go to the [Marketplace published package](https://marketplace.visualstudio.com/items?itemName=EMicklei.melrose-for-vscode)


## [all platforms] Go SDK.

You need to install the [Go SDK](https://golang.org/dl/) for compiling the program on your machine.

## Mac OSX

	go install -v github.com/emicklei/melrose/cmd/melrose...@latest

After installing `melrōse`, you can start the tool in a Terminal using:

	$ melrose

If this command cannot be found then you need to add `$GOPATH/bin` to your `PATH`.

## Linux

On Ubuntu / Debian

	sudo apt-get install libasound2-dev

You need to install the [Go SDK](https://golang.org/dl/) for compiling the program on your machine.

	go install -v github.com/emicklei/melrose/cmd/melrose...@latest

After installing both `libasound2` and `melrōse`, you can start the tool in a Terminal using:

	$ melrose

If this command cannot be found then you need to add `$GOPATH/bin` to your `PATH`.

## Windows

Look at the build script (`.travis.yml`) of [melrose-windows](https://github.com/emicklei/melrose-windows) for detailed steps to build an executable from source.
