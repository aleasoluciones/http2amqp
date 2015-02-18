package mocks

import "github.com/stretchr/testify/mock"

type AMQPPublisher struct {
	mock.Mock
}

func (m *AMQPPublisher) Publish(routingKey string, message []byte) {
	m.Called(routingKey, message)
}
