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
	Id           int
	JsonResponse string
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
	messages := d.amqpConsumer.Receive(RESPONSE_EXCHANGE, []string{"#"}, "responses_q", 30*time.Minute)
	for message := range messages {
		log.Println("RECEIVE1", message)

		var m ResponseMessage
		err := json.Unmarshal([]byte(message.Body), &m)
		if err != nil {
			log.Println("R2 Error", err, message.Body)
		} else {
			// FIXME extract id
			d.responses <- ResponseMessage{m.Id, string(message.Body)}
		}
	}
}

func (d *HttpDispatcher) dispatch() {
	var outputChannels [MAX_INFLY]chan ResponseMessage

	id := 0
	for {
		log.Println("Antes select")
		select {
		case inputMessage := <-d.input:
			log.Println("D1 Dispatch input", id, inputMessage.queryMessage, inputMessage.outputChannel)
			outputChannels[id] = inputMessage.outputChannel
			inputMessage.queryMessage.Id = id
			jsonQuery, _ := json.Marshal(inputMessage.queryMessage)
			log.Println("D2 Dispatch input")
			d.amqpPublisher.Publish(inputMessage.queryMessage.Topic, []byte(jsonQuery))
			break

		case response := <-d.responses:
			log.Println("Dispatch response", response, response.Id)
			if outputChannels[response.Id] != nil {
				outputChannels[response.Id] <- response
			}
			break
		}

		log.Println("Despues select")

		id = id + 1
		if id == MAX_INFLY {
			id = 0
		}
	}
}

func (dispatcher *HttpDispatcher) httpHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("H1", r.URL.Path[1:])

	topic := r.URL.Path[1:]
	queryValues, _ := url.ParseQuery(r.URL.RawQuery)

	queryMessage := QueryMessage{
		Topic:  topic,
		Values: queryValues,
	}

	jsonQuery, err := json.Marshal(queryMessage)

	log.Println("H2", string(jsonQuery), err)

	ouput := make(chan ResponseMessage)
	log.Println("H3")
	dispatcher.input <- InputQueryMessage{queryMessage, ouput}

	log.Println("H4")
	response := <-ouput
	log.Println("H5 Response", response, response.Id)

	fmt.Fprintf(w, "Hi there, I love")
}
