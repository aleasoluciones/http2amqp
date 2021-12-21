package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"aleasoluciones/http2amqp"

	"github.com/aleasoluciones/simpleamqp"
)

func main() {
	amqpuri := flag.String("amqpuri", LocalBrokerUri(), "AMQP connection uri")
	exchange := flag.String("exchange", "events", "AMQP exchange name")
	topic := flag.String("topic", "get.#", "topic to subscribe")
	flag.Parse()

	fmt.Println("Using broker:", *amqpuri)
	amqpPublisher := simpleamqp.NewAmqpPublisher(*amqpuri, *exchange)
	amqpConsumer := simpleamqp.NewAmqpConsumer(*amqpuri)

	inMessages := amqpConsumer.ReceiveWithoutTimeout(
		*exchange,
		[]string{*topic},
		"",
		simpleamqp.QueueOptions{Durable: false, Delete: true, Exclusive: true})

	for inMessage := range inMessages {
		log.Println("Message Received. Topic:", inMessage.RoutingKey)
		var request http2amqp.AmqpRequestMessage
		err := json.Unmarshal([]byte(inMessage.Body), &request)
		if err != nil {
			log.Panic("Error deserializing ", err)
		}

		fmt.Println("Receive", request)
		log.Println("Message Received. Id:", request.ID, request.Request.Method, request.Request.URL)
		log.Println("Body:", string(request.Request.Body))

		response := http2amqp.AmqpResponseMessage{
			ID: request.ID,
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
