# goxidation

An exercise/thought experiment of recreating the core Result and Option types from Rust in Go.

There are many like it, but this one is mine.

This is not intended for use or adoption.

## Goals

Mimic base protectons of the Rust Result and Option types:

- an Err Result[T] or None Option[T] cannot be coerced into a T without causing a panic

## Non-Goals

- Exact API

## Design

- `Result[T]` and `Option[T]` are implemented as Interfaces.
  - These interfaces provide methods to determine which variant an instance is
  - A `GetOr(T) T` convenience method for retrieving the contained value or the provided default value
  - And an unexported method that references an internal type prevents external implementations of `Result[T]`, `Option[T]`, and the variant interferfaces
- Implementation types of `Result[T]` and `Option[T]` are unexported to prevent misuse
- `Ok(T)` and `Err[T](error)` instantiate the `Result[T]` implementation types
- `Some(T)` and `None[T]()` instantiate the `Option[T]` implementation types
- `Result[T]` and `Option[T]` instances can be type-asserted into the correct interface variant
  - Incorrect single-return type-assertions panic immediately preventing misuse
  - Incorrect two-return type-assertions result in a nil interface which will panic upon use, preventing misuse
- `OkResult[T]` provides a `Get() T` method to retrieve the value
- `ErrResult[T]` provides methods to implement the `Error` and `Unwrapper` interfaces for compatibility with `errors.Is` and `errors.As`. `Unwrap()` is also used to get the original error for handling.
  - The overlap of the `Unwrap() error` method name is why I use `Get() T` instead of `Unwrap() T`
- `SomeOption[T]` provides a `Get() T` method to retrieve the value
- `NoneOption[T]` has no usable methods

## Usage

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nabowler/goxidation"
)

func main() {
	result := Get("https://example.com")
	if result.IsError() {
		asError := result.(goxidation.ErrResult[*http.Response])
		os.Stderr.WriteString(asError.Error())
		os.Exit(1)
	}

	asOk := result.(goxidation.OkResult[*http.Response])
	resp := asOk.Get()
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

//nolint:bodyclose
func Get(url string) goxidation.Result[*http.Response] {
	resp, err := http.Get(url)
	if err != nil {
		return goxidation.Error[*http.Response](err)
	}
	return goxidation.Ok(resp)
}
```