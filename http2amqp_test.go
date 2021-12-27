//go:build integration

package http2amqp_test

import (
	. "aleasoluciones/http2amqp"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func envBrokerURI() string {

	brokerURI := os.Getenv("BROKER_URI")

	if len(brokerURI) == 0 {
		brokerURI = "amqp://guest:guest@localhost/"
	}

	return brokerURI
}

var service = NewService(envBrokerURI(), "events", 1*time.Second)

func TestHttpSuccessfullGetToEchoServer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(NewHTTPServerFunc(service)))
	defer ts.Close()

	// Sleep to allow the consumer.Receive to connect to rabbitmq
	time.Sleep(5 * time.Second)

	response, err := http.Get(ts.URL + "/test/ok")

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)
}

func TestHttpTimeoutGetToEchoServer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(NewHTTPServerFunc(service)))
	defer ts.Close()

	// Sleep to allow the consumer.Receive to connect to rabbitmq
	time.Sleep(5 * time.Second)
	response, err := http.Get(ts.URL + "/test/timeout")

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 408)
}

func TestHttpSuccessfullPostToEchoServer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(NewHTTPServerFunc(service)))
	defer ts.Close()

	// Sleep to allow the consumer.Receive to connect to rabbitmq
	time.Sleep(5 * time.Second)
	response, err := http.PostForm(ts.URL+"/test/ok", url.Values{"key": {"Value"}, "id": {"123"}})

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(body), "id=123&key=Value")
}
