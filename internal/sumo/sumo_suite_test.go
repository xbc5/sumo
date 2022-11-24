package sumo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSumo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sumo Suite")
}
