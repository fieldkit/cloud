all: build build/example

build:
	mkdir -p build

example/secrets.go:
	cp example/secrets.go.template example/secrets.go

build/example: example/*.go example/secrets.go
	go build -o build/example example/*.go

deps: example/secrets.go
	go get ./...

clean:
	rm -rf build

.PHONY: all clean
