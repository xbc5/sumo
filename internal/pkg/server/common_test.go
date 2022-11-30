package server_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xbc5/sumo/internal/pkg/server"
)

func serv() server.ServerBuilder {
	s := server.Server{}
	return s.NewTest()
}

var _ = Describe("Common", func() {
	Context("when starting the server", func() {
		It("should start", func() {
			s := serv().Build().StartTest()
			defer s.Close()
			req, _ := http.NewRequest("GET", s.URL, nil)
			res, resErr := http.DefaultClient.Do(req)
			Expect(resErr).To(BeNil())
			defer res.Body.Close()
		})
	})
})
