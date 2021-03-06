LATEST_TAG := $(shell git describe --abbrev=0)

all: installl snippets grammar

test:
	go vet ./...
	go test -race -cover ./...

unused:
	# env GO111MODULE=on go get honnef.co/go/tools/cmd/staticcheck@v0.0.1-2020.1.4
	staticcheck --unused.whole-program=true -- ./...

build: 
	cd cmd/melrose && go build -ldflags "-s -w -X main.BuildTag=$(LATEST_TAG)" -o ../../target/melrose
	
install: test
	go install github.com/emicklei/melrose/cmd/melrose

restart: 
	# ps | awk '/melrose$$/ {print $$1}' | xargs kill
	go install github.com/emicklei/melrose/cmd/melrose
	melrose

snippets:
	cd cmd/vsc && go run *.go snippets > ../../../melrose-for-vscode/snippets/snippets.json

grammar:
	cd cmd/vsc && go run *.go grammar  \
		../../../melrose-for-vscode/syntaxes/melrose.tmGrammar.json.template \
		../../../melrose-for-vscode/syntaxes/melrose.tmGrammar.json

clean:
	rm -rf target
	mkdir target

APP := /Applications/Melrose
package: clean build  
	# prepare target, build results is already in target
	# cp /usr/lib/libSystem.B.dylib target
	# cp /usr/lib/libc++.1.dylib target
	echo "$(LATEST_TAG)" > target/version.txt
	# copy to APP
	rm -rf $(APP)
	mkdir -p $(APP)	
	cp packaging/macosx/*.sh $(APP)
	cp target/* $(APP)
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