package config_test

import (
	"strings"
	"testing"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/stretchr/testify/assert"
)

func TestGetSpecificKey(t *testing.T) {
	cmd := strings.Join([]string{common.Consumer, "Carts"}, ".")
	key := config.GetHystrixSpecificKey(config.NamespaceIsolation, cmd, config.PropertyTimeoutInMilliseconds)
	assert.Equal(t, "cse.isolation.Consumer.Carts."+config.PropertyTimeoutInMilliseconds, key)
	t.Log(key)
}
