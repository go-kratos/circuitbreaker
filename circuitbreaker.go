package circuitbreaker

// State .
type State int

const (
	// StateOpen when circuit breaker open, request not allowed, after sleep
	// some duration, allow one single request for testing the health, if ok
	// then state reset to closed, if not continue the step.
	StateOpen State = iota
	// StateClosed when circuit breaker closed, request allowed, the breaker
	// calc the succeed ratio, if request num greater request setting and
	// ratio lower than the setting ratio, then reset state to open.
	StateClosed
	// StateHalfopen when circuit breaker open, after slepp some duration, allow
	// one request, but not state closed.
	StateHalfopen
)

// CircuitBreaker .
type CircuitBreaker interface {
	Allow() error
	MarkSuccess()
	MarkFailed()
}

// Options .
type Options func(o *options)

type options struct {
	State func(o, n State)
}

// OnState .
func OnState(fn func(o, n State)) Options {
	return func(o *options) { o.State = fn }
}

// New .
func New(opts ...Options) CircuitBreaker {
	return nil
}

// Do .
// err := Do("xx", func() error {
//    return Ignore(errors.New("ok"))
// })
func Do(name string, fn func() error, fbs ...func(error) error) error {
	return nil
}

type ignore struct {
	error
}

// Ignore .
func Ignore(err error) error {
	return ignore{err}
}

type drop struct {
	error
}

// Drop .
func Drop(err error) error {
	return drop{err}
}
