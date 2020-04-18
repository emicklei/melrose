run: test build
	melrose

test:
	go test -cover ./...

build:
	go install github.com/emicklei/melrose/cmd/melrose

snippets:
	cd cmd/vsc && go run snippets.go > ../../../melrose-for-vscode/snippets/snippets.json