all: deps build test

deps:
	go get -t -v ./...
	go get github.com/stretchr/testify/assert

test:
	- pkill -f echoservice
	./echoservice -topic '*.test.ok'  &
	go test -v -tags integration -parallel 2 ./...
	- pkill -f echoservice

build:
	go vet
	go build .
	go build -a -installsuffix cgo -o http2amqp httpserver/http2amqp.go
	go build -a -installsuffix cgo  -o echoservice examples/echoservice/echoservice.go

.PHONY: deps test
