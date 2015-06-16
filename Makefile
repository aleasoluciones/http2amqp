all: deps build test

deps:
	go get -t -v ./...

test:
	#go test -v ./...

build:
	go build .
	go build -o http2amqp httpserver/http2amqp.go
	go build -o echoservice examples/echoservice/echoservice.go

.PHONY: deps test
