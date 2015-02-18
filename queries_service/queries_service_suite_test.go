package queries_service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestQueriesService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "QueriesService Suite")
}
