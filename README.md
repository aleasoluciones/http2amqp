# http2amqp

[![Build Status](https://travis-ci.org/aleasoluciones/http2amqp.svg)](https://travis-ci.org/aleasoluciones/http2amqp)

##Disclaimer
Any documentation is obsolete from the very moment it is written down, so mistrust it and check it still applies.

##Features
* It publishes an amqp message for each http request received and process the corresponding amqp responses (it waits for it) in order to answer to the original http request.
* The topic of the message it publishes comes from the URL Path of the HTTP request, using the HTTP method, network and the path replacing '/' with '.'
  * E.g.: 'get.arl.cpe'
* The TTL for the messages published is 1000 ms
* The exchange used by default is 'events'


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

