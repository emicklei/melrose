LATEST_TAG := $(shell git describe --abbrev=0)

all: installl snippets grammar dslmd

udp:
	cd cmd/melrose/ && go build -tags=udp -o melrose-udp 

test:
	go vet ./...
	go test -race -cover ./...

unused:
	env GO111MODULE=on go get honnef.co/go/tools/cmd/staticcheck@v0.0.1-2020.1.4
	staticcheck --unused.whole-program=true -- ./...

build:
	export LATEST_TAG=$(shell git describe --abbrev=0)
	cd cmd/melrose && go build -ldflags "-s -w -X core.BuildTag=$(LATEST_TAG)" -o ../../target/melrose
	
install: test
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

APP := /Applications/Melrose
package: clean build  
	# prepare target
	cp /usr/local/opt/portmidi/lib/libportmidi.dylib target 
	echo "$(LATEST_TAG)" > target/version.txt
	# copy to APP
	rm -rf $(APP)
	mkdir -p $(APP)	
	cp target/melrose $(APP)
	cp packaging/macosx/*.sh $(APP)
	cp target/version.txt $(APP)
	# package it up
	/usr/local/bin/packagesbuild --package-version "$(LATEST_TAG)" packaging/macosx/Melrose.pkgproj
	mv packaging/macosx/Melrose.pkg "packaging/macosx/Melrose-$(LATEST_TAG).pkg"

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