package server_test

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

const input = `{ "fake": "fake value" }`

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
				err := conn.WriteJSON(input)

				Expect(err).To(BeNil())
			})

			It("should receive an expected value", func() {
				s := serv().Build().StartTest()
				defer s.Close()
				var result interface{}

				conn, _, _ := websocket.DefaultDialer.Dial(wsUrl(s.URL), nil)
				defer close(conn)
				conn.WriteJSON(input)

				// should be an echo; also blocks (will hang test on error)
				go conn.ReadJSON(&result)

				Eventually(func() interface{} {
					return result
				}).WithTimeout(time.Second).Should(Equal(input))
			})
		})

		Context("the server", func() {
			It("should emit one response", func(ctx SpecContext) {
				called := 0
				evt := server.NewOkEvtStub()
				evt.Sub(server.ResReady, func(_ context.Context, _ any) {
					called++
				})
				s := serv().Evt(evt).Build().StartTest()
				defer s.Close()

				conn, _, _ := websocket.DefaultDialer.Dial(wsUrl(s.URL), nil)
				defer close(conn)
				conn.WriteJSON(input)
				Eventually(func() int { return called }).WithTimeout(time.Second).Should(Equal(1))
			})
		})
	})

	Context("when sending an invalid message", func() {
		Context("the server", func() {
			It("should recover and allow subsequent connections", func() {
				s := serv().Build().StartTest()
				defer s.Close()

				conn, _, _ := websocket.DefaultDialer.Dial(wsUrl(s.URL), nil)
				defer close(conn)
				conn.WriteMessage(websocket.TextMessage, []byte("bad message"))

				time.Sleep(time.Millisecond * 100)
				_, _, err := websocket.DefaultDialer.Dial(wsUrl(s.URL), nil)

				Expect(err).To(BeNil()) // errors if StatusCode != 101 (amongst other things)
			})
		})
	})

	Context("with multiple concurrent connections", func() {
		Context("the server", func() {
			It("should handle them all", func() {
				s := serv().Build().StartTest()
				defer s.Close()
				max := 100
				results := []interface{}{}
				var mu sync.Mutex

				conns := []*websocket.Conn{}
				for i := 0; i < max; i++ {
					conn, _, _ := websocket.DefaultDialer.Dial(wsUrl(s.URL), nil)
					conns = append(conns, conn)
				}

				var wg sync.WaitGroup
				wg.Add(max)
				for _, conn := range conns {
					go func(wg *sync.WaitGroup, conn *websocket.Conn, results *[]interface{}) {
						defer close(conn)
						defer wg.Done()
						for i := 0; i < max; i++ {
							conn.WriteJSON(input)
							var r interface{}
							conn.ReadJSON(&r)
							mu.Lock()
							*results = append(*results, r)
							mu.Unlock()
						}
					}(&wg, conn, &results)
				}
				wg.Wait()

				Eventually(func() int { return len(results) }).Should(Equal(max * max))
			})
		})
	})
})
