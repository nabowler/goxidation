package goxidation_test

import (
	"fmt"
	"testing"

	"github.com/nabowler/goxidation"
	"github.com/stretchr/testify/assert"
)

const (
	getOrValue = "some alternate value"
)

func anyTestCaseValues() []any {
	return []any{
		"",
		nil,
		1,
		3.14,
		true,
		false,
		(*int)(nil),
		(*string)(nil),
		[]any{"1", 2, true},
	}
}

func TestOk(t *testing.T) {
	for _, tc := range anyTestCaseValues() {
		tc := tc
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			some := goxidation.Some(tc)
			assert.True(t, some.IsSome())
			assert.False(t, some.IsNone())

			val := some.GetOr(getOrValue)
			assert.Equal(t, tc, val)

			someOpt := some.(goxidation.SomeOption[any])
			val = someOpt.Get()
			assert.Equal(t, tc, val)
		})
	}
}

func TestNone(t *testing.T) {
	none := goxidation.None[any]()
	assert.False(t, none.IsSome())
	assert.True(t, none.IsNone())

	val := none.GetOr(getOrValue)
	assert.Equal(t, getOrValue, val)

	_, ok := none.(goxidation.NoneOption[any])
	assert.True(t, ok)
}

func TestNoneCoercion(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "expected panic")
	}()
	none := goxidation.None[any]()
	// "unsafe" assertion variant should panic
	_ = none.(goxidation.SomeOption[any])
	t.Errorf("Expected panic by now")
}

func TestNoneCoercionWithSafeCast(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "expected panic")
	}()
	none := goxidation.None[any]()
	// use the "safe" assertion varianbt to try to get an SomeOption value
	someRes, ok := none.(goxidation.SomeOption[any])
	assert.False(t, ok)
	assert.Nil(t, someRes)
	// trying to use the someRes will panic on a nil pointer
	_ = someRes.Get()
	t.Errorf("Expected panic by now")
}

func TestSomeCoercion(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "expected panic")
	}()
	some := goxidation.Some[any]("val")
	// "unsafe" assertion variant should panic
	_ = some.(goxidation.NoneOption[any])
	t.Errorf("Expected panic by now")
}

func TestSomeCoercionWithSafeCast(t *testing.T) {
	some := goxidation.Some[any]("val")
	// use the "safe" assertion varianbt to try to get an SomeOption value
	noneRes, ok := some.(goxidation.NoneOption[any])
	assert.False(t, ok)
	assert.Nil(t, noneRes)

	// NoneRes provides no usable methods so there's not really a way to cause an NPE
}
