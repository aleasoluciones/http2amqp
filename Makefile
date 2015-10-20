all: deps build test

deps:
	go get -t -v ./...
	go get github.com/stretchr/testify/assert

test:
	- pkill -f echoservice
	./echoservice -topic '*.test.ok'  &
	BROKER_URI=amqp://guest:guest@localhost/ go test -v -tags integration -parallel 2 ./...
	- pkill -f echoservice

build:
	go vet
	go build .
	go build -o http2amqp httpserver/http2amqp.go
	go build -a -installsuffix cgo httpserver/http2amqp.go -o http2amqp
	go build -a -installsuffix cgo examples/echoservice/echoservice.go -o echoservice

.PHONY: deps test
