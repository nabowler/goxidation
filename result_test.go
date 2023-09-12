package goxidation_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/nabowler/goxidation"
	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	for _, tc := range anyTestCaseValues() {
		tc := tc
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			ok := goxidation.Ok(tc)
			assert.True(t, ok.IsOk())
			assert.False(t, ok.IsError())

			val := ok.GetOr(getOrValue)
			assert.Equal(t, tc, val)

			okRes := ok.(goxidation.OkResult[any])
			val = okRes.Get()
			assert.Equal(t, tc, val)
		})
	}
}

func TestError(t *testing.T) {
	for _, tc := range []error{
		nil,
		io.EOF,
		fmt.Errorf("This is a test of the emergency broadcast system."),
	} {
		tc := tc
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			err := goxidation.Error[any](tc)
			assert.False(t, err.IsOk())
			assert.True(t, err.IsError())

			val := err.GetOr(getOrValue)
			assert.Equal(t, getOrValue, val)

			errRes := err.(goxidation.ErrResult[any])
			innerErr := errRes.Unwrap()
			assert.Equal(t, tc, innerErr)
			assert.Equal(t, tc != nil, errors.Is(errRes, tc))

			expectedErrorMsg := goxidation.DefaultErrorResultMessage
			if tc != nil {
				expectedErrorMsg = tc.Error()
			}
			assert.Equal(t, expectedErrorMsg, errRes.Error())
		})
	}
}

func TestErrorCoercion(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "expected panic")
	}()
	err := goxidation.Error[any](fmt.Errorf("some error"))
	// "unsafe" assertion variant should panic
	_ = err.(goxidation.OkResult[any])
	t.Errorf("Expected panic by now")
}

func TestErrorCoercionWithSafeCast(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "expected panic")
	}()
	err := goxidation.Error[any](fmt.Errorf("some error"))
	// use the "safe" assertion to try to get an OkResult value
	okRes, ok := err.(goxidation.OkResult[any])
	assert.False(t, ok)
	assert.Nil(t, okRes)
	// trying to use the okRes will panic on a nil pointer
	_ = okRes.Get()
	t.Errorf("Expected panic by now")
}

func TestOkCoercion(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "expected panic")
	}()
	ok := goxidation.Ok[any]("val")
	// "unsafe" assertion variant should panic
	_ = ok.(goxidation.ErrResult[any])
	t.Errorf("Expected panic by now")
}

func TestOkCoercionWithSafeCast(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "expected panic")
	}()
	okRes := goxidation.Ok[any]("val")
	// use the "safe" assertion to try to get an OkResult value
	errRes, ok := okRes.(goxidation.ErrResult[any])
	assert.False(t, ok)
	assert.Nil(t, errRes)
	// trying to use the okRes will panic on a nil pointer
	_ = errRes.Unwrap()
	t.Errorf("Expected panic by now")
}
