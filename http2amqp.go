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
	//"github.com/kr/pretty"
)

const (
	MAX_INFLY = 1000
)

type QueryMessage struct {
	Id     int                 `json:"queryId"`
	Topic  string              `json:"topic"`
	Values map[string][]string `json:"values"`
}

type InputQueryMessage struct {
	queryMessage  QueryMessage
	outputChannel chan ResponseMessage
}

type ResponseMessage struct {
	id           int
	jsonResponse string
}

type httpDispatcher struct {
	input chan InputQueryMessage
}

func (d *httpDispatcher) dispatch() {
	var outputChannels [MAX_INFLY]chan ResponseMessage

	var responses chan ResponseMessage
	responses = make(chan ResponseMessage, 10000)

	id := 0
	for {

		select {
		case inputMessage := <-d.input:
			fmt.Println("Dispatch input", id, inputMessage.queryMessage, inputMessage.outputChannel)
			outputChannels[id] = inputMessage.outputChannel

			responses <- ResponseMessage{id, fmt.Sprintf("Response %s", inputMessage.queryMessage)}

			break

		case response := <-responses:
			fmt.Println("Dispatch response", response)
			outputChannels[response.id] <- response
			close(outputChannels[response.id])
		}

		id = id + 1
		if id == MAX_INFLY {
			id = 0
		}
	}
}

func handler(dispatcher *httpDispatcher, w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	topic := r.URL.Path[1:]
	queryValues, _ := url.ParseQuery(r.URL.RawQuery)

	queryMessage := QueryMessage{
		Topic:  topic,
		Values: queryValues,
	}

	jsonQuery, err := json.Marshal(queryMessage)

	fmt.Println("EFA2", string(jsonQuery), err)

	ouput := make(chan ResponseMessage)
	fmt.Println("Waiting to put input message!")
	dispatcher.input <- InputQueryMessage{queryMessage, ouput}

	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

	fmt.Println("Waiting response!")
	response := <-ouput
	fmt.Println("Response", response)
	fmt.Fprintf(w, "Hi there, I love %s!", response)
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

	dispatcher := httpDispatcher{make(chan InputQueryMessage)}
	go dispatcher.dispatch()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(&dispatcher, w, r)
	})

	addressAndPort := fmt.Sprintf("%s:%s", *address, *port)
	log.Println("Server running in", addressAndPort, " ...")
	log.Fatal(http.ListenAndServe(addressAndPort, nil))
}
