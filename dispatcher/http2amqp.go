// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aleasoluciones/http2amqp"
)

func main() {
	address := flag.String("address", "0.0.0.0", "listen address")
	port := flag.String("port", "8080", "listen port")
	amqpuri := flag.String("amqpuri", "amqp://guest:guest@localhost/", "AMQP connection uri")
	flag.Parse()

	httpDispatcher := http2amqp.NewHttpDispatcher(*amqpuri)
	addressAndPort := fmt.Sprintf("%s:%s", *address, *port)
	log.Fatal(httpDispatcher.ListenAndServe(addressAndPort))
}
