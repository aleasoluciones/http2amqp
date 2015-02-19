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

type Result interface{}

type QueriesService interface {
	Query(topic string, criteria Criteria) (Result, error)
}

func NewQueriesService(amqpPublisher simpleamqp.AMQPPublisher, amqpConsumer simpleamqp.AMQPConsumer, timeout time.Duration) QueriesService {
	service := queriesService{
		amqpConsumer:  amqpConsumer,
		amqpPublisher: amqpPublisher,
		queryTimeout:  timeout,

		queries:   make(chan query),
		responses: make(chan responseMessage),
	}

	go service.dispatch()
	go service.receiveResponses()

	return &service
}

type queriesService struct {
	amqpConsumer  simpleamqp.AMQPConsumer
	amqpPublisher simpleamqp.AMQPPublisher
	queryTimeout  time.Duration

	queries   chan query
	responses chan responseMessage
}

type query struct {
	Id              int
	RoutingKey      string
	CriteriaMessage Criteria
	Responses       chan responseMessage
}
type amqpQueryMessage struct {
	Id              int      `json:"id"`
	CriteriaMessage Criteria `json:"criteria"`
}

type responseMessage struct {
	Id      int
	Message interface{}
}

func (service *queriesService) receiveResponses() {
	amqpResponses := service.amqpConsumer.Receive("an exchange", []string{"a topic"}, "no se", simpleamqp.QueueOptions{}, 1*time.Minute)

	var deserialized map[string]interface{}

	for message := range amqpResponses {
		_ = json.Unmarshal([]byte(message.Body), &deserialized)

		if deserialized["id"] != nil {
			service.responses <- responseMessage{
				Id:      int(deserialized["id"].(float64)),
				Message: deserialized["content"],
			}
		}
	}
}

func (service *queriesService) dispatch() {
	var responses chan responseMessage
	var found bool

	queryResponses := map[int]chan responseMessage{}
	id := 0

	for {
		select {
		case query := <-service.queries:
			query.Id = id
			queryResponses[id] = query.Responses
			messageToSend := amqpQueryMessage{Id: id, CriteriaMessage: query.CriteriaMessage}
			messageToSendJson, _ := json.Marshal(messageToSend)
			service.amqpPublisher.Publish(query.RoutingKey, messageToSendJson)
			id += 1
		case response := <-service.responses:
			responses, found = queryResponses[response.Id]
			if found {
				responses <- response
			}
		}
	}
}

type Criteria map[string]string

func (service *queriesService) Query(topic string, criteria Criteria) (Result, error) {
	responses := make(chan responseMessage)
	service.queries <- query{
		RoutingKey:      topic,
		CriteriaMessage: criteria,
		Responses:       responses,
	}

	afterTimeoutTicker := time.NewTicker(service.queryTimeout)
	defer afterTimeoutTicker.Stop()
	afterTimeout := afterTimeoutTicker.C

	select {
	case response := <-responses:
		return response.Message, nil
	case <-afterTimeout:
		return nil, errors.New("Timeout")
	}
}
