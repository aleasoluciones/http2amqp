# http2amqp

[![Build Status](https://travis-ci.org/aleasoluciones/http2amqp.svg)](https://travis-ci.org/aleasoluciones/http2amqp)

## Usage
```
$ ./http2amqp --help
Usage of ./http2amqp:
  -address="0.0.0.0": Listen address
  -amqpuri="amqp://guest:guest@localhost/": AMQP connection uri
  -exchange="events": AMQP exchange name
  -port="18080": Listen port
  -timeout=1000: Queries timeout in milliseconds
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

