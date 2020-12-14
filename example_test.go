package circuitbreaker_test

import (
	"errors"
	"fmt"

	circuitbreaker "github.com/go-kratos/circuitbreaker/v1"
)

// This is a example of using a circuit breaker Do() when return nil.
func ExampleDo() {
	err := circuitbreaker.Do("do", func() error {
		// dosomething
		return nil
	})

	fmt.Printf("err=%v", err)
	// Output: err=<nil>
}

// This is a example of using a circuit breaker fn failed then call fallback.
func ExampleDo_fallback() {
	err := circuitbreaker.Do("do", func() error {
		// dosomething
		return errors.New("fallback")
	}, func(err error) error {
		return nil
	})

	fmt.Printf("err=%v", err)
	// Output: err=fallback
}

// This is a example of using a circuit breaker fn failed but ignore error mark
// as success.
func ExampleDo_ignore() {
	err := circuitbreaker.Do("do", func() error {
		// dosomething
		return circuitbreaker.Ignore(errors.New("fallback"))
	}, func(err error) error {
		return errors.New("fallback") // won`t touch here
	})

	fmt.Printf("err=%v", err)
	// Output: err=<nil>
}
