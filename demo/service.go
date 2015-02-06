package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aleasoluciones/http2amqp"
	"github.com/aleasoluciones/simpleamqp"
)

const (
	QUERY_EXCHANGE    = "queries"
	RESPONSE_EXCHANGE = "responses"
)

func main() {
	name := flag.String("name", "service1", "Service name")
	amqpuri := flag.String("amqpuri", "amqp://guest:guest@localhost/", "AMQP connection uri")
	flag.Parse()

	amqpPublisher := simpleamqp.NewAmqpPublisher(*amqpuri, "responses")
	amqpConsumer := simpleamqp.NewAmqpConsumer(*amqpuri)
	messages := amqpConsumer.Receive("queries", []string{"#"}, *name, 30*time.Minute)

	cont := 0
	for message := range messages {
		log.Println(message.Body)

		var m http2amqp.QueryMessage
		err := json.Unmarshal([]byte(message.Body), &m)
		if err == nil {
			response := http2amqp.ResponseMessage{m.Id, fmt.Sprintf("name <%s> cont %d", *name, cont)}

			json, _ := json.Marshal(response)
			amqpPublisher.Publish(m.Topic, []byte(json))
		}

		cont += cont + 1
	}

}
