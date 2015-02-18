package queries_service_test

import (
	"encoding/json"
	"github.com/aleasoluciones/http2amqp/queries_service/mocks"
	"github.com/aleasoluciones/simpleamqp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"time"
)

const (
	A_TOPIC = "a-topic"
)

var _ = Describe("Queries service", func() {
	It("query returns response when a response for the query has been received", func() {
		// eventPublisher.publish A_TPIC, values
		// waits for a respons
		amqpResponses := make(chan simpleamqp.AmqpMessage)
		amqpConsumer := new(mocks.AMQPConsumer)
		amqpConsumer.On("Receive", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(amqpResponses)

		go func() {
			time.Sleep(1 * time.Second)
			amqpResponses <- newAmqpResponse([]string{"foo", "bar"})
		}()

		queriesService := NewQueriesService(amqpConsumer)
		response := queriesService.Query(A_TOPIC, map[string]interface{}{"q": "foo"})

		Expect(response.Content).To(ConsistOf("foo", "bar"))
	})
})

func NewQueriesService(amqpConsumer simpleamqp.AMQPConsumer) QueriesService {
	return &queriesServiceT{amqpConsumer: amqpConsumer}
}

type QueriesService interface {
	Query(topic string, criteria map[string]interface{}) queryResponse
}

type queriesServiceT struct {
	amqpConsumer simpleamqp.AMQPConsumer
}

func (service *queriesServiceT) Query(topic string, criteria map[string]interface{}) queryResponse {
	amqpResponses := service.amqpConsumer.Receive("foo exchange", []string{"foo routing key"}, " no se", simpleamqp.QueueOptions{}, 10)
	amqpResponse := <-amqpResponses
	var content interface{}
	_ = json.Unmarshal([]byte(amqpResponse.Body), &content)
	return queryResponse{Content: content}
}

type queryResponse struct {
	Content interface{}
}

func newAmqpResponse(response interface{}) simpleamqp.AmqpMessage {
	serialized, _ := json.Marshal(response)
	return simpleamqp.AmqpMessage{Body: string(serialized)}
}
