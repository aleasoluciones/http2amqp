// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_service

import (
	"time"

	"github.com/aleasoluciones/simpleamqp"
)

func NewQueriesServiceFactory(brokerUri string, queriesExchange string, responsesExchange string, timeout time.Duration) QueriesService {
	return NewQueriesService(
		simpleamqp.NewAmqpPublisher(brokerUri, queriesExchange),
		simpleamqp.NewAmqpConsumer(brokerUri),
		NewIdsRepository(),
		responsesExchange,
		timeout)
}
