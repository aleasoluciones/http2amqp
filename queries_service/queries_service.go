// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_service

import (
	"encoding/json"
	"errors"
	"log"
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
		queryResults:  safemap.NewSafeMap(),
	}

	go service.receiveResponses()

	return &service
}

type queriesService struct {
	amqpConsumer  simpleamqp.AMQPConsumer
	amqpPublisher simpleamqp.AMQPPublisher
	idsRepository IdsRepository
	exchange      string
	queryTimeout  time.Duration
	queryResults  safemap.SafeMap
}

type Result interface{}

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
	var value safemap.Value
	var results chan Result
	var found bool

	for message := range amqpResponses {
		_ = json.Unmarshal([]byte(message.Body), &deserialized)

		value, found = service.queryResults.Find(deserialized.Id)
		if found {
			results = value.(chan Result)
			results <- deserialized.Result
		}
	}
}

func (service *queriesService) publishQuery(id Id, topic string, criteria Criteria) {
	serialized, _ := json.Marshal(amqpQueryMessage{
		Id:             id,
		CriteriaValues: criteria,
	})

	log.Println("[queries_service] Query id:", id, "topic:", topic, "criteria:", criteria)

	service.amqpPublisher.Publish("queries.query."+topic, serialized)
}

type Criteria map[string]string

func (service *queriesService) Query(topic string, criteria Criteria) (Result, error) {
	id := service.idsRepository.Next()
	results := make(chan Result)
	service.queryResults.Insert(id, results)
	defer service.queryResults.Delete(id)
	service.publishQuery(id, topic, criteria)

	timeoutTicker := time.NewTicker(service.queryTimeout)
	defer timeoutTicker.Stop()
	afterTimeout := timeoutTicker.C

	select {
	case result := <-results:
		return result, nil
	case <-afterTimeout:
		log.Println("[queries_service] Timeout for query id:", id)
		return nil, errors.New("Timeout")
	}
}
