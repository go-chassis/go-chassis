package archaius_test

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSpecificKey(t *testing.T) {
	cmd := handler.NewHystrixCmd("vmall", common.Consumer, "Carts", "cartService", "get")
	key := archaius.GetHystrixSpecificKey(archaius.NamespaceIsolation, cmd, archaius.PropertyTimeoutInMilliseconds)
	assert.Equal(t, "cse.isolation.vmall.Consumer.Carts.cartService.get."+archaius.PropertyTimeoutInMilliseconds, key)

}
