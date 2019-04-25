package loadbalancer_test

import (
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"gopkg.in/go-playground/assert.v1"
	"testing"
	"time"
)

func TestProtocolStats_SaveLatency(t *testing.T) {
	s := &loadbalancer.ProtocolStats{
		Latency: make([]time.Duration, 0),
		Addr:    "127.0.0.1:8080",
	}
	s.CalculateAverageLatency()

	s.SaveLatency(1 * time.Second)
	s.SaveLatency(2 * time.Second)
	s.CalculateAverageLatency()
	assert.Equal(t, 1500*time.Millisecond, s.AvgLatency)

	s.SaveLatency(2 * time.Second)
	s.SaveLatency(2 * time.Second)
	s.SaveLatency(2 * time.Second)
	s.SaveLatency(2 * time.Second)
	s.SaveLatency(2 * time.Second)
	s.SaveLatency(2 * time.Second)
	s.SaveLatency(2 * time.Second)
	s.SaveLatency(2 * time.Second)
	s.CalculateAverageLatency()
	s.SaveLatency(3 * time.Second)
	s.SaveLatency(3 * time.Second)
	s.CalculateAverageLatency()
	assert.Equal(t, 2200*time.Millisecond, s.AvgLatency)

}
