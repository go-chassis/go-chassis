package client

import (
	"testing"
)

func BenchmarkRoundRobin(b *testing.B) {
	eps := []string{"172.0.0.1", "172.0.0.2", "172.0.0.3", "172.0.0.4",
		"172.0.0.5", "172.0.0.6", "172.0.0.7", "172.0.0.8", "172.0.0.9",
		"172.0.0.10", "172.0.0.11", "172.0.0.12"}
	next := RoundRobin(eps)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = next()
	}
}
