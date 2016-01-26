// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package http2amqp

import (
	"errors"
	"log"
	"strconv"
	"time"

	"encoding/json"

	"github.com/aleasoluciones/goaleasoluciones/safemap"
	"github.com/aleasoluciones/simpleamqp"
)

const (
	responsesQueue = "queries_responses"
	responseTopic  = "queries.response"
)

// NewService return the http2amqp service. This service publish a amqp message for each http request
// and process the corresponding amqp responses to answer to the original http request
func NewService(brokerURI, exchange string, timeout time.Duration) *Service {

	service := Service{
		amqpConsumer:   simpleamqp.NewAmqpConsumer(brokerURI),
		amqpPublisher:  simpleamqp.NewAmqpPublisher(brokerURI, exchange),
		idsGenerator:   NewUUIDIdsGenerator(),
		exchange:       exchange,
		queryTimeout:   timeout,
		queryResponses: safemap.NewSafeMap(),
	}

	go service.receiveResponses(service.amqpConsumer.ReceiveWithoutTimeout(
		service.exchange,
		[]string{responseTopic},
		responsesQueue,
		simpleamqp.QueueOptions{Durable: false, Delete: true, Exclusive: true}))

	return &service
}

// Service http2amqp service
type Service struct {
	amqpConsumer   simpleamqp.AMQPConsumer
	amqpPublisher  simpleamqp.AMQPPublisher
	idsGenerator   IdsGenerator
	exchange       string
	queryTimeout   time.Duration
	queryResponses safemap.SafeMap
}

func (service *Service) receiveResponses(amqpResponses chan simpleamqp.AmqpMessage) {
	var deserialized AmqpResponseMessage
	var value safemap.Value
	var responses chan Response
	var found bool

	for message := range amqpResponses {
		log.Println("Response received")
		err := json.Unmarshal([]byte(message.Body), &deserialized)

		if err != nil {
			log.Println("Error unmarshaling json response:", err)
			continue
		}

		log.Println("Response received ID:", deserialized.ID)
		value, found = service.queryResponses.Find(deserialized.ID)
		if found {
			log.Println("Pending request found for", deserialized.ID)
			responses = value.(chan Response)
			responses <- deserialized.Response
		}
	}
}

func (service *Service) publishQuery(id string, topic string, request Request, ttl time.Duration) {
	serialized, _ := json.Marshal(AmqpRequestMessage{
		ID:            id,
		Request:       request,
		ResponseTopic: responseTopic,
	})
	log.Println("[queries_service] Query id:", id, "topic:", topic, "request:", request, "ttl:", ttl)
	service.amqpPublisher.PublishWithTTL(topic, serialized, durationToMilliseconds(ttl))
}

func durationToMilliseconds(value time.Duration) int {
	return int(value.Nanoseconds() / 1000000)
}

// DispatchHTTPRequest process a request. Send the request to the broker using the
// given topic and wait for the response (or the timeout)
func (service *Service) DispatchHTTPRequest(topic string, request Request) (Response, error) {
	id := service.idsGenerator.Next()
	responses := make(chan Response)
	service.queryResponses.Insert(id, responses)
	defer service.queryResponses.Delete(id)

	timeout := service.queryTimeout
	for k, v := range request.URL.Query() {
		if k == "timeout" {
			milliseconds, _ := strconv.Atoi(v[0])
			timeout = time.Duration(milliseconds) * time.Millisecond
		}
	}

	service.publishQuery(id, topic, request, timeout)
	log.Println("Request published", id)

	timeoutTicker := time.NewTicker(timeout)
	defer timeoutTicker.Stop()
	afterTimeout := timeoutTicker.C

	select {
	case response := <-responses:
		return response, nil
	case <-afterTimeout:
		log.Println("[queries_service] Timeout for query id:", id)
		return Response{}, errors.New("Timeout")
	}
}
