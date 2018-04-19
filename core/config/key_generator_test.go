package config

import (
	"strings"
	"testing"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/stretchr/testify/assert"
)

func TestGetSpecificKey(t *testing.T) {
	cmd := strings.Join([]string{common.Consumer, "Carts"}, ".")
	key := GetHystrixSpecificKey(NamespaceIsolation, cmd, PropertyTimeoutInMilliseconds)
	assert.Equal(t, "cse.isolation.Consumer.Carts."+PropertyTimeoutInMilliseconds, key)
}
