.PHONY: build test clean

build:
	go build -o bin/codegraft ./...

test:
	go test ./...

clean:
	rm -rf bin/
