package common_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/stretchr/testify/assert"
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
func TestXCseHeader(t *testing.T) {
	req := &http.Request{
		Header: map[string][]string{},
	}
	t.Run("set value of X into req header , use GetXCSEHeader func get value reply X", func(t *testing.T) {
		common.SetXCSEHeader(common.HeaderSourceName, "test1", req)
		s := common.GetXCSEHeader(common.HeaderSourceName, req)
		assert.NotEmpty(t, s)
		assert.Equal(t, s, "test1")
	})

	t.Run("test the empty value did not overwrite old value", func(t *testing.T) {
		m := map[string]string{
			common.HeaderSourceName: "test2",
		}
		b, _ := json.Marshal(m)

		req.Header.Set(common.HeaderXCseContent, string(b))
		common.SetXCSEHeader(common.HeaderSourceName, "", req)
		s := common.GetXCSEHeader(common.HeaderSourceName, req)
		assert.NotEmpty(t, s)
		assert.Equal(t, s, "test2")

		common.SetXCSEHeader(common.HeaderSourceName, "test3", nil)
		s = common.GetXCSEHeader(common.HeaderSourceName, nil)
		assert.Empty(t, s)

		common.SetXCSEHeader("", "test3", req)
		s = common.GetXCSEHeader("", req)
		assert.Empty(t, s)
	})
	t.Run("test new value will overwrite old value", func(t *testing.T) {
		common.SetXCSEHeader(common.HeaderSourceName, "test4", req)
		s := common.GetXCSEHeader(common.HeaderSourceName, req)
		assert.NotEmpty(t, s)
		assert.Equal(t, s, "test4")
	})

	t.Run("the same value did not overwrite the old value",
		func(t *testing.T) {
			common.SetXCSEHeader(common.HeaderSourceName, "test4", req)
			s := common.GetXCSEHeader(common.HeaderSourceName, req)
			assert.NotEmpty(t, s)
			assert.Equal(t, s, "test4")
		})
}
