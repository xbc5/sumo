package server_test

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func wsUrl(u string) string {
	_url, err := url.ParseRequestURI(u)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse URL: %q", u))
	}
	return fmt.Sprintf("ws://%s/ws", _url.Host)
}

const json string = `
{
  "test": "test value"
}
`

var _ = Describe("/ws route", func() {
	Context("when connecting", func() {
		It("should upgrade connection", func() {
			s := serv().Build().StartTest()
			defer s.Close()

			fmt.Printf(s.URL)
			conn, res, err := websocket.DefaultDialer.Dial(wsUrl(s.URL), nil)
			defer conn.Close()

			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(101)) // upgrading to WS
		})
	})

	Context("when sending a message", func() {
		It("should not error", func() {
			s := serv().Build().StartTest()
			defer s.Close()

			fmt.Printf(s.URL)
			conn, _, _ := websocket.DefaultDialer.Dial(wsUrl(s.URL), nil)
			defer conn.Close()
			err := conn.WriteMessage(websocket.TextMessage, []byte(json))

			Expect(err).To(BeNil())
		})
	})
})
