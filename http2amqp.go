// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/kr/pretty"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	pretty.Println(r)
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	address := flag.String("address", "0.0.0.0", "listen address")
	port := flag.String("port", "8080", "listen port")
	flag.Parse()

	http.HandleFunc("/", handler)

	addressAndPort := fmt.Sprintf("%s:%s", *address, *port)
	log.Println("Server running in", addressAndPort, " ...")
	log.Fatal(http.ListenAndServe(addressAndPort, nil))
}
