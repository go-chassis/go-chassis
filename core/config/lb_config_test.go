package config_test

import (
	"github.com/go-chassis/go-chassis/v2/resilience/retry"
	// "github.com/go-chassis/go-chassis/v2/core/common"
	"testing"

	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/stretchr/testify/assert"
)

func TestGetStrategyName(t *testing.T) {
	config.ReadLBFromArchaius()
	check := config.GetStrategyName("service")
	assert.Equal(t, "WeightedResponse", check)

	t.Run("TestGetRetryOnNext", func(t *testing.T) {
		check := config.GetRetryOnNext("source", "service")
		assert.Equal(t, 2, check)
	})

	t.Run("TestRetryEnabled", func(t *testing.T) {
		b := config.RetryEnabled("source", "service")
		assert.Equal(t, false, b)
	})

	t.Run("TestBackOffKind", func(t *testing.T) {
		s := config.BackOffKind("service")
		assert.Equal(t, retry.KindConstant, s)
	})

	t.Run("TestBackOffMaxMs", func(t *testing.T) {
		max := config.BackOffMaxMs("source", "service")
		assert.Equal(t, 400, max)
	})

	t.Run("TestBackOffMinMs",
		func(t *testing.T) {
			min := config.BackOffMinMs("source", "service")
			assert.Equal(t, 200, min)
		})
}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}

// GetServerListFilters get server list filters
func BenchmarkGetServerListFilters(b *testing.B) {

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

	err := config.InitArchaius()
	assert.NoError(b, err)
	config.ReadLBFromArchaius()
	b.Log(config.GetStrategyName(""))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GetStrategyName("")
	}
}
