// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"net/url"

	//"github.com/aleasoluciones/simpleamqp"
	"github.com/kr/pretty"
)

type QueryMessage struct {
	id     int
	topic  string
	values map[string][]string
}

type httpDispatcher struct {
}

func (d *httpDispatcher) dispatch() {

}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	pretty.Println(r)

	topic := r.URL.Path[1:]
	queryValues, _ := url.ParseQuery(r.URL.RawQuery)

	queryMessage := QueryMessage{
		topic:  topic,
		values: queryValues,
	}

	jsonQuery, err := json.Marshal(queryMessage)

	pretty.Println("EFA2", jsonQuery, err)

	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

}

func main() {
	address := flag.String("address", "0.0.0.0", "listen address")
	port := flag.String("port", "8080", "listen port")
	//amqpuri := flag.String("amqpuri", "amqp://guest:guest@localhost/", "AMQP connection uri")
	flag.Parse()

	// amqpPublisher := simpleamqp.NewAmqpPublisher(amqpuri, "events")
	// amqpConsumer := simpleamqp.NewAmqpConsumer(amqpuri)
	// messages := amqpConsumer.Receive("events", []string{"efa1", "efa2"}, "efa", 30*time.Second)
	// for message := range messages {
	// 	log.Println(message)
	// }

	// amqpPublisher.Publish("efa2", []byte(messageBody))

	http.HandleFunc("/", handler)

	addressAndPort := fmt.Sprintf("%s:%s", *address, *port)
	log.Println("Server running in", addressAndPort, " ...")
	log.Fatal(http.ListenAndServe(addressAndPort, nil))
}
