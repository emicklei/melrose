run: test build snippets grammar dslmd
	melrose

fast: build
	melrose

test:
	go test -cover ./...

build:
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

package: build vsc
	rm -rf target
	mkdir target
	cp /usr/local/opt/portmidi/lib/libportmidi.dylib target
	cp ${GOPATH}/bin/melrose target
	cp run.sh target
	cp ../melrose-for-vscode/*vsix target
	