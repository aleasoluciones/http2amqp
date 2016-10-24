all: deps build test

deps:
	go get -d -v ./...
	go get -d -v github.com/stretchr/testify/assert
	go get -v github.com/golang/lint/golint

update_deps:
	go get -d -v -u ./...
	go get -d -v -u github.com/stretchr/testify/assert
	go get -v -u github.com/golang/lint/golint


test:
	golint ./...
	go vet ./...
	- pkill -f echoservice
	./echoservice -topic '*.test.ok'  &
	go test -v -tags integration -parallel 2 ./...
	- pkill -f echoservice

build:
	go build .
	go build -a -installsuffix cgo -o http2amqp httpserver/http2amqp.go
	go build -a -installsuffix cgo  -o echoservice examples/echoservice/echoservice.go

.PHONY: deps update_deps test build
