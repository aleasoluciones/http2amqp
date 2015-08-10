// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package http2amqp

import (
	"errors"
	"log"
	"time"
	"strconv"

	"encoding/json"

	"github.com/aleasoluciones/goaleasoluciones/safemap"
	"github.com/aleasoluciones/simpleamqp"
)

const (
	AMQP_RECEIVE_TIMEOUT = 30 * time.Minute
	RESPONSES_QUEUE      = "queries_responses"
	RESPONSE_TOPIC       = "queries.response"
)

func NewHttp2AmqpService(brokerUri, exchange string, timeout time.Duration) *http2amqpService {

	service := http2amqpService{
		amqpConsumer:   simpleamqp.NewAmqpConsumer(brokerUri),
		amqpPublisher:  simpleamqp.NewAmqpPublisher(brokerUri, exchange),
		idsRepository:  NewIdsRepository(),
		exchange:       exchange,
		queryTimeout:   timeout,
		queryResponses: safemap.NewSafeMap(),
	}

	go service.receiveResponses(service.amqpConsumer.Receive(
		service.exchange,
		[]string{RESPONSE_TOPIC},
		RESPONSES_QUEUE,
		simpleamqp.QueueOptions{Durable: false, Delete: true, Exclusive: true},
		AMQP_RECEIVE_TIMEOUT))

	return &service
}

type http2amqpService struct {
	amqpConsumer   simpleamqp.AMQPConsumer
	amqpPublisher  simpleamqp.AMQPPublisher
	idsRepository  IdsRepository
	exchange       string
	queryTimeout   time.Duration
	queryResponses safemap.SafeMap
}

func (service *http2amqpService) receiveResponses(amqpResponses chan simpleamqp.AmqpMessage) {
	var deserialized AmqpResponseMessage
	var value safemap.Value
	var responses chan Response
	var found bool

	for message := range amqpResponses {
		_ = json.Unmarshal([]byte(message.Body), &deserialized)

		value, found = service.queryResponses.Find(deserialized.Id)
		if found {
			responses = value.(chan Response)
			responses <- deserialized.Response
		}
	}
}

func (service *http2amqpService) publishQuery(id string, topic string, request Request) {
	serialized, _ := json.Marshal(AmqpRequestMessage{
		Id:            id,
		Request:       request,
		ResponseTopic: RESPONSE_TOPIC,
	})
	log.Println("[queries_service] Query id:", id, "topic:", topic, "request:", request)
	service.amqpPublisher.Publish(topic, serialized)
}

func (service *http2amqpService) Query(topic string, request Request) (Response, error) {
	id := service.idsRepository.Next()
	responses := make(chan Response)
	service.queryResponses.Insert(id, responses)
	defer service.queryResponses.Delete(id)

	timeout := service.queryTimeout
	for k, v := range request.URL.Query() {
	    if k == "timeout" {
	       seconds, _ := strconv.Atoi(v[0])
	       timeout = time.Duration(seconds) * time.Second
	    }
	}
	service.publishQuery(id, topic, request)

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
