# http2amqp

[![Build Status](https://travis-ci.org/aleasoluciones/http2amqp.svg)](https://travis-ci.org/aleasoluciones/http2amqp)

## Profiling

The HTTP server exposes an endpoint to access [profiling data](http://golang.org/pkg/net/http/pprof/).

```
go tool pprof http2amqp http://localhost:18080/debug/pprof/heap
```

To access from a remote machine you can perform a [ssh tunnel](http://en.wikipedia.org/wiki/Tunneling_protocol#Secure_Shell_tunneling).

```
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no <your-remote-host> -L 18080:localhost:18080 -N
```

And then:

```
go tool pprof http2amqp http://localhost:18080/debug/pprof/heap
```
