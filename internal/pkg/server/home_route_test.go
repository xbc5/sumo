package server_test

import (
	"io"

	"github.com/bradleyjkemp/cupaloy"
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

		It("should return text/html Content-Type", func() {
			s := serv().Build().StartTest()
			c := s.Client()
			defer s.Close()

			res, _ := c.Get(s.URL + "/")
			if res != nil {
				defer res.Body.Close()
			}

			Expect(res.Header.Get("Content-Type")).To(Equal("text/html; charset=utf-8"))
		})

		It("should return an expected result", func() {
			s := serv().Build().StartTest()
			c := s.Client()
			defer s.Close()

			res, _ := c.Get(s.URL + "/")
			if res != nil {
				defer res.Body.Close()
			}
			body, err := io.ReadAll(res.Body)

			Expect(err).To(BeNil())
			cupaloy.SnapshotT(GinkgoT(), body)
		})
	})
})
