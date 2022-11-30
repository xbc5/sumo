package server_test

import (
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("/", func() {
	Context("when connecting to /", func() {
		It("should not error", func() {
			s := serv().Build().StartTest()
			c := s.Client()
			defer s.Close()

			res, err := c.Get(s.URL + "/")
			if res != nil {
				defer res.Body.Close()
			}

			Expect(err).To(BeNil())
		})

		It("should return 200 OK", func() {
			s := serv().Build().StartTest()
			c := s.Client()
			defer s.Close()

			res, _ := c.Get(s.URL + "/")
			if res != nil {
				defer res.Body.Close()
			}

			Expect(res.StatusCode).To(Equal(200))
		})

		It("should return an expected result", func() {
			// FIXME: mae test better:
			// this just tests that a stub handler returns something; yet the real handler is untested
			s := serv().Build().StartTest()
			c := s.Client()
			defer s.Close()

			res, _ := c.Get(s.URL + "/")
			if res != nil {
				defer res.Body.Close()
			}
			body, err := io.ReadAll(res.Body)

			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("response stub"))
		})
	})
})
