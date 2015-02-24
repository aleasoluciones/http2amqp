// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aleasoluciones/goaleasoluciones/safemap"
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

		queryResponses: safemap.NewSafeMap(),
		responses:      make(chan response),
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

	queryResponses safemap.SafeMap
	responses      chan response
}

type Result interface{}

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
	var value safemap.Value
	var responses chan response
	var found bool

	for {
		select {
		case response_ := <-service.responses:
			value, found = service.queryResponses.Find(response_.Id)
			if found {
				responses = value.(chan response)
				responses <- response_
			}
		}
	}
}

func (service *queriesService) publishQuery(id Id, topic string, criteria Criteria) {
	serialized, _ := json.Marshal(amqpQueryMessage{
		Id:             id,
		CriteriaValues: criteria,
	})

	service.amqpPublisher.Publish("queries.query."+topic, serialized)
}

type Criteria map[string]string

func (service *queriesService) Query(topic string, criteria Criteria) (Result, error) {
	id := service.idsRepository.Next()
	responses := make(chan response)
	service.queryResponses.Insert(id, responses)
	defer service.queryResponses.Delete(id)
	service.publishQuery(id, topic, criteria)

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
