# http2amqp

[![Build Status](https://travis-ci.org/aleasoluciones/http2amqp.svg)](https://travis-ci.org/aleasoluciones/http2amqp)

## Usage
```
$ ./http2amqp --help
Usage of ./http2amqp:
  -amqpuri="amqp://guest:guest@localhost/": AMQP connection uri
  -exchange="events": AMQP exchange name
  -timeout=1: Queries timeout in seconds
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

You can specify the timeout (in seconds) for waiting for the response
```
curl -X GET http://localhost:18080/test?timeout=2 -d 'hello world'
```

## Generating new version

Update code and commit changes.
Generate a new tag and push the tag. The version will be automatically upload to [github releases](https://github.com/aleasoluciones/http2amqp/releases)

Example:
```
git tag v0.3.0
git push
git push --tags
```

Will generate the 0.3.0 version at https://github.com/aleasoluciones/http2amqp/releases/download/v0.3.0/http2amqp

## TODO
 - test timeout parameter for each request
 - implement delay parameter for echo server to allow tests timeouts
 - implement timeout in milliseconds
 
