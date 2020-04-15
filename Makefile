run: test build
	melrose

test:
	go test -v -cover ./...

build:
	go install github.com/emicklei/melrose/cmd/melrose