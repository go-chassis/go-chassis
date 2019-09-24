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
	req := &http.Request{}
	t.Run("set value of X into req header , use GetXCSEContext func get value reply X", func(t *testing.T) {
		common.SetXCSEContext(map[string]string{common.HeaderSourceName: "test1"}, req)
		s := common.GetXCSEContext(common.HeaderSourceName, req)
		assert.NotEmpty(t, s)
		assert.Equal(t, s, "test1")
	})

	t.Run("test the empty value ", func(t *testing.T) {
		m := map[string]string{
			common.HeaderSourceName: "test2",
		}
		b, _ := json.Marshal(m)

		req.Header.Set(common.HeaderXCseContent, string(b))
		common.SetXCSEContext(map[string]string{common.HeaderSourceName: ""}, req)
		s := common.GetXCSEContext(common.HeaderSourceName, req)
		assert.Empty(t, s)
		assert.Equal(t, s, "")

		common.SetXCSEContext(map[string]string{common.HeaderSourceName: "test3"}, nil)
		common.SetXCSEContext(map[string]string{"": "test3"}, req)
		s = common.GetXCSEContext("", req)
		assert.Equal(t, s, "test3")

	})
	t.Run("input param req is nil or req.Header is nil , will return empty", func(t *testing.T) {
		s := common.GetXCSEContext(common.HeaderSourceName, nil)
		assert.Empty(t, s)

		s = common.GetXCSEContext(common.HeaderSourceName, &http.Request{})
		assert.Empty(t, s)

	})
	t.Run("test new value will overwrite old value", func(t *testing.T) {
		common.SetXCSEContext(map[string]string{common.HeaderSourceName: "test4"}, req)
		s := common.GetXCSEContext(common.HeaderSourceName, req)
		assert.NotEmpty(t, s)
		assert.Equal(t, s, "test4")
	})

	t.Run("the same value did not overwrite the old value",
		func(t *testing.T) {
			common.SetXCSEContext(map[string]string{common.HeaderSourceName: "test4"}, req)
			s := common.GetXCSEContext(common.HeaderSourceName, req)
			assert.NotEmpty(t, s)
			assert.Equal(t, s, "test4")
		})
	t.Run("test old version and new version transfer source service name ", func(t *testing.T) {
		req.Header = make(map[string][]string)
		req.Header.Set(common.HeaderSourceName, "test5")
		s := common.GetXCSEContext(common.HeaderSourceName, req)
		assert.NotEmpty(t, s)
		assert.Equal(t, s, "test5")
	})
}
