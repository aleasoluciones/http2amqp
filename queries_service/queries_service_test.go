// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_service_test

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/aleasoluciones/http2amqp/queries_service"
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
			amqpPublisher  *mocks.AMQPPublisher
			queriesService QueriesService
		)

		BeforeEach(func() {
			amqpResponses = make(chan simpleamqp.AmqpMessage)
			amqpConsumer := new(mocks.AMQPConsumer)
			amqpConsumer.On("Receive", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(amqpResponses)

			amqpPublisher = new(mocks.AMQPPublisher)
			amqpPublisher.On("Publish", mock.Anything, mock.Anything).Return(nil)

			queriesService = NewQueriesService(amqpPublisher, amqpConsumer, A_QUERY_TIMEOUT)
		})

		Context("Response management", func() {
			It("returns response when a response for the query has been received", func() {
				go func() {
					time.Sleep(A_QUERY_TIMEOUT / time.Duration(2))
					amqpResponses <- newAmqpResponse(0, []string{"foo", "bar"})
				}()

				result, err := queriesService.Query(A_TOPIC, Criteria{"q": "foo"})
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(ConsistOf("foo", "bar"))
			})

			It("returns timeout error when response for the query has not been received in time", func() {
				go func() {
					time.Sleep(A_QUERY_TIMEOUT * time.Duration(3))
					amqpResponses <- newAmqpResponse(0, []string{"foo", "bar"})
				}()

				result, err := queriesService.Query(A_TOPIC, Criteria{"q": "foo"})

				Expect(err).To(MatchError("Timeout"))
				Expect(result).To(BeNil())
			})

			It("returns timeout error when response for another query", func() {
				go func() {
					time.Sleep(A_QUERY_TIMEOUT / time.Duration(2))
					amqpResponses <- newAmqpResponse(1, []string{"foo", "bar"})
				}()

				result, err := queriesService.Query(A_TOPIC, Criteria{"q": "foo"})

				Expect(err).To(MatchError("Timeout"))
				Expect(result).To(BeNil())
			})
		})

		Context("Query management", func() {
			It("publish the query to amqp", func() {
				// Recibe una petición =>
				// crear ueryMessage
				// Publicarlo

				_, _ = queriesService.Query(A_TOPIC, Criteria{"q": "foo"})

				expectedCriteriaJson := fmt.Sprintf(`{"%s":"%s"}`, "q", "foo")
				expectedQueryJson := fmt.Sprintf(`{"id":%d,"criteria":%s}`, 0, expectedCriteriaJson)
				amqpPublisher.AssertCalled(GinkgoT(), "Publish", A_TOPIC, []byte(expectedQueryJson))
			})
		})
	})
})

func newAmqpResponse(id int, response interface{}) simpleamqp.AmqpMessage {
	serialized, _ := json.Marshal(map[string]interface{}{
		"id":      id,
		"content": response,
	})
	return simpleamqp.AmqpMessage{Body: string(serialized)}
}
