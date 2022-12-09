package server_test

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xbc5/sumo/internal/pkg/event"
	"github.com/xbc5/sumo/internal/pkg/server"
)

func wsUrl(u string) string {
	_url, err := url.ParseRequestURI(u)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse URL: %q", u))
	}
	return fmt.Sprintf("ws://%s/ws", _url.Host)
}

func close(conn *websocket.Conn) {
	closeMsg := websocket.FormatCloseMessage(1000, "finished")
	if err := conn.WriteMessage(websocket.CloseMessage, closeMsg); err != nil {
		panic(fmt.Sprintf("Cannot close connection: %s", err))
	}
}

func write(
	conn *websocket.Conn,
	results chan any,
	msgs int,
) {
	defer close(conn)
	for i := 0; i < msgs; i++ {
		conn.WriteJSON(input())
		var r interface{}
		conn.ReadJSON(&r)
		results <- r
	}
}

func echo() (*event.Msg, *httptest.Server, *websocket.Conn, *int) {
	s := serv().Build()
	called := 0

	go func(s *server.Server, called *int) {
		for msg := range s.Evt.NewReq {
			s.Evt.NewRes <- msg
			*called++
		}
	}(s, &called)

	server := s.StartTest()

	conn, _, _ := websocket.DefaultDialer.Dial(wsUrl(server.URL), nil)
	conn.WriteJSON(input())

	// should be an echo; also blocks (will hang test on error)
	result := event.Msg{}
	go conn.ReadJSON(&result)

	return &result, server, conn, &called
}

// const input = `{ "fake": "fake value" }`
type fakeInput struct {
	Fake string `json:"fake"`
}

func fakeMapResult() map[string]any {
	// use "any" because event.Msg.Data is any, which is the DTO that we use.
	// This makes it easier to match against. The data can literally be anything,
	// we do not test this at all. We only test that a generic message can be passed
	// to and from the server.
	return map[string]any{
		"fake": "fake value",
	}
}

func input() fakeInput {
	return fakeInput{
		Fake: "fake value",
	}
}

var _ = Describe("/ws route", func() {
	Context("when connecting", func() {
		It("should upgrade connection", func() {
			s := serv().Build().StartTest()
			defer s.Close()

			conn, res, err := websocket.DefaultDialer.Dial(wsUrl(s.URL), nil)
			defer close(conn)

			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(101)) // upgrading to WS
		})
	})

	Context("when sending a valid message", func() {
		Context("the client", func() {
			It("should not error", func() {
				s := serv().Build().StartTest()
				defer s.Close()

				conn, _, _ := websocket.DefaultDialer.Dial(wsUrl(s.URL), nil)
				defer close(conn)
				err := conn.WriteJSON(input())

				Expect(err).To(BeNil())
			})

			It("should receive an expected data value", func() {
				result, server, conn, _ := echo()
				defer server.Close()
				defer close(conn)

				Eventually(func() any {
					return result.Data
				}).WithTimeout(time.Second).Should(Equal(fakeMapResult()))
			})

			It("should receive a UUID value", func() {
				result, server, conn, _ := echo()
				defer server.Close()
				defer close(conn)

				uuidPat := `[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}`
				Eventually(func() any {
					return fmt.Sprintf("%s", result.UUID)
				}).WithTimeout(time.Second).Should(MatchRegexp(uuidPat))
			})
		})

		Context("the server", func() {
			It("should emit one response", func(ctx SpecContext) {
				_, server, conn, called := echo()
				defer server.Close()
				defer close(conn)

				// Eventually() will return as soon as called == 1; we want to make sure it doesn't go
				// beyond 1.
				time.Sleep(time.Millisecond * 500)

				Expect(*called).To(Equal(1))
			})
		})
	})

	Context("when sending an invalid message", func() {
		Context("the server", func() {
			It("should not close the socket (error)", func() {
				s := serv().Build().StartTest()
				defer s.Close()

				conn, _, _ := websocket.DefaultDialer.Dial(wsUrl(s.URL), nil)

				// do first message
				defer close(conn)
				conn.WriteMessage(websocket.TextMessage, []byte("bad message"))

				// try another message
				time.Sleep(time.Millisecond * 100)
				err := conn.WriteMessage(1, []byte("another bad message"))

				Expect(err).To(BeNil())
			})
		})
	})

	Context("with multiple concurrent connections", func() {
		Context("the server", func() {
			It("should handle them all", func() {
				threads := 1000 // concurrent connections
				msgs := 10      // messages per thread
				expected := threads * msgs

				s := serv().Build()
				httpserv := s.StartTest()
				defer httpserv.Close()

				go func(s *server.Server) {
					for msg := range s.Evt.NewReq {
						s.Evt.NewRes <- msg
					}
				}(s)

				rchan := make(chan interface{}, expected*2)

				// create T connections, where T = threads;
				conns := []*websocket.Conn{}
				for i := 0; i < threads; i++ {
					conn, _, _ := websocket.DefaultDialer.Dial(wsUrl(httpserv.URL), nil)
					conns = append(conns, conn)
				}

				// send M messages to T threads, where threads = T = M => TM
				// e.g. threads = 100; send 100 messages, for each of 100 threads.
				for _, conn := range conns {
					go write(conn, rchan, msgs)
				}

				Eventually(
					func() int { return len(rchan) },
				).WithTimeout(time.Second * 5).
					Should(Equal(expected))
			})
		})
	})
})
