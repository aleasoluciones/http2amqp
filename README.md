# http2amqp

[![Build Status](https://app.travis-ci.com/aleasoluciones/http2amqp.svg?branch=master)](https://app.travis-ci.com/github/aleasoluciones/http2amqp)
[![GoDoc](https://godoc.org/github.com/aleasoluciones/http2amqp?status.png)](http://godoc.org/github.com/aleasoluciones/http2amqp)
[![License](https://img.shields.io/github/license/aleasoluciones/http2amqp)](https://github.com/aleasoluciones/http2amqp/blob/master/LICENSE)

HTTP interface to AMPQ.

## Features

- It publishes an AMQP request message for each HTTP request received. Then it waits for the corresponding AMQP response message in order to answer to the original HTTP request.
- The topic of the message that it publishes comes from the method and the path of the HTTP request, replacing slashes (/) with dots (.).
  - e.g.: `GET http://localhost:18080/net/test` → `get.net.test`
- The topic of the response messages to close the connection is `queries.response`. These response messages must carry the same ID as the request they intend to close.
- The default TTL for the published messages is 1000 ms. After that, the connection will be closed with a timeout if a response did not arrive.
- The exchange used by default is 'events'.


## Build

You need a Go runtime installed in your system which supports [modules](https://tip.golang.org/doc/go1.16#modules). A nice way to have multiple Go versions and switch easily between them is the [g](https://github.com/stefanmaric/g) application.

A [Makefile](Makefile) is available, so you only need to run:

```sh
make build
```

## Running tests

Make sure that you have built the binaries from the previous section, because the tests will run the `echoservice` binary. 

Load environment variables to set the BROKER_URI environment variable.

```sh
source dev/env_develop
```

Start a RabbitMQ service with default configuration (specified in [`/dev/env_develop`](/dev/env_develop)).

```sh
make start_dependencies
```

Run tests. They will only work if the RabbitMQ container is up.

```sh
make test
```

## Building docker images

The Jenkins [pipeline](Jenkinsfile) will generate two containers, which we can also build locally. One for compiling the app (http2amqp-builder), and the other to deploy its binary (http2amqp). This avoids the need of having Go installed in the host system.  Take a look at the [Dockerfile](Dockerfile) for more details.

```sh
make build_images
```

Once built, if we want to run the tests inside the builder image, we can do the following (remember to have the environment variables loaded and the RabbitMQ container up too):

```sh
docker run --rm -it --net=host aleasoluciones/http2amqp:GIT_REV integration-tests
```

## Usage

```sh
$ ./http2amqp -help
Usage of ./http2amqp:
  -address string
    	HTTP listen IP address (default "0.0.0.0")
  -brokeruri string
    	AMQP broker connection URI (default "amqp://guest:guest@localhost:5665/")
  -exchange string
    	AMQP broker exchange name (default "events")
  -port string
    	HTTP listen port (default "18080")
  -timeout int
    	AMQP broker queries timeout in milliseconds (default 1000)
  -verbose
    	Verbose mode, enable logging
```

## Execution example

Make sure that the environment variables are loaded before executing each command. Also that you have a RabbitMQ container up and running with those parameters.

Start the htt2amqp server in a terminal.

```sh
./http2amqp -verbose
```

Start the echo service in another terminal. This will publish a response event each time a request event arrives, with the same ID and payload (ping pong).

```sh
./echoservice
```

Make a request with a payload. The echo service will answer back if running, otherwise it will timeout.

```
curl -X GET http://localhost:18080/net/test -d 'hello world'
```

You can specify the timeout (in milliseconds) as a query param.

```
curl -X GET http://localhost:18080/net/test?timeout=200 -d 'hello world'
```

## TODO

- Test timeout parameter for each request.
- Implement delay parameter for echo server to allow tests timeouts.

## Notes

* We compile in the Dockerfile with CGO_ENABLED=0 because in the scratch image there are not some C libraries it needs. So we compile the application with them in order to run it properly. See https://go.dev/blog/cgo for more info.
