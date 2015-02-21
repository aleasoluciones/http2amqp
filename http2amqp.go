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
	maxInFly         = 1000
	queryExchange    = "queries"
	responseExchange = "responses"
)

type QueryMessage struct {
	ID     int                 `json:"queryId"`
	Topic  string              `json:"topic"`
	Values map[string][]string `json:"values"`
}

type inputQueryMessage struct {
	QueryMessage  QueryMessage
	outputChannel chan ResponseMessage
}

type ResponseMessage struct {
	ID           int
	JsonResponse string
}

// HTTPDispatcher struct maintains the channels and configuration for the server
type HTTPDispatcher struct {
	input         chan inputQueryMessage
	responses     chan ResponseMessage
	amqpPublisher simpleamqp.AmqpPublisher
	amqpConsumer  simpleamqp.AmqpConsumer
	timeout       time.Duration
}

// NewHTTPDispatcher return a new http that can dispatch the http queries to a amqp backends
func NewHTTPDispatcher(amqpuri string, timeout time.Duration) *HTTPDispatcher {
	dispatcher := HTTPDispatcher{
		input:         make(chan inputQueryMessage),
		responses:     make(chan ResponseMessage, 10000),
		amqpPublisher: *simpleamqp.NewAmqpPublisher(amqpuri, queryExchange),
		amqpConsumer:  *simpleamqp.NewAmqpConsumer(amqpuri),
		timeout:       timeout,
	}

	return &dispatcher
}

// ListenAndServe bind the http server to the given address and port and start to handle the requests
func (d *HTTPDispatcher) ListenAndServe(addressAndPort string) error {
	go d.dispatch()
	go d.receiveResponses()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d.httpHandle(w, r)
	})

	log.Println("Server running in", addressAndPort, " ...")
	return http.ListenAndServe(addressAndPort, nil)
}

func (d *HTTPDispatcher) receiveResponses() {
	messages := d.amqpConsumer.Receive(responseExchange,
		[]string{"#"},
		"", simpleamqp.QueueOptions{Durable: false, Delete: true, Exclusive: true},
		30*time.Minute)
	for message := range messages {
		log.Println("RECEIVE1", message)

		var m ResponseMessage
		err := json.Unmarshal([]byte(message.Body), &m)
		if err != nil {
			log.Println("R2 Error", err, message.Body)
		} else {
			// FIXME extract id
			d.responses <- ResponseMessage{m.ID, string(message.Body)}
		}
	}
}

func (d *HTTPDispatcher) dispatch() {
	var outputChannels [maxInFly]chan ResponseMessage

	id := 0
	for {
		log.Println("Antes select")
		select {
		case inputMessage := <-d.input:
			log.Println("D1 Dispatch input", id, inputMessage.QueryMessage, inputMessage.outputChannel)
			outputChannels[id] = inputMessage.outputChannel
			inputMessage.QueryMessage.ID = id
			jsonQuery, _ := json.Marshal(inputMessage.QueryMessage)
			log.Println("D2 Dispatch input")
			d.amqpPublisher.Publish(inputMessage.QueryMessage.Topic, []byte(jsonQuery))
			break

		case response := <-d.responses:
			log.Println("Dispatch response", response, response.ID)
			if outputChannels[response.ID] != nil {
				select {
				case outputChannels[response.ID] <- response:
					break
				default:
					log.Println("Response discarded (timeout)", response)
				}
			}
			break
		}

		log.Println("Despues select")

		id = id + 1
		if id == maxInFly {
			id = 0
		}
	}
}

func (d *HTTPDispatcher) httpHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("H1", r.URL.Path[1:])

	topic := r.URL.Path[1:]
	queryValues, _ := url.ParseQuery(r.URL.RawQuery)

	QueryMessage := QueryMessage{
		Topic:  topic,
		Values: queryValues,
	}

	jsonQuery, err := json.Marshal(QueryMessage)

	log.Println("H2", string(jsonQuery), err)

	ouput := make(chan ResponseMessage)
	log.Println("H3")
	d.input <- inputQueryMessage{QueryMessage, ouput}

	log.Println("H4")

	timeoutTimer := time.NewTimer(d.timeout)
	defer timeoutTimer.Stop()
	afterTimeout := timeoutTimer.C

	select {
	case response := <-ouput:
		log.Println("Wrinting response")
		fmt.Fprintf(w, response.JsonResponse)
		break
	case <-afterTimeout:
		log.Println("Wrinting timeout error")
		http.Error(w, "Internal timeout Error!", 500)
		log.Println("Error")
		break
	}

}
