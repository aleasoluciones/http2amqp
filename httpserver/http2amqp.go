package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	alealog "github.com/aleasoluciones/goaleasoluciones/log"
	"github.com/aleasoluciones/http2amqp"
)

func main() {
	alealog.Init()
	alealog.DisableLogging()
	defaultVerbose := verboseMode()

	verbose := flag.Bool("verbose", defaultVerbose, "Verbose mode, enable logging")
	amqpuri := flag.String("amqpuri", localBrokerUri(), "AMQP connection uri")
	address := flag.String("address", "0.0.0.0", "Listen address")
	port := flag.String("port", http2amqpPort(), "Listen port")
	exchange := flag.String("exchange", "events", "AMQP exchange name")
	timeout := flag.Int("timeout", 1000, "Queries timeout in milliseconds")
	flag.Parse()

	if *verbose {
		alealog.EnableLogging()
		log.Println("[http2amqp] verbose mode enabled")
	}

	service := http2amqp.NewService(*amqpuri, *exchange, time.Duration(*timeout)*time.Millisecond)

	http.HandleFunc("/", http2amqp.NewHTTPServerFunc(service))
	addressAndPort := fmt.Sprintf("%s:%s", *address, *port)
	log.Println("[http2amqp] Starting HTTP server at ", addressAndPort)
	http.ListenAndServe(addressAndPort, nil)
}

func http2amqpPort() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "18080"
	}

	return port
}

func localBrokerUri() string {
	brokerUri := os.Getenv("BROKER_URI")

	if len(brokerUri) == 0 {
		brokerUri = "amqp://guest:guest@localhost/"
	}

	return brokerUri
}

func isTrue(value string) bool {
	value = strings.ToLower(value)
	trulies := []string{"1", "y", "yes", "on"}

	for _, truly := range trulies {
		if value == truly {
			return true
		}
	}
	return false
}

func verboseMode() bool {
	verbose := os.Getenv("HTTP2AMQP_VERBOSE")
	return isTrue(verbose)
}
