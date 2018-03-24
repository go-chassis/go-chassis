package archaius_test

import (
	"strings"
	"testing"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/stretchr/testify/assert"
)

func TestGetSpecificKey(t *testing.T) {
	cmd := strings.Join([]string{common.Consumer, "Carts"}, ".")
	key := archaius.GetHystrixSpecificKey(archaius.NamespaceIsolation, cmd, archaius.PropertyTimeoutInMilliseconds)
	assert.Equal(t, "cse.isolation.Consumer.Carts."+archaius.PropertyTimeoutInMilliseconds, key)

}
