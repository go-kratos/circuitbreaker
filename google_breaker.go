package circuitbreaker

import (
	"math/rand"
	"sync"
	"time"
)

// googleBreaker is a google sre CircuitBreaker pattern.
type googleBreaker struct {
	// stat metric.RollingCounter
	r *rand.Rand
	// rand.New(...) returns a non thread safe object
	randLock sync.Mutex

	opts *options
	k    float64

	state State
}

func newGoogleBreaker(opts *options) CircuitBreaker {
	return &googleBreaker{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),

		opts: opts,
		k:    1 / opts.Success,

		state: StateClosed,
	}
}

func (b *googleBreaker) Allow() error {
	/*
		success, total := b.summary()
		k := b.k * float64(success)
		if log.V(5) {
			log.Info("breaker: request: %d, succee: %d, fail: %d", total, success, total-success)
		}
		// check overflow requests = K * success
		if total < b.request || float64(total) < k {
			if atomic.LoadInt32(&b.state) == StateOpen {
				atomic.CompareAndSwapInt32(&b.state, StateOpen, StateClosed)
			}
			return nil
		}
		if atomic.LoadInt32(&b.state) == StateClosed {
			atomic.CompareAndSwapInt32(&b.state, StateClosed, StateOpen)
		}
		dr := math.Max(0, (float64(total)-k)/float64(total+1))
		b.randLock.Lock()
		drop = b.r.Float64() < dr
		b.randLock.Unlock()
		drop := b.trueOnProba(dr)
		if log.V(5) {
			log.Info("breaker: drop ratio: %f, drop: %t", dr, drop)
		}
		if drop {
			return ecode.ServiceUnavailable
		}
	*/
	return nil
}

func (b *googleBreaker) MarkSuccess() {
	// b.stat.Add(1)
}

func (b *googleBreaker) MarkFailed() {
	// NOTE: when client reject requets locally, continue add counter let the
	// drop ratio higher.
	// b.stat.Add(0)
}
