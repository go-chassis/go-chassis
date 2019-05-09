package common_test

import (
	"context"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := common.NewContext(map[string]string{
		"1": "2",
	})
	m := common.FromContext(ctx)
	assert.Equal(t, "2", m["1"])

	ctx = common.WithContext(ctx, "3", "4")

	m = common.FromContext(ctx)
	assert.Equal(t, "2", m["1"])
	assert.Equal(t, "4", m["3"])

	ctx = common.NewContext(nil)
	m = common.FromContext(ctx)
	assert.NotNil(t, m)

	ctx = common.WithContext(nil, "test", "1")
	m = common.FromContext(ctx)
	assert.Equal(t, "1", m["test"])

	t.Run("convert nil context, it return new map", func(t *testing.T) {
		m = common.FromContext(nil)
		assert.Equal(t, 0, len(m))
	})
	t.Run("set kv with wrong context, it return context", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), common.ContextHeaderKey{}, 1)
		ctx = common.WithContext(ctx, "os", "mac")
		assert.Equal(t, "mac", common.FromContext(ctx)["os"])
	})
}
