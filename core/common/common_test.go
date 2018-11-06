package common_test

import (
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
}
