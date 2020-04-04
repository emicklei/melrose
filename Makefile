run: test
	go install github.com/emicklei/melrose/cmd/melrose && melrose

test:
	go test -v -cover ./...