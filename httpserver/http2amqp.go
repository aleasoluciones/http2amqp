package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/aleasoluciones/http2amqp"
)

func main() {
	brokerUri, exchange, timeout := parseArgs()

	http2amqpService := http2amqp.NewHttp2AmqpService(brokerUri, exchange, timeout)

	http.HandleFunc("/", http2amqp.NewHTTPServerFunc(http2amqpService))
	log.Println("[http2amqp] Starting HTTP server at 127.0.0.1:18080 ...")
	http.ListenAndServe("127.0.0.1:18080", nil)

}

func parseArgs() (string, string, time.Duration) {
	amqpuri := flag.String("amqpuri", localBrokerUri(), "AMQP connection uri")
	exchange := flag.String("exchange", "events", "AMQP exchange name")
	timeout := flag.Int("timeout", 1, "Queries timeout in seconds")
	flag.Parse()

	return *amqpuri, *exchange, time.Duration(*timeout) * time.Second
}

func localBrokerUri() string {
	brokerUri := os.Getenv("BROKER_URI")

	if len(brokerUri) == 0 {
		brokerUri = "amqp://guest:guest@localhost/"
	}

	return brokerUri
}
