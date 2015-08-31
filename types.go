// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package http2amqp

import (
	"net/http"
	"net/url"
)

// Request http request info to publish to amqp broker
type Request struct {
	Method string      `json:"method"`
	URL    *url.URL    `json:"url"`
	Header http.Header `json:"header"`
	Body   []byte      `json:"body"`
}

// Response info to generate a http response to return to the original http client
type Response struct {
	Status int         `json:"status"`
	Header http.Header `json:"header"`
	Body   []byte      `json:"body"`
}

// AmqpRequestMessage wrapper type for a Request including a ID and the topic to use
type AmqpRequestMessage struct {
	ID            string  `json:"id"`
	Request       Request `json:"request"`
	ResponseTopic string  `json:"responseTopic"`
}

// AmqpResponseMessage wrapper type for a Response including a ID and the response
type AmqpResponseMessage struct {
	ID       string   `json:"id"`
	Response Response `json:"response"`
}
