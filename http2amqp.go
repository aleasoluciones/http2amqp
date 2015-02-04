// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package http2amqp

import (
	"fmt"
	"log"
	"time"

	"encoding/json"
	"net/http"
	"net/url"

	"github.com/aleasoluciones/simpleamqp"
)

const (
	MAX_INFLY         = 1000
	QUERY_EXCHANGE    = "queries"
	RESPONSE_EXCHANGE = "responses"
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

type HttpDispatcher struct {
	input         chan InputQueryMessage
	responses     chan ResponseMessage
	amqpPublisher simpleamqp.AmqpPublisher
	amqpConsumer  simpleamqp.AmqpConsumer
}

func NewHttpDispatcher(amqpuri string) *HttpDispatcher {
	dispatcher := HttpDispatcher{
		input:         make(chan InputQueryMessage),
		responses:     make(chan ResponseMessage, 10000),
		amqpPublisher: *simpleamqp.NewAmqpPublisher(amqpuri, QUERY_EXCHANGE),
		amqpConsumer:  *simpleamqp.NewAmqpConsumer(amqpuri),
	}

	return &dispatcher
}

func (dispatcher *HttpDispatcher) ListenAndServe(addressAndPort string) error {
	go dispatcher.dispatch()
	go dispatcher.receiveResponses()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dispatcher.httpHandle(w, r)
	})

	log.Println("Server running in", addressAndPort, " ...")
	return http.ListenAndServe(addressAndPort, nil)
}

func (d *HttpDispatcher) receiveResponses() {
	messages := d.amqpConsumer.Receive(RESPONSE_EXCHANGE, []string{"#"}, "responses_q", 30*time.Second)
	for message := range messages {
		log.Println("RECEIVE1", message)

		// FIXME extract id
		d.responses <- ResponseMessage{0, string(message.Body)}
	}
}

func (d *HttpDispatcher) dispatch() {
	var outputChannels [MAX_INFLY]chan ResponseMessage

	id := 0
	for {

		select {
		case inputMessage := <-d.input:
			fmt.Println("Dispatch input", id, inputMessage.queryMessage, inputMessage.outputChannel)
			outputChannels[id] = inputMessage.outputChannel
			inputMessage.queryMessage.Id = id
			jsonQuery, _ := json.Marshal(inputMessage.queryMessage)
			d.amqpPublisher.Publish(inputMessage.queryMessage.Topic, []byte(jsonQuery))
			break

		case response := <-d.responses:
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

func (dispatcher *HttpDispatcher) httpHandle(w http.ResponseWriter, r *http.Request) {
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
