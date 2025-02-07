package pubsub

import (
	"context"
	"sync"
)

type Action interface {
	Execute(ctx context.Context) error
}

type Executor interface {
	Execute(ctx context.Context, action Action)
}

type ActionFunc func(ctx context.Context) error

func (fn ActionFunc) Execute(ctx context.Context) error {
	return fn(ctx)
}

type Waiter interface {
	Wait()
	WaitWithContext(ctx context.Context) error
}

type Handler[T any] interface {
	Handle(ctx context.Context, event *Event[T])
}

type HandlerFunc[T any] func(ctx context.Context, event *Event[T])

func (fn HandlerFunc[T]) Handle(ctx context.Context, event *Event[T]) {
	fn(ctx, event)
}

type Subscriber[T any] interface {
	Handler[T]
	ID() string
	Close(ctx context.Context)
}

type Event[T any] struct {
	ctx     context.Context
	wg      sync.WaitGroup
	subject string
	entity  T
}

func (that *Event[T]) Enter() {
	that.wg.Add(1)
}

func (that *Event[T]) Leave() {
	that.wg.Done()
}

func (that *Event[T]) Context() context.Context {
	return that.ctx
}

func (that *Event[T]) Subject() string {
	return that.subject
}

func (that *Event[T]) Entity() T {
	return that.entity
}

func (that *Event[T]) Wait() {
	that.wg.Wait()
}

func (that *Event[T]) WaitWithContext(ctx context.Context) error {
	done := make(chan struct{}, 1)

	go func() {
		that.wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type defaultExecutor struct {
}

func (that *defaultExecutor) Execute(
	ctx context.Context,
	action Action,
) {
	_ = action.Execute(ctx)
}

var defExecutor = &defaultExecutor{}
