package mocks

import "github.com/aleasoluciones/http2amqp/queries_service"
import "github.com/stretchr/testify/mock"

type IdsRepository struct {
	mock.Mock
}

func (m *IdsRepository) Next() queries_service.Id {
	ret := m.Called()

	r0 := ret.Get(0).(queries_service.Id)

	return r0
}
