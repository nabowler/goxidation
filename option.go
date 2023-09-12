package goxidation

import "github.com/nabowler/goxidation/internal"

type (
	Option[T any] interface {
		IsSome() bool
		IsNone() bool
		// GetOr returns the contained Some value, or the provided def value if the result is an error result.
		GetOr(def T) T
		optionIF() internal.T
	}

	SomeOption[T any] interface {
		Get() T
		optionIF() internal.T
	}

	NoneOption[T any] interface {
		// noneIF causes `_, ok := some.(NoneOption[any])` to return false for the second return value
		noneIF() internal.T
	}

	some[T any] struct {
		val T
	}

	none[T any] struct{}
)

func Some[T any](val T) Option[T] {
	return some[T]{val: val}
}

func None[T any]() Option[T] {
	return none[T]{}
}

func (s some[T]) IsSome() bool {
	return true
}

func (s some[T]) IsNone() bool {
	return false
}

func (s some[T]) Get() T {
	return s.val
}

func (s some[T]) GetOr(def T) T {
	return s.val
}

//nolint:unused
func (s some[T]) optionIF() internal.T {
	t := new(internal.T)
	return *t
}

func (n none[T]) IsSome() bool {
	return false
}

func (n none[T]) IsNone() bool {
	return true
}

func (n none[T]) GetOr(def T) T {
	return def
}

//nolint:unused
func (n none[T]) optionIF() internal.T {
	t := new(internal.T)
	return *t
}

//nolint:unused
func (n none[T]) noneIF() internal.T {
	return n.optionIF()
}
