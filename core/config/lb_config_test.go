package config_test

import (
	// "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStrategyName(t *testing.T) {
	check := config.GetStrategyName("source", "service")
	assert.Equal(t, "WeightedResponse", check)
}

func TestGetRetryOnNext(t *testing.T) {
	check := config.GetRetryOnNext("source", "service")
	assert.Equal(t, 2, check)
}

// GetServerListFilters get server list filters
func BenchmarkGetServerListFilters(b *testing.B) {
	lager.Initialize("", "INFO", "", "size",
		true, 1, 10, 7)

	err := config.InitArchaius()
	assert.NoError(b, err)
	f := config.GetServerListFilters()
	b.Log(f)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.GetServerListFilters()
	}
}

// GetServerListFilters get server list filters
func BenchmarkGetServerListFilters2(b *testing.B) {
	lager.Initialize("", "INFO", "", "size",
		true, 1, 10, 7)

	err := config.InitArchaius()
	assert.NoError(b, err)
	config.ReadLBFromArchaius()
	b.Log(config.GetLoadBalancing().Filters)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GetLoadBalancing().Filters
	}
}
func BenchmarkGetStrategyName(b *testing.B) {
	lager.Initialize("", "INFO", "", "size",
		true, 1, 10, 7)

	err := config.InitArchaius()
	assert.NoError(b, err)
	config.ReadLBFromArchaius()
	b.Log(config.GetStrategyName("", ""))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GetStrategyName("", "")
	}
}
