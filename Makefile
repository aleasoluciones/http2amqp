all: build test clean

update_dep:
	go get $(DEP)
	go mod tidy

update_all_deps:
	go get -u all
	go mod tidy

format:
	go fmt ./...

test:
	go vet ./...
	- pkill -f echoservice
	./echoservice -topic '*.test.ok' &
	go clean -testcache
	go test -v -tags integration -parallel 2 ./... -timeout 60s
	- pkill -f echoservice

build:
	go build -o http2amqp examples/httpserver/http2amqp.go
	go build -o echoservice examples/echoservice/echoservice.go

build_images:
	docker build . --no-cache --target http2amqp-builder -t aleasoluciones/http2amqp-builder:${GIT_REV}
	docker build . --target http2amqp -t aleasoluciones/http2amqp:${GIT_REV}

clean:
	rm -f http2amqp
	rm -f echoservice

start_dependencies:
	docker-compose -f dev/http2amqp_devdocker/docker-compose.yml up -d

stop_dependencies:
	docker-compose -f dev/http2amqp_devdocker/docker-compose.yml stop

rm_dependencies:
	docker-compose -f dev/http2amqp_devdocker/docker-compose.yml down -v


.PHONY: all update_dep update_all_deps format build test clean start_dependencies stop_dependencies rm_dependencies
