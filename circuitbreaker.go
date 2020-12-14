package circuitbreaker

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	_group = &Group{New: func() CircuitBreaker { return New() }}
	// ErrNotAllowed error not allowed.
	ErrNotAllowed = errors.New("circuitbreaker: not allowed for circuit open")
)

func init() {
	_group.val.Store(make(map[string]CircuitBreaker))
}

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

// Group .
type Group struct {
	mutex sync.Mutex
	val   atomic.Value

	New func() CircuitBreaker
}

// Get .
func (g *Group) Get(name string) CircuitBreaker {
	v := g.val.Load().(map[string]CircuitBreaker)
	cb, ok := v[name]
	if ok {
		return cb
	}
	// slowpath for group don`t have specified name breaker.
	g.mutex.Lock()
	nv := make(map[string]CircuitBreaker, len(v)+1)
	for i, j := range v {
		nv[i] = j
	}
	cb = g.New()
	nv[name] = cb
	g.val.Store(nv)
	g.mutex.Unlock()
	return cb
}

// Do runs your function in a synchronous manner, blocking until either your
// function succeeds or an error is returned, including circuit errors.
func Do(name string, fn func() error, fbs ...func(error) error) error {
	cb := _group.Get(name)
	err := cb.Allow()
	if err != nil {
		return err
	}
	if err = fn(); err == nil {
		cb.MarkSuccess()
		return nil
	}
	switch err.(type) {
	case ignore:
		cb.MarkSuccess()
		return err.(ignore).error
	case drop:
		return err.(drop).error
	default:
		cb.MarkFailed()
	}
	oerr := err // origin error
	for _, fb := range fbs {
		if err = fb(oerr); err == nil {
			return nil
		}
	}
	return err
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
