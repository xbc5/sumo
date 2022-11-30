package server_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xbc5/sumo/internal/pkg/server"
)

func serv() server.ServerBuilder {
	s := server.Server{}
	return s.NewTest()
}

var _ = Describe("Common", func() {
	Context("after starting the server", func() {
		It("should not error", func() {
			s := serv().Build().StartTest()
			c := s.Client()
			defer s.Close()
			_, err := c.Get(s.URL)
			Expect(err).To(BeNil())
		})
	})
})
