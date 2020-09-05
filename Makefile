LATEST_TAG := $(shell git describe --abbrev=0)

refresh: test install snippets grammar dslmd

run: install
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

clean:
	rm -rf target
	mkdir target

vsc:
	cd ../melrose-for-vscode && vsce package

package: clean build vsc
	cp /usr/local/opt/portmidi/lib/libportmidi.dylib target
	cp run.sh target
	cp ../melrose-for-vscode/*vsix target
	cd target && zip -mr macosx-melrose-$(LATEST_TAG).zip . && md5 macosx-melrose-$(LATEST_TAG).zip > macosx-melrose-$(LATEST_TAG).zip.md5


# go get -u -v github.com/aktau/github-release
# export GITHUB_TOKEN=$(kiya me get github/emicklei/macbookhub)
.PHONY: createrelease
createrelease:
	github-release info -u emicklei -r melrose
	github-release release \
		--user emicklei \
		--repo melrose \
		--tag $(LATEST_TAG) \
		--name "melrose" \
		--description "melrōse - program your melodies"

.PHONY: uploadrelease
uploadrelease:
	github-release upload \
		--user emicklei \
		--repo melrose \
		--tag $(LATEST_TAG) \
		--name "macosx-melrose-$(LATEST_TAG).zip" \
		--file target/macosx-melrose-$(LATEST_TAG).zip	