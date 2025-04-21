all: install snippets grammar

test:
	go vet ./...
	go test -cover ./...

unused:
	# go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...
	
install: test q

# quickly get me a new binary
q:
	go install github.com/emicklei/melrose/cmd/melrose

snippets:
	cd cmd/vsc && go run *.go snippets > ../../../melrose-for-vscode/snippets/snippets.json

grammar:
	cd cmd/vsc && go run *.go grammar  \
		../../../melrose-for-vscode/syntaxes/melrose.tmGrammar.json.template \
		../../../melrose-for-vscode/syntaxes/melrose.tmGrammar.json

since:
	git log --oneline v0.55.1..@ > since.log