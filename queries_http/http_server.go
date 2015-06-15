// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_http

import (
	"log"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aleasoluciones/http2amqp/queries_service"
)

func NewHTTPServer(queriesService queries_service.QueriesService) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		topic := topicFor(r)

		request := queries_service.Request{
			Method: r.Method,
			URL:    r.URL,
			Header: r.Header,
		}

		var err error
		request.Body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("[http2amqp] Error reading body", err)
			newJsonError(w, err.Error(), 400)
			return
		}

		result, err := queriesService.Query(topic, request)

		if err != nil {
			newJsonError(w, err.Error(), 404)
			return
		}

		serialized, err := json.Marshal(jsonResponse{Result: result})

		if err != nil {
			log.Println("[http2amqp] Error marshaling query result", err)
			newJsonError(w, "Internal Server Error", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(serialized)
	})

	log.Println("[http2amqp] Starting HTTP server at 127.0.0.1:18080 ...")
	http.ListenAndServe("127.0.0.1:18080", nil)
}

func topicFor(r *http.Request) string {
	return r.URL.Path[1:]
}

func newJsonError(w http.ResponseWriter, message string, status int) {
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

type jsonResponse struct {
	Result queries_service.Result `json:"result"`
}
