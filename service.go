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
	amqpReceiveTimeout = 30 * time.Minute
	responsesQueue     = "queries_responses"
	responseTopic      = "queries.response"
)

func NewService(brokerURI, exchange string, timeout time.Duration) *Service {

	service := Service{
		amqpConsumer:   simpleamqp.NewAmqpConsumer(brokerURI),
		amqpPublisher:  simpleamqp.NewAmqpPublisher(brokerURI, exchange),
		idsRepository:  NewIdsRepository(),
		exchange:       exchange,
		queryTimeout:   timeout,
		queryResponses: safemap.NewSafeMap(),
	}

	go service.receiveResponses(service.amqpConsumer.Receive(
		service.exchange,
		[]string{responseTopic},
		responsesQueue,
		simpleamqp.QueueOptions{Durable: false, Delete: true, Exclusive: true},
		amqpReceiveTimeout))

	return &service
}

type Service struct {
	amqpConsumer   simpleamqp.AMQPConsumer
	amqpPublisher  simpleamqp.AMQPPublisher
	idsRepository  IdsRepository
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
		_ = json.Unmarshal([]byte(message.Body), &deserialized)

		value, found = service.queryResponses.Find(deserialized.ID)
		if found {
			responses = value.(chan Response)
			responses <- deserialized.Response
		}
	}
}

func (service *Service) publishQuery(id string, topic string, request Request) {
	serialized, _ := json.Marshal(AmqpRequestMessage{
		ID:            id,
		Request:       request,
		ResponseTopic: responseTopic,
	})
	log.Println("[queries_service] Query id:", id, "topic:", topic, "request:", request)
	service.amqpPublisher.Publish(topic, serialized)
}

func (service *Service) Query(topic string, request Request) (Response, error) {
	id := service.idsRepository.Next()
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
