LATEST_TAG := $(shell git describe --abbrev=0)

run: test install snippets grammar dslmd
	melrose

test:
	go test -cover ./...

build:
	export LATEST_TAG=$(shell git describe --abbrev=0)
	cd cmd/melrose && go build -ldflags "-s -w -X main.version=$(LATEST_TAG)" -o ../../target/melrose
	
install:
	go install github.com/emicklei/melrose/cmd/melrose

snippets:
	cd cmd/vsc && go run *.go snippets > ../../../melrose-for-vscode/snippets/snippets.json

grammar:
	cd cmd/vsc && go run *.go grammar  \
		../../../melrose-for-vscode/syntaxes/melrose.tmGrammar.json.template \
		../../../melrose-for-vscode/syntaxes/melrose.tmGrammar.json

dslmd:
	cd cmd/vsc && go run *.go dslmd

vsc: snippets grammar
	cd ../melrose-for-vscode && vsce package

clean:
	rm -rf target
	mkdir target

package: clean build vsc
	cp /usr/local/opt/portmidi/lib/libportmidi.dylib target
	cp run.sh target
	cp ../melrose-for-vscode/*vsix target
	cd target && zip -mr macosx-melrose-$(LATEST_TAG).zip .