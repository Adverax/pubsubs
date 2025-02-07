package pubsub

import (
	"context"
	"fmt"
	"sync/atomic"
)

type Subscription[T any] struct {
	id string
	Handler[T]
}

func NewSubscription[T any](handler Handler[T]) *Subscription[T] {
	return &Subscription[T]{
		id:      newIdentifier(),
		Handler: handler,
	}
}

func (that *Subscription[T]) ID() string {
	return that.id
}

func (that *Subscription[T]) Close(ctx context.Context) {
	// do nothing
}

type ChannelSubscription[T any] struct {
	id string
	ch chan *Event[T]
}

func NewChannelSubscription[T any](cap int) *ChannelSubscription[T] {
	return &ChannelSubscription[T]{
		id: newIdentifier(),
		ch: make(chan *Event[T], cap),
	}
}

func (that *ChannelSubscription[T]) ID() string {
	return that.id
}

func (that *ChannelSubscription[T]) Close(ctx context.Context) {
	close(that.ch)
}

func (that *ChannelSubscription[T]) Handle(ctx context.Context, event *Event[T]) {
	that.ch <- event
}

func (that *ChannelSubscription[T]) Channel() <-chan *Event[T] {
	return that.ch
}

func (that *ChannelSubscription[T]) Serve(ctx context.Context, handler Handler[T]) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-that.ch:
			handler.Handle(event.ctx, event)
		}
	}
}

var newIdentifier = func() func() string {
	var counter int32
	return func() string {
		val := atomic.AddInt32(&counter, 1)
		return fmt.Sprintf("%d", val)
	}
}()
