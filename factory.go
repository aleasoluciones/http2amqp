// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package http2amqp

import (
	"time"

	"github.com/aleasoluciones/simpleamqp"
)

func NewQueriesServiceFactory(brokerUri string, exchange string, timeout time.Duration) QueriesService {
	return NewQueriesService(
		simpleamqp.NewAmqpPublisher(brokerUri, exchange),
		simpleamqp.NewAmqpConsumer(brokerUri),
		NewIdsRepository(),
		exchange,
		timeout)
}
