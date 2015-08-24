// +build integration

package http2amqp_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	. "github.com/aleasoluciones/http2amqp"
	"github.com/stretchr/testify/assert"
)

var service = NewHTTP2AmqpService(os.Getenv("BROKER_URI"), "events", 1*time.Second)

func TestHttpSuccessfullGetToEchoServer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(NewHTTPServerFunc(service)))
	defer ts.Close()

	response, err := http.Get(ts.URL + "/test/ok")

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)
}

func TestHttpTimeoutGetToEchoServer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(NewHTTPServerFunc(service)))
	defer ts.Close()

	response, err := http.Get(ts.URL + "/test/timeout")

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 404)
}

func TestHttpSuccessfullPostToEchoServer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(NewHTTPServerFunc(service)))
	defer ts.Close()

	response, err := http.PostForm(ts.URL+"/test/ok", url.Values{"key": {"Value"}, "id": {"123"}})

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(body), "id=123&key=Value")
}
