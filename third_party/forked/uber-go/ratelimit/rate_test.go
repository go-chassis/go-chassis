package ratelimit_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ServiceComb/go-chassis/third_party/forked/benbjohnson/clock"
	"github.com/ServiceComb/go-chassis/third_party/forked/uber-go/ratelimit"
	"github.com/stretchr/testify/assert"
)

// Int32 is an atomic wrapper around an int32.
type Int32 struct{ v int32 }

// NewInt32 creates an Int32.
func NewInt32(i int32) *Int32 {
	return &Int32{i}
}

// Load atomically loads the wrapped value.
func (i *Int32) Load() int32 {
	return atomic.LoadInt32(&i.v)
}

// Inc atomically increments the wrapped int32 and returns the new value.
func (i *Int32) Inc() int32 {
	return i.Add(1)
}

// Add atomically adds to the wrapped int32 and returns the new value.
func (i *Int32) Add(n int32) int32 {
	return atomic.AddInt32(&i.v, n)
}

func ExampleRatelimit() {
	rl := ratelimit.New(100) // per second

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		if i > 0 {
			fmt.Println(i, now.Sub(prev))
		}
		prev = now
	}

	// Output:
	// 1 10ms
	// 2 10ms
	// 3 10ms
	// 4 10ms
	// 5 10ms
	// 6 10ms
	// 7 10ms
	// 8 10ms
	// 9 10ms
}

//func TestRateLimiter(t *testing.T) {
//	var wg sync.WaitGroup
//	wg.Add(1)
//	defer wg.Wait()
//
//	clock := clock.NewMock()
//	rl := qpslimiter.New(100, qpslimiter.WithClock(clock), qpslimiter.WithoutSlack)
//
//	count := NewInt32(0)
//
//	// Until we're done...
//	done := make(chan struct{})
//	defer close(done)
//
//	// Create copious counts concurrently.
//	go job(rl, count, done)
//	go job(rl, count, done)
//	go job(rl, count, done)
//	go job(rl, count, done)
//
//	clock.AfterFunc(1*time.Second, func() {
//		assert.InDelta(t, 100, count.Load(), 10, "count within rate limit")
//	})
//
//	clock.AfterFunc(2*time.Second, func() {
//		assert.InDelta(t, 200, count.Load(), 10, "count within rate limit")
//	})
//
//	clock.AfterFunc(3*time.Second, func() {
//		assert.InDelta(t, 300, count.Load(), 10, "count within rate limit")
//		wg.Done()
//	})
//
//	clock.Add(4 * time.Second)
//
//	clock.Add(5 * time.Second)
//}

func TestDelayedRateLimiter(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	clock := clock.NewMock()
	slow := ratelimit.New(10, ratelimit.WithClock(clock))
	fast := ratelimit.New(100, ratelimit.WithClock(clock))

	count := NewInt32(0)

	// Until we're done...
	done := make(chan struct{})
	defer close(done)

	// Run a slow job
	go func() {
		for {
			slow.Take()
			fast.Take()
			count.Inc()
			select {
			case <-done:
				return
			default:
			}
		}
	}()

	// Accumulate slack for 10 seconds,
	clock.AfterFunc(20*time.Second, func() {
		// Then start working.
		go job(fast, count, done)
		go job(fast, count, done)
		go job(fast, count, done)
		go job(fast, count, done)
	})

	clock.AfterFunc(30*time.Second, func() {
		assert.InDelta(t, 1200, count.Load(), 10, "count within rate limit")
		wg.Done()
	})

	clock.Add(40 * time.Second)
}

func job(rl ratelimit.Limiter, count *Int32, done <-chan struct{}) {
	for {
		rl.Take()
		count.Inc()
		select {
		case <-done:
			return
		default:
		}
	}
}
