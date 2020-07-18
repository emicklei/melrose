LATEST_TAG := $(shell git describe --abbrev=0)

refresh: test install snippets grammar dslmd menu

run: refresh
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

menu:
	cd cmd/vsc && go run *.go menu

vsc: snippets grammar
	cd ../melrose-for-vscode && vsce package

clean:
	rm -rf target
	mkdir target

package: clean build vsc
	cp /usr/local/opt/portmidi/lib/libportmidi.dylib target
	cp run.sh target
	cp ../melrose-for-vscode/*vsix target
	cd target && zip -mr macosx-melrose-$(LATEST_TAG).zip . && md5 macosx-melrose-$(LATEST_TAG).zip > macosx-melrose-$(LATEST_TAG).zip.md5