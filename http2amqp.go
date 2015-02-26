package main

import (
	"flag"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/aleasoluciones/http2amqp/queries_http"
	"github.com/aleasoluciones/http2amqp/queries_service"
)

func main() {
	brokerUri, exchange, timeout := parseArgs()

	queriesService := queries_service.NewQueriesServiceFactory(brokerUri, exchange, timeout)
	queries_http.NewHTTPServer(queriesService)
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
