package qpslimiter

import (
	"sync"
	"time"

	"github.com/ServiceComb/go-chassis/third_party/forked/benbjohnson/clock"
)

// Limiter is used to rate-limit some process, possibly across goroutines.
// The process is expected to call Take() before every iteration, which
// may block to throttle the goroutine.
type Limiter interface {
	// Take should block to make sure that the RPS is met.
	Take() time.Time
}

// Clock is the minimum necessary interface to instantiate a rate limiter with
// a clock or mock clock, compatible with clocks created using
// github.com/andres-erbsen/clock.
type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

type limiter struct {
	sync.Mutex
	last       time.Time
	sleepFor   time.Duration
	perRequest time.Duration
	maxSlack   time.Duration
	clock      Clock
}

// Option configures a Limiter.
type Option func(l *limiter)

// New returns a Limiter that will limit to the given RPS.
func New(rate int, opts ...Option) Limiter {
	l := &limiter{
		perRequest: time.Second / time.Duration(rate),
		maxSlack:   -10 * time.Second / time.Duration(rate),
	}
	for _, opt := range opts {
		opt(l)
	}
	if l.clock == nil {
		l.clock = clock.New()
	}
	return l
}

// WithClock returns an option for ratelimit.New that provides an alternate
// Clock implementation, typically a mock Clock for testing.
func WithClock(clock Clock) Option {
	return func(l *limiter) {
		l.clock = clock
	}
}

// WithoutSlack is an option for ratelimit.New that initializes the limiter
// without any initial tolerance for bursts of traffic.
var WithoutSlack Option = withoutSlackOption

func withoutSlackOption(l *limiter) {
	l.maxSlack = 0
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (t *limiter) Take() time.Time {
	t.Lock()
	defer t.Unlock()

	now := t.clock.Now()

	// If this is our first request, then we allow it.
	if t.last.IsZero() {
		t.last = now
		return t.last
	}

	// sleepFor calculates how much time we should sleep based on
	// the perRequest budget and how long the last request took.
	// Since the request may take longer than the budget, this number
	// can get negative, and is summed across requests.
	t.sleepFor += t.perRequest - now.Sub(t.last)

	// We shouldn't allow sleepFor to get too negative, since it would mean that
	// a service that slowed down a lot for a short period of time would get
	// a much higher RPS following that.
	if t.sleepFor < t.maxSlack {
		t.sleepFor = t.maxSlack
	}

	// If sleepFor is positive, then we should sleep now.
	if t.sleepFor > 0 {
		t.clock.Sleep(t.sleepFor)
		t.last = now.Add(t.sleepFor)
		t.sleepFor = 0
	} else {
		t.last = now
	}

	return t.last
}
