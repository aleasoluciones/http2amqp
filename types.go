// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package http2amqp

import (
	"net/http"
	"net/url"
)

type Request struct {
	Method string      `json:"method"`
	URL    *url.URL    `json:"url"`
	Header http.Header `json:"header"`
	Body   []byte      `json:"body"`
}

type Response struct {
	Status int         `json:"status"`
	Header http.Header `json:"header"`
	Body   []byte      `json:"body"`
}

type AmqpRequestMessage struct {
	Id            string  `json:"id"`
	Request       Request `json:"request"`
	ResponseTopic string  `json:"responseTopic"`
}

type AmqpResponseMessage struct {
	Id       string   `json:"id"`
	Response Response `json:"response"`
}
