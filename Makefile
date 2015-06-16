all: deps build test

deps:
	go get -t -v ./...

test:
	- pkill -f echoservice
	./echoservice >> salida &
	BROKER_URI=amqp://guest:guest@localhost/ go test -v -tags integration ./...

build:
	go vet
	go build .
	go build -o http2amqp httpserver/http2amqp.go
	go build -o echoservice examples/echoservice/echoservice.go

.PHONY: deps test
