package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aleasoluciones/simpleamqp"
)

func main() {
	name := flag.String("name", "service1", "Service name")
	amqpuri := flag.String("amqpuri", "amqp://guest:guest@localhost/", "AMQP connection uri")
	flag.Parse()

	amqpPublisher := simpleamqp.NewAmqpPublisher(*amqpuri, "events")
	amqpConsumer := simpleamqp.NewAmqpConsumer(*amqpuri)
	messages := amqpConsumer.Receive("events", []string{"#"}, *name, 30*time.Second)

	cont := 0
	for message := range messages {
		log.Println(message)

		messageBody := fmt.Sprintf("name <%s> cont %d", *name, cont)

		amqpPublisher.Publish("efa2", []byte(messageBody))

		cont += cont + 1
	}

}
