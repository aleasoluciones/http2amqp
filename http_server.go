// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package http2amqp

import (
	"log"
	"strings"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

func NewHTTPServerFunc(HTTP2amqpService *HTTP2amqpService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := topicFor(r)

		request := Request{
			Method: r.Method,
			URL:    r.URL,
			Header: r.Header,
		}

		var err error
		request.Body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("[http2amqp] Error reading body", err)
			newJSONError(w, err.Error(), 400)
			return
		}

		response, err := HTTP2amqpService.Query(topic, request)

		if err != nil {
			newJSONError(w, err.Error(), 404)
			return
		}

		for header := range response.Header {
			w.Header().Set(header, response.Header.Get(header))
		}
		w.WriteHeader(response.Status)
		w.Write(response.Body)
	}

}

func topicFor(r *http.Request) string {
	return strings.ToLower(r.Method) + "." + strings.Replace(r.URL.Path[1:], "/", ".", -1)
}

func newJSONError(w http.ResponseWriter, message string, status int) {
	serialized, err := json.Marshal(jsonError{Status: status, Error: message})

	if err != nil {
		log.Println("[http2amqp] Error marshaling error message", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(serialized)
}

type jsonError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
