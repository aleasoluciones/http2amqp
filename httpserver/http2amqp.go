package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/aleasoluciones/http2amqp"
)

func main() {
	amqpuri := flag.String("amqpuri", localBrokerUri(), "AMQP connection uri")
	address := flag.String("address", "0.0.0.0", "Listen address")
	port := flag.String("port", "18080", "Listen port")
	exchange := flag.String("exchange", "events", "AMQP exchange name")
	timeout := flag.Int("timeout", 1000, "Queries timeout in milliseconds")
	flag.Parse()

	http2amqpService := http2amqp.NewHttp2AmqpService(*amqpuri, *exchange, time.Duration(*timeout)*time.Millisecond)

	http.HandleFunc("/", http2amqp.NewHTTPServerFunc(http2amqpService))
	addressAndPort := fmt.Sprintf("%s:%s", *address, *port)
	log.Println("[http2amqp] Starting HTTP server at ", addressAndPort)
	http.ListenAndServe(addressAndPort, nil)
}

func localBrokerUri() string {
	brokerUri := os.Getenv("BROKER_URI")

	if len(brokerUri) == 0 {
		brokerUri = "amqp://guest:guest@localhost/"
	}

	return brokerUri
}
