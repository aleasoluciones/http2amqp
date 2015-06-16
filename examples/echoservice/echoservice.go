package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aleasoluciones/http2amqp"
	"github.com/aleasoluciones/simpleamqp"
)

func main() {
	amqpuri := flag.String("amqpuri", LocalBrokerUri(), "AMQP connection uri")
	flag.Parse()

	fmt.Println("Using broker:", *amqpuri)
	amqpPublisher := simpleamqp.NewAmqpPublisher(*amqpuri, "EX")
	amqpConsumer := simpleamqp.NewAmqpConsumer(*amqpuri)

	inMessages := amqpConsumer.Receive(
		"EX",
		[]string{"GET.#"},
		"",
		simpleamqp.QueueOptions{Durable: false, Delete: true, Exclusive: true},
		5*60*time.Second)

	for inMessage := range inMessages {
		log.Println("Message Received. Topic:", inMessage.RoutingKey)
		var request http2amqp.AmqpRequestMessage
		err := json.Unmarshal([]byte(inMessage.Body), &request)
		if err != nil {
			log.Panic("Error deserializing ", err)
		}

		fmt.Println("Receive", request)
		log.Println("Message Received. Id:", request.Id, request.Request.Method, request.Request.URL)
		log.Println("Body:", string(request.Request.Body))

		response := http2amqp.AmqpResponseMessage{
			Id: request.Id,
			Response: http2amqp.Response{
				Body:   request.Request.Body,
				Status: 200,
			},
		}

		serializedResponse, _ := json.Marshal(response)
		fmt.Println("Sending response", string(serializedResponse))
		amqpPublisher.Publish(request.ResponseTopic, serializedResponse)
	}
}
func LocalBrokerUri() string {
	brokerUri := os.Getenv("BROKER_URI")
	if len(brokerUri) == 0 {
		brokerUri = "amqp://guest:guest@localhost/"
	}

	return brokerUri
}
