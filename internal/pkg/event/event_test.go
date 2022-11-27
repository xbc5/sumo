package event_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xbc5/sumo/internal/pkg/event"
)

func nMsgs(msgs ...string) (results []string) {
	evt := event.Evt[string]{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(len(msgs))
	evt.Sub("fakeEvt", func(_ context.Context, msg string) {
		defer mu.Unlock()
		mu.Lock()
		results = append(results, msg)
		wg.Done()
	})

	ctx := context.TODO()
	for _, msg := range msgs {
		evt.Pub(ctx, "fakeEvt", msg)
	}

	wg.Wait()

	return
}

func errMsg(id int) string {
	return fmt.Sprintf("Cannot remove subscription -- ID not found: %d", id)
}

var _ = Describe("Evt", func() {
	Context("publishing one event", func() {
		It("should yield one message", func() {
			Expect(nMsgs("one")).To(HaveLen(1))
		})

		It("should yield an expected message", func() {
			Expect(nMsgs("one")[0]).To(Equal("one"))
		})
	})

	Context("publishing two events", func() {
		It("should yield two messages", func() {
			Expect(nMsgs("one", "two")).To(HaveLen(2))
		})

		It("should yield the expected messages", func() {
			Expect(nMsgs("one", "two")).To(ConsistOf("one", "two"))
		})
	})

	Context("when subbing more than once", func() {
		It("should not crash when there's a publish occuring in parallel", func() {
			evt := event.Evt[string]{}
			called := 0

			// do continuous subs asynchronously
			go func(called *int) {
				for {
					go func(called *int) {
						evt.Sub("fakeEvt", func(_ context.Context, msg string) {
							*called++
						})
					}(called)
				}
			}(&called)

			time.Sleep(time.Millisecond * 100)

			ctx := context.TODO()
			evt.Pub(ctx, "fakeEvt", "msg")

			Consistently(called, 1).ShouldNot(Equal(0))
		})

		It("should provide unique IDs for each sub", func() {
			evt := event.Evt[string]{}
			ids := map[int]any{}
			max := 500

			for i := 1; i <= max; i++ {
				id := evt.Sub("fakeEvt", func(_ context.Context, msg string) {})
				ids[id] = nil
			}

			Expect(ids).To(HaveLen(max))
		})

		It("should safely handle when in parallel", func() {
			evt := event.Evt[string]{}
			called := 0
			max := 50000 // racy is hard to detect, go full nuke
			var mu sync.Mutex

			// sub lots in parallel
			var wg sync.WaitGroup
			wg.Add(max)
			for i := 1; i <= max; i++ {
				go func(wg *sync.WaitGroup, called *int, mu *sync.Mutex) {
					evt.Sub("fakeEvt", func(_ context.Context, msg string) {
						defer mu.Unlock()
						mu.Lock()
						*called++
					})
					wg.Done()
				}(&wg, &called, &mu)
			}
			wg.Wait()

			ctx := context.TODO()
			evt.Pub(ctx, "fakeEvt", "msg")

			// wait until all calls made
			for i := 0; i <= 10; i++ {
				if called == max {
					break
				}
				time.Sleep(time.Millisecond * 100)
			}

			// just check that all calls made; concurrent sub writes will panic on its own
			Expect(called).To(Equal(max))
		})
	})

	Context("when unsubbing", func() {
		Context("and the ID doesn't exist", func() {
			It("should return an error", func() {
				evt := event.Evt[string]{}

				Expect(evt.Unsub(1)).Should(MatchError(errMsg(1)))
			})
		})

		Context("after removing a sub", func() {
			It("should not result in a call", func(ctx SpecContext) {
				called := 0
				evt := event.Evt[string]{}

				id := evt.Sub("fakeEvt", func(_ context.Context, msg string) {
					called++
				})

				evt.Unsub(id)
				c := context.TODO()
				evt.Pub(c, "fakeEvt", "msg")

				for i := 0; i < 10; i++ {
					if called > 0 {
						break
					}
					time.Sleep(time.Millisecond * 10)
				}

				Expect(called).To(Equal(0))
			})

			It("should return an error on subsequent unsub", func() {
				evt := event.Evt[string]{}
				id := evt.Sub("fakeEvt", func(_ context.Context, msg string) {})

				Expect(evt.Unsub(id)).To(BeNil())

				Expect(evt.Unsub(1)).Should(MatchError(errMsg(1)))
			})
		})

		Context("in parallel", func() {
			It("should unsub safely", func() {
				// we're really only checking that concurrent unsubs don't Panic
				called := 0
				evt := event.Evt[string]{}
				max := 500           // the number of subs
				ids := map[int]any{} // so we know when unsubs are done
				var mu sync.Mutex

				// sub lots
				for i := 1; i <= max; i++ {
					id := evt.Sub("fakeEvt", func(_ context.Context, msg string) {
						defer mu.Unlock()
						mu.Lock()
						called++
					})
					mu.Lock()
					ids[id] = nil // track subs
					mu.Unlock()
				}

				// unsub all
				for id := range ids {
					go evt.Unsub(id) // concurrent
					mu.Lock()
					delete(ids, id) // indicate unsub
					mu.Unlock()
				}

				// poll for unsubs
				for i := 0; i < 5; i++ {
					if len(ids) == 0 {
						break
					}
					time.Sleep(time.Millisecond * 100)
				}

				Expect(ids).To(HaveLen(0)) // just check that ubsubs are done
			})
		})

		It("should unsub many", func(ctx SpecContext) {
			called := 0
			evt := event.Evt[string]{}
			max := 500 // the number of subs
			var ids []int
			var mu sync.Mutex

			// sub lots
			for i := 1; i <= max; i++ {
				id := evt.Sub("fakeEvt", func(_ context.Context, msg string) {
					defer mu.Unlock()
					mu.Lock()
					called++
				})
				mu.Lock()
				ids = append(ids, id)
				mu.Unlock()
			}

			// unsub all
			c := context.TODO()
			for _, id := range ids {
				evt.Unsub(id)
			}

			evt.Pub(c, "fakeEvt", "msg")

			// poll
			for i := 0; i < 5; i++ {
				if called > 0 {
					break
				}
				time.Sleep(time.Millisecond * 100)
			}

			Expect(called).To(Equal(0))
		})

		It("should unsub all but keep some", func(ctx SpecContext) {
			called := 0
			evt := event.Evt[string]{}
			max := 500    // the number of subs
			expected := 2 // the subs that we're left with
			ids := map[int]any{}
			var mu sync.Mutex

			// sub lots
			for i := 1; i <= max; i++ {
				id := evt.Sub("fakeEvt", func(_ context.Context, msg string) {
					defer mu.Unlock()
					mu.Lock()
					called++
				})
				mu.Lock()
				ids[id] = nil
				mu.Unlock()
			}

			// unsub most
			c := context.TODO()
			count := 1 // we break after almost all
			for id := range ids {
				evt.Unsub(id)
				delete(ids, id)
				if count == max-expected {
					break
				}
				count++
			}

			evt.Pub(c, "fakeEvt", "msg")

			// poll
			for i := 0; i < 5; i++ {
				if called == expected {
					break
				}
				time.Sleep(time.Millisecond * 100)
			}

			Expect(called).To(Equal(expected))
		})
	})
})
