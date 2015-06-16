// +build integration

package http2amqp_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	. "github.com/aleasoluciones/http2amqp"
	"github.com/stretchr/testify/assert"
)

func TestHttpSuccessfullGetToEchoServer(t *testing.T) {
	t.Parallel()

	queriesService := NewQueriesServiceFactory(os.Getenv("BROKER_URI"), "events", 1*time.Second)

	ts := httptest.NewServer(http.HandlerFunc(NewHTTPServerFunc(queriesService)))
	defer ts.Close()

	response, err := http.Get(ts.URL + "/test/ok")

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)
}
