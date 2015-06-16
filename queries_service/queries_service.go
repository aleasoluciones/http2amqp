// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_service

import (
	"errors"
	"log"
	"time"

	"encoding/json"
	"net/http"
	"net/url"

	"github.com/aleasoluciones/goaleasoluciones/safemap"
	"github.com/aleasoluciones/simpleamqp"
)

const (
	AMQP_RECEIVE_TIMEOUT = 30 * time.Minute
	RESPONSES_QUEUE      = "queries_responses"
)

type Request struct {
	Method string
	URL    *url.URL
	Header http.Header
	Body   []byte
}

type Response struct {
	Status int
	Header http.Header
	Body   []byte
}

type amqpRequestMessage struct {
	Id      Id      `json:"id"`
	Request Request `json:"request"`
}

type amqpResponseMessage struct {
	Id       Id       `json:"id"`
	Response Response `json:"response"`
}

type QueriesService interface {
	Query(topic string, request Request) (Response, error)
}

func NewQueriesService(amqpPublisher simpleamqp.AMQPPublisher, amqpConsumer simpleamqp.AMQPConsumer, idsRepository IdsRepository, exchange string, timeout time.Duration) QueriesService {
	service := queriesService{
		amqpConsumer:   amqpConsumer,
		amqpPublisher:  amqpPublisher,
		idsRepository:  idsRepository,
		exchange:       exchange,
		queryTimeout:   timeout,
		queryResponses: safemap.NewSafeMap(),
	}

	go service.receiveResponses()

	return &service
}

type queriesService struct {
	amqpConsumer   simpleamqp.AMQPConsumer
	amqpPublisher  simpleamqp.AMQPPublisher
	idsRepository  IdsRepository
	exchange       string
	queryTimeout   time.Duration
	queryResponses safemap.SafeMap
}

func (service *queriesService) receiveResponses() {
	amqpResponses := service.amqpConsumer.Receive(
		service.exchange,
		[]string{"queries.response"},
		RESPONSES_QUEUE,
		simpleamqp.QueueOptions{Durable: false, Delete: true, Exclusive: true},
		AMQP_RECEIVE_TIMEOUT)

	var deserialized amqpResponseMessage
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

func (service *queriesService) publishQuery(id Id, topic string, request Request) {
	serialized, _ := json.Marshal(amqpRequestMessage{
		Id:      id,
		Request: request,
	})
	log.Println("[queries_service] Query id:", id, "topic:", topic, "request:", request)
	service.amqpPublisher.Publish(topic, serialized)
}

func (service *queriesService) Query(topic string, request Request) (Response, error) {
	id := service.idsRepository.Next()
	responses := make(chan Response)
	service.queryResponses.Insert(id, responses)
	defer service.queryResponses.Delete(id)
	service.publishQuery(id, topic, request)

	timeoutTicker := time.NewTicker(service.queryTimeout)
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
