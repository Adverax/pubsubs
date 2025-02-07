package pubsub

import (
	"fmt"
)

type Builder[T any] struct {
	ps *PubSub[T]
}

func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{
		ps: &PubSub[T]{
			executor: DefaultExecutor,
		},
	}
}

func (that *Builder[T]) Subject(subject string) *Builder[T] {
	that.ps.subject = subject
	return that
}

func (that *Builder[T]) Executor(executor Executor) *Builder[T] {
	that.ps.executor = executor
	return that
}

func (that *Builder[T]) Build() (*PubSub[T], error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	return that.ps, nil
}

func (that *Builder[T]) checkRequiredFields() error {
	if that.ps.subject == "" {
		return ErrFieldSubjectIsRequired
	}

	return nil
}

var (
	ErrFieldSubjectIsRequired = fmt.Errorf("Field 'subject' is required")
)
