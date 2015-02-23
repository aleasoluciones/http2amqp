// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/aleasoluciones/http2amqp/queries_service"
)

func NewHTTPServer(queriesService queries_service.QueriesService) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		topic := topicFor(r)
		criteria := criteriaFor(r)
		log.Println("[http2amqp] Query", topic, criteria)

		result, err := queriesService.Query(topic, criteria)

		if err != nil {
			log.Println("[http2amqp] Query error", err)
			http.Error(w, err.Error(), 500)
			return
		}

		serialized, err := json.Marshal(result)

		if err != nil {
			log.Println("[http2amqp] Error marshaling query result", err)
			http.Error(w, "Internal Server Error", 500)
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

func criteriaFor(r *http.Request) queries_service.Criteria {
	criteria, _ := url.ParseQuery(r.URL.RawQuery)
	log.Println("JGIL", criteria)

	return queries_service.Criteria{}
}
