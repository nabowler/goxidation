package goxidation

import "github.com/nabowler/goxidation/internal"

type (
	Result[T any] interface {
		IsOk() bool
		IsError() bool
		// GetOr returns the contained Ok value, or the provided def value if the result is an error result.
		GetOr(def T) T
		resultIF() internal.T
	}

	OkResult[T any] interface {
		Get() T
		resultIF() internal.T
	}

	ErrResult[T any] interface {
		Unwrap() error
		Error() string
		resultIF() internal.T
	}

	ok[T any] struct {
		val T
	}

	err[T any] struct {
		inner error
	}
)

const (
	DefaultErrorResultMessage = "Error Result"
)

func Ok[T any](val T) Result[T] {
	return ok[T]{val: val}
}

func Error[T any](errr error) Result[T] {
	return err[T]{inner: errr}
}

func (o ok[T]) IsOk() bool {
	return true
}

func (o ok[T]) IsError() bool {
	return false
}

func (o ok[T]) Get() T {
	return o.val
}

func (o ok[T]) GetOr(def T) T {
	return o.val
}

//nolint:unused
func (o ok[T]) resultIF() internal.T {
	t := new(internal.T)
	return *t
}

func (e err[T]) IsOk() bool {
	return false
}

func (e err[T]) IsError() bool {
	return true
}

func (e err[T]) GetOr(def T) T {
	return def
}

func (e err[T]) Unwrap() error {
	return e.inner
}

func (e err[T]) Error() string {
	if e.inner == nil {
		return DefaultErrorResultMessage
	}
	return e.inner.Error()
}

//nolint:unused
func (e err[T]) resultIF() internal.T {
	t := new(internal.T)
	return *t
}
