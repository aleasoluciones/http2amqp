// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aleasoluciones/simpleamqp"
)

const (
	AMQP_RECEIVE_TIMEOUT = 30 * time.Minute
	RESPONSES_QUEUE      = "queries_responses"
)

type QueriesService interface {
	Query(topic string, criteria Criteria) (Result, error)
}

func NewQueriesService(amqpPublisher simpleamqp.AMQPPublisher, amqpConsumer simpleamqp.AMQPConsumer, idsRepository IdsRepository, exchange string, timeout time.Duration) QueriesService {
	service := queriesService{
		amqpConsumer:  amqpConsumer,
		amqpPublisher: amqpPublisher,
		idsRepository: idsRepository,
		exchange:      exchange,
		queryTimeout:  timeout,

		queries:   make(chan query),
		responses: make(chan response),
	}

	go service.dispatch()
	go service.receiveResponses()

	return &service
}

type queriesService struct {
	amqpConsumer  simpleamqp.AMQPConsumer
	amqpPublisher simpleamqp.AMQPPublisher
	idsRepository IdsRepository
	exchange      string
	queryTimeout  time.Duration

	queries   chan query
	responses chan response
}

type Result interface{}

type query struct {
	RoutingKey     string
	CriteriaValues Criteria
	Responses      chan response
}

type response struct {
	Id     Id
	Result Result
}

type amqpQueryMessage struct {
	Id             Id       `json:"id"`
	CriteriaValues Criteria `json:"criteria"`
}

type amqpResponseMessage struct {
	Id     Id     `json:"id"`
	Result Result `json:"result"`
}

func (service *queriesService) receiveResponses() {
	amqpResponses := service.amqpConsumer.Receive(
		service.exchange,
		[]string{"queries.response"},
		RESPONSES_QUEUE,
		simpleamqp.QueueOptions{Durable: false, Delete: true, Exclusive: true},
		AMQP_RECEIVE_TIMEOUT)

	var deserialized amqpResponseMessage

	for message := range amqpResponses {
		_ = json.Unmarshal([]byte(message.Body), &deserialized)

		service.responses <- response{
			Id:     deserialized.Id,
			Result: deserialized.Result,
		}
	}
}

func (service *queriesService) dispatch() {
	var id Id
	var responses chan response
	var found bool

	queryResponses := map[Id]chan response{}

	for {
		select {
		case query := <-service.queries:
			id = service.idsRepository.Next()
			queryResponses[id] = query.Responses
			service.publishQuery(id, query)
		case response := <-service.responses:
			responses, found = queryResponses[response.Id]
			if found {
				responses <- response
			}
		}
	}
}

func (service *queriesService) publishQuery(id Id, query query) {
	serialized, _ := json.Marshal(amqpQueryMessage{
		Id:             id,
		CriteriaValues: query.CriteriaValues,
	})

	service.amqpPublisher.Publish("queries.query."+query.RoutingKey, serialized)
}

type Criteria map[string]string

func (service *queriesService) Query(topic string, criteria Criteria) (Result, error) {
	responses := make(chan response)
	service.queries <- query{
		RoutingKey:     topic,
		CriteriaValues: criteria,
		Responses:      responses,
	}

	timeoutTicker := time.NewTicker(service.queryTimeout)
	defer timeoutTicker.Stop()
	afterTimeout := timeoutTicker.C

	select {
	case response := <-responses:
		return response.Result, nil
	case <-afterTimeout:
		return nil, errors.New("Timeout")
	}
}
