package clock_test

import (
	"github.com/ServiceComb/go-chassis/core/qpslimiter"
	"github.com/ServiceComb/go-chassis/third_party/forked/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"

	"sync/atomic"
	"time"
	//"fmt"
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

func TestClock(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	clock := clock.NewMock()
	rl := qpslimiter.New(100, qpslimiter.WithClock(clock), qpslimiter.WithoutSlack)

	count := NewInt32(0)

	// Until we're done...
	done := make(chan struct{})
	defer close(done)

	// Create copious counts concurrently.
	go job(rl, count, done)
	go job(rl, count, done)
	go job(rl, count, done)
	go job(rl, count, done)

	clock.AfterFunc(1*time.Second, func() {
		assert.InDelta(t, 100, count.Load(), 10, "count within rate limit")
	})

	clock.AfterFunc(2*time.Second, func() {
		assert.InDelta(t, 200, count.Load(), 10, "count within rate limit")
	})

	clock.AfterFunc(3*time.Second, func() {
		assert.InDelta(t, 300, count.Load(), 10, "count within rate limit")
		wg.Done()
	})

	clock.Add(4 * time.Second)

	clock.Add(5 * time.Second)
}

func TestTimer(t *testing.T) {
	clock := clock.NewMock()
	ti := clock.Timer(3 * time.Second)
	assert.NotEmpty(t, ti)
	//fmt.Println("ti is ******************", ti)
	//fmt.Println("ti is %+v******************", ti)
	//fmt.Println("ti is ******************", &ti)
	//assert.Equal(t, ti, ti)

}
func job(rl qpslimiter.Limiter, count *Int32, done <-chan struct{}) {
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
