package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aleasoluciones/simpleamqp"
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

type amqpRequestMessage struct {
	Id            string  `json:"id"`
	Request       Request `json:"request"`
	ResponseTopic string  `json:"responseTopic"`
}

type amqpResponseMessage struct {
	Id       string   `json:"id"`
	Response Response `json:"response"`
}

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
		var request amqpRequestMessage
		err := json.Unmarshal([]byte(inMessage.Body), &request)
		if err != nil {
			log.Panic("Error deserializing ", err)
		}

		fmt.Println("Receive", request)
		log.Println("Message Received. Id:", request.Id, request.Request.Method, request.Request.URL)
		log.Println("Body:", string(request.Request.Body))

		response := amqpResponseMessage{
			Id: request.Id,
			Response: Response{
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
