all: clean build test

travis: clean build test

jenkins: clean build

update_dep:
	go get $(DEP)
	go mod tidy

update_all_deps:
	go get -u
	go mod tidy

test:
	go vet ./...
	- pkill -f echoservice
	./echoservice -topic '*.test.ok'  &
	go test -v -tags integration -parallel 2 ./...
	- pkill -f echoservice

build:
	go build -a -installsuffix cgo -o http2amqp examples/httpserver/http2amqp.go
	go build -a -installsuffix cgo -o echoservice examples/echoservice/echoservice.go

clean:
	rm -f http2amqp
	rm -f echoservice


.PHONY: all travis jenkins update_dep update_all_deps test build clean
