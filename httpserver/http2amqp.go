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
	HTTPAddress := flag.String("address", defaultHTTPAddress(), "Listen address")
	HTTPPort := flag.String("port", defaultHTTPPort(), "Listen port")
	exchange := flag.String("exchange", "events", "AMQP exchange name")
	timeout := flag.Int("timeout", 1000, "Queries timeout in milliseconds")
	flag.Parse()

	if *verbose {
		alealog.EnableLogging()
		log.Println("[http2amqp] verbose mode enabled")
	}

	service := http2amqp.NewService(*amqpuri, *exchange, time.Duration(*timeout)*time.Millisecond)

	http.HandleFunc("/", http2amqp.NewHTTPServerFunc(service))
	addressAndPort := fmt.Sprintf("%s:%s", *HTTPAddress, *HTTPPort)
	log.Println("[http2amqp] Starting HTTP server at ", addressAndPort)
	http.ListenAndServe(addressAndPort, nil)
}

func defaultHTTPPort() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "18080"
	}

	return port
}

func defaultHTTPAddress() string {
	address := os.Getenv("ADDRESS")

	if len(address) == 0 {
		address = "0.0.0.0"
	}

	return address
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
