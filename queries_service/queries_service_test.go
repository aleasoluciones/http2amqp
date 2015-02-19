// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_service_test

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aleasoluciones/http2amqp/queries_service/mocks"
	"github.com/aleasoluciones/simpleamqp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

const (
	A_TOPIC         = "a-topic"
	A_QUERY_TIMEOUT = 10 * time.Millisecond
)

var _ = Describe("Queries service", func() {
	Describe("Query", func() {
		var (
			amqpResponses  chan simpleamqp.AmqpMessage
			queriesService QueriesService
		)

		BeforeEach(func() {
			amqpResponses = make(chan simpleamqp.AmqpMessage)
			amqpConsumer := new(mocks.AMQPConsumer)
			amqpConsumer.On("Receive", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(amqpResponses)

			queriesService = NewQueriesService(amqpConsumer, A_QUERY_TIMEOUT)
		})

		It("returns response when a response for the query has been received", func() {
			go func() {
				time.Sleep(A_QUERY_TIMEOUT / time.Duration(2))
				amqpResponses <- newAmqpResponse([]string{"foo", "bar"})
			}()

			response := queriesService.Query(A_TOPIC, map[string]interface{}{"q": "foo"})

			Expect(response.Error).NotTo(HaveOccurred())
			Expect(response.Content).To(ConsistOf("foo", "bar"))
		})

		It("returns timeout error when no response for the query has been received in time", func() {
			go func() {
				time.Sleep(A_QUERY_TIMEOUT * time.Duration(2))
				amqpResponses <- newAmqpResponse([]string{"foo", "bar"})
			}()

			response := queriesService.Query(A_TOPIC, map[string]interface{}{"q": "foo"})

			Expect(response.Error).To(MatchError("Timeout"))
			Expect(response.Content).To(BeNil())
		})
	})
})

func NewQueriesService(amqpConsumer simpleamqp.AMQPConsumer, timeout time.Duration) QueriesService {
	return &queriesServiceT{
		amqpConsumer: amqpConsumer,
		queryTimeout: timeout,
	}
}

type QueriesService interface {
	Query(topic string, criteria map[string]interface{}) queryResponse
}

type queriesServiceT struct {
	amqpConsumer simpleamqp.AMQPConsumer
	queryTimeout time.Duration
}

func (service *queriesServiceT) Query(topic string, criteria map[string]interface{}) queryResponse {
	var content interface{}
	var err error

	amqpResponses := service.amqpConsumer.Receive("foo exchange", []string{"foo routing key"}, " no se", simpleamqp.QueueOptions{}, 10)

	afterTimeoutTicker := time.NewTicker(service.queryTimeout)
	defer afterTimeoutTicker.Stop()
	afterTimeout := afterTimeoutTicker.C

	select {
	case amqpResponse := <-amqpResponses:
		_ = json.Unmarshal([]byte(amqpResponse.Body), &content)
	case <-afterTimeout:
		err = errors.New("Timeout")
	}

	return queryResponse{Content: content, Error: err}
}

type queryResponse struct {
	Content interface{}
	Error   error
}

func newAmqpResponse(response interface{}) simpleamqp.AmqpMessage {
	serialized, _ := json.Marshal(response)
	return simpleamqp.AmqpMessage{Body: string(serialized)}
}
