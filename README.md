# http2amqp
[![Build Status](https://app.travis-ci.com/aleasoluciones/http2amqp.svg?branch=master)](https://app.travis-ci.com/github/aleasoluciones/http2amqp)


HTTP interface to AMPQ. 

## Features
* It publishes an amqp message for each http request received and process the corresponding amqp responses (it waits for it) in order to answer to the original http request.
* The topic of the message it publishes comes from the URL Path of the HTTP request, using the HTTP method, network and the path replacing '/' with '.'
  * E.g.: 'get.arl.cpe'
* The TTL for the messages published is 1000 ms
* The exchange used by default is 'events'


## Build
A Makefile is available so you only need to run
```
make
```

## Running tests
Before running test:

* Ensure you have `echoservice` binary is built. To build `echoservice` run `make build` 
```
make build
```
* Start a rabbitmq service with default credentials

```
docker run --rm -d --name http2amqp-rabbit -p5672:5672 rabbitmq:3
```

You can use BROKER_URI env for setting your custom rabbitmq values

```
BROKER_URI="amqp://guest:guest@localhost/" make test
```

* Run tests with `make test`. `Makefile` has a test section for running tests.
```
make test
```

## Building docker image
Also there is a script for building a docker image

```
./build.sh
```

## Usage
```
$ ./http2amqp --help
Usage of ./http2amqp:
  -address="0.0.0.0": Listen address
  -amqpuri="amqp://guest:guest@localhost/": AMQP connection uri
  -exchange="events": AMQP exchange name
  -port="18080": Listen port
  -timeout=1000: Queries timeout in milliseconds
  -verbose :Enable logging, false by default
```

## Debugging

Running the container with verbose mode to debugg whats happening

```
 docker-compose -f docker-compose.yml run http2amqp -verbose
```

## Execution example
With a rabbitmq running with the default credentials...

Start the htt2amqp server in a terminal
```
./http2amqp
```

Start in another terminal the echo service
```
./echoservice
```

Tests diferent get requests to be served by the echo service
```
curl -X GET http://localhost:18080/test -d 'hello world'
```

You can specify the timeout (in milliseconds) for waiting for the response
```
curl -X GET http://localhost:18080/test?timeout=200 -d 'hello world'
```

## TODO
 - test timeout parameter for each request
 - implement delay parameter for echo server to allow tests timeouts

