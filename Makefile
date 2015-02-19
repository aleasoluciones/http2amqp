all: deps test

deps:
	go get -t -v ./...

test:
	go test -v ./...

.PHONY: deps test
