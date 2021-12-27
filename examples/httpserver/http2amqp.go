package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"aleasoluciones/http2amqp"

	alealog "github.com/aleasoluciones/goaleasoluciones/log"
)

const defaultTimeoutMs = 1000
const defaultExchange = "events"
const defaultHTTPPort = "18080"
const defaultHTTPAddress = "0.0.0.0"
const defaultBrokerURI = "amqp://guest:guest@localhost/"

func main() {
	alealog.Init()
	alealog.DisableLogging()
	defaultVerbose := verboseMode()

	verbose := flag.Bool("verbose", defaultVerbose, "Verbose mode, enable logging")
	brokeruri := flag.String("brokeruri", envBrokerURI(), "AMQP broker connection URI")
	HTTPAddress := flag.String("address", envHTTPAddress(), "HTTP listen IP address")
	HTTPPort := flag.String("port", envHTTPPort(), "HTTP listen port")
	exchange := flag.String("exchange", defaultExchange, "AMQP broker exchange name")
	timeout := flag.Int("timeout", defaultTimeoutMs, "AMQP broker queries timeout in milliseconds")
	flag.Parse()

	if *verbose {
		alealog.EnableLogging()
		log.Println("[http2amqp] verbose mode enabled")
	}

	service := http2amqp.NewService(*brokeruri, *exchange, time.Duration(*timeout)*time.Millisecond)

	http.HandleFunc("/", http2amqp.NewHTTPServerFunc(service))
	addressAndPort := fmt.Sprintf("%s:%s", *HTTPAddress, *HTTPPort)
	log.Println("[http2amqp] Starting HTTP server at ", addressAndPort)
	http.ListenAndServe(addressAndPort, nil)
}

func envHTTPPort() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = defaultHTTPPort
	}

	return port
}

func envHTTPAddress() string {
	address := os.Getenv("ADDRESS")

	if len(address) == 0 {
		address = defaultHTTPAddress
	}

	return address
}

func envBrokerURI() string {
	brokerURI := os.Getenv("BROKER_URI")

	if len(brokerURI) == 0 {
		brokerURI = defaultBrokerURI
	}

	return brokerURI
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
