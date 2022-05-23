.SILENT: build test

build:
	go vet ./... && go fmt ./... && go build -o a.out

test: build
	./a.out ./.test/test.go
