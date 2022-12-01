package event

import (
	"context"
	"fmt"
	"sync"
)

type IEvt[T any] interface {
	Sub(event string, fn func(ctx context.Context, msg T)) int
	Unsub(id int) error
	Pub(ctx context.Context, event string, msg T) *Evt[T]
}

type Evt[T any] struct {
	subs  map[string]map[int]func(ctx context.Context, msg T)
	cache map[int]string
	mu    sync.RWMutex
	id    int
}

func (this *Evt[T]) tryNew() {
	if this.subs == nil {
		this.subs = map[string]map[int]func(ctx context.Context, msg T){}
	}
	if this.cache == nil {
		this.cache = map[int]string{}
	}
}

func (this *Evt[T]) Sub(event string, fn func(ctx context.Context, msg T)) int {
	defer this.mu.Unlock()
	this.mu.Lock()

	this.tryNew()

	if this.subs[event] == nil {
		this.subs[event] = map[int]func(ctx context.Context, msg T){}
	}

	this.id++
	this.subs[event][this.id] = fn
	this.cache[this.id] = event // for each unsubs

	return this.id
}

func (this *Evt[_]) Unsub(id int) error {
	defer this.mu.Unlock()
	this.mu.Lock()

	event := this.cache[id]
	if event == "" {
		return fmt.Errorf("Cannot remove subscription -- ID not found: %d", id)
	}

	delete(this.subs[event], id)
	delete(this.cache, id)

	return nil
}

func (this *Evt[T]) Pub(ctx context.Context, event string, msg T) *Evt[T] {
	// since pub will be the most frequently called behaviour, concurrent reads
	// should alleviate any performance impact with mutexes.
	defer this.mu.RUnlock()
	this.mu.RLock()

	this.tryNew()

	for _, pub := range this.subs[event] {
		go pub(ctx, msg)
	}

	return this
}
