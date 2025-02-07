package pubsub

import (
	"context"
	"sync"
)

type PubSub[T any] struct {
	mx       sync.RWMutex
	subs     []Subscriber[T]
	executor Executor
	subject  string
	closed   bool
}

func (that *PubSub[T]) Subject() string {
	return that.subject
}

func (that *PubSub[T]) Close(ctx context.Context) {
	that.mx.Lock()
	defer that.mx.Unlock()

	if !that.closed {
		that.closed = true
		for _, sub := range that.subs {
			sub.Close(ctx)
		}
	}
}

func (that *PubSub[T]) Subscribe(sub Subscriber[T]) {
	that.mx.Lock()
	defer that.mx.Unlock()

	that.subs = append(that.subs, sub)
}

func (that *PubSub[T]) Unsubscribe(ctx context.Context, subID string) {
	that.mx.Lock()
	defer that.mx.Unlock()

	for i, sub := range that.subs {
		if sub.ID() == subID {
			sub.Close(ctx)
			that.subs[i] = nil
			that.subs = append(that.subs[:i], that.subs[i+1:]...)
			break
		}
	}
}

func (that *PubSub[T]) Publish(ctx context.Context, entity T) Waiter {
	event := &Event[T]{
		ctx:     ctx,
		subject: that.subject,
		entity:  entity,
	}

	event.Enter()
	go that.publish(ctx, event)

	return event
}

func (that *PubSub[T]) publish(ctx context.Context, event *Event[T]) {
	defer event.Leave()

	that.mx.RLock()
	defer that.mx.RUnlock()

	if that.closed {
		return
	}

	for _, sub := range that.subs {
		that.pub(ctx, sub, event)
	}
}

func (that *PubSub[T]) pub(ctx context.Context, sub Subscriber[T], event *Event[T]) {
	event.Enter()

	that.executor.Execute(
		ctx,
		ActionFunc(
			func(ctx context.Context) error {
				defer event.Leave()

				sub.Handle(ctx, event)
				return nil
			},
		),
	)
}
