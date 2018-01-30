all: clean build test

jenkins: install_dep_tool install_go_linter production_restore_deps clean build test

install_dep_tool:
	go get github.com/tools/godep

install_go_linter:
	go get -u -v github.com/golang/lint/golint

initialize_deps:
	go get -d -v ./...
	go get -d -v github.com/stretchr/testify/assert
	go get -v github.com/golang/lint/golint
	godep save ./...

update_deps:
	godep go get -d -v ./...
	godep go get -d -v github.com/stretchr/testify/assert
	godep go get -v github.com/golang/lint/golint
	godep update ./...

test:
	golint ./...
	godep go vet ./...
	- pkill -f echoservice
	./echoservice -topic '*.test.ok'  &
	godep go test -v -tags integration -parallel 2 ./...
	- pkill -f echoservice

build:
	godep go build -a -installsuffix cgo -o http2amqp httpserver/http2amqp.go
	godep go build -a -installsuffix cgo -o echoservice examples/echoservice/echoservice.go

clean:
	rm -rf http2amqp
	rm -rf echoservice

production_restore_deps:
	godep restore

.PHONY: deps update_deps test build
