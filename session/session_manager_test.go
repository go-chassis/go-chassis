package session

import (
	"context"
	"testing"

	"net/http"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
)

func TestGenSessionID(t *testing.T) {
	id, _ := GenerateSessionID()
	t.Log(id)
}
func TestAdd_GetSessionStickinessCache(t *testing.T) {
	SessionStickinessCache = initCache()
	testSlice := []string{"services-cookies1", "services-cookies2"}
	t.Run("set session id to cache and get this", func(t *testing.T) {
		AddSessionStickinessToCache(testSlice[0], "n1")
		cookie := GetSessionID("n1")
		assert.Equal(t, "services-cookies1", cookie)

		AddSessionStickinessToCache(testSlice[1], "n2")
		cookie = GetSessionID("n2")
		assert.Equal(t, "services-cookies2", cookie)

	})
	t.Run("get the namespace but it not exist", func(t *testing.T) {
		cookie := GetSessionID("n3")
		assert.Empty(t, cookie)
	})
}
func TestGetContextMetadata(t *testing.T) {
	t.Run("get value from ctx", func(t *testing.T) {
		ctx := SetContextMetadata(context.Background(), "version", "0.0.1")
		s := GetContextMetadata(ctx, "version")
		assert.Equal(t, s, "0.0.1")

		s = GetContextMetadata(ctx, "name1")
		assert.Empty(t, s)
	})

	t.Run("get value for ctx but ContextHeaderKey did not map[string]string", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), common.ContextHeaderKey{}, map[string]int{
			"version": 1.0,
		})
		s := GetContextMetadata(ctx, "version")
		assert.Empty(t, s)
	})
}
func TestGetSessionFromResp(t *testing.T) {
	resp := &http.Response{
		Header: make(map[string][]string),
	}
	cookie := &http.Cookie{
		Name:  "k1",
		Value: "v1",
	}
	resp.Header.Set("Set-Cookie", cookie.String())
	t.Run("get session form response cookie", func(t *testing.T) {
		s := GetSessionFromResp("k1", resp)
		assert.Equal(t, s, "v1")
		s = GetSessionFromResp("k2", resp)
		assert.Empty(t, s)
	})
}

func TestSaveSessionIDFromContext(t *testing.T) {

	t.Run("save session id from context , cache exist", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), common.ContextHeaderKey{}, map[string]string{
			common.LBSessionID: common.LBSessionID + "=09060c01-040a-4f03-490b-0307060e0008",
		})
		Cache.Set("09060c01-040a-4f03-490b-0307060e0008", "127.0.0.1:8080", 0)

		ctx = SaveSessionIDFromContext(ctx, "127.0.0.1:8080", 10)
		cacheSessionID := getLBCookie(common.LBSessionID)
		assert.NotEmpty(t, cacheSessionID)
		assert.Equal(t, cacheSessionID, common.LBSessionID+"=09060c01-040a-4f03-490b-0307060e0008")
	})

	t.Run("save session id from context , cache not exist", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), common.ContextHeaderKey{}, map[string]string{
			common.LBSessionID: common.LBSessionID + "=09060c01-040a-4f03-490b-0307060e0001",
		})
		ctx = SaveSessionIDFromContext(ctx, "127.0.0.1:8080", 10)
		cacheSessionID := getLBCookie(common.LBSessionID)
		assert.NotEmpty(t, cacheSessionID)
		assert.NotEqual(t, cacheSessionID, common.LBSessionID+"=09060c01-040a-4f03-490b-0307060e0001")
	})

}

func TestSaveSessionIDFromHTTP(t *testing.T) {

	resp := &http.Response{
		Header: make(map[string][]string),
	}
	req := &http.Request{
		Header: make(map[string][]string),
	}

	cookie := &http.Cookie{
		Name: common.LBSessionID,
	}

	t.Run("save session id from response", func(t *testing.T) {

		cookie.Value = "09060c01-040a-4f03-490b-0307060e0008"
		resp.Header.Add("Set-Cookie", cookie.String())
		SaveSessionIDFromHTTP("127.0.0.1:8080", 0, resp, req)
		c, ok := Cache.Get("09060c01-040a-4f03-490b-0307060e0008")
		assert.True(t, ok)
		v, ok := c.(string)
		assert.True(t, ok)
		assert.Equal(t, v, "127.0.0.1:8080")
	})

	Cache.Flush()

	t.Run("save session id from request", func(t *testing.T) {
		resp.Header = make(map[string][]string)
		cookie.Value = "07020b07-0b08-4e0e-410c-0908060a0e0e"
		req.AddCookie(cookie)
		Cache.Set("07020b07-0b08-4e0e-410c-0908060a0e0e", "127.0.0.1:8888", 0)
		SaveSessionIDFromHTTP("127.0.0.1:8888", 0, resp, req)
		c, ok := Cache.Get(cookie.Value)
		assert.True(t, ok)
		v, ok := c.(string)
		assert.True(t, ok)
		assert.Equal(t, v, "127.0.0.1:8888")
	})
	Cache.Flush()
	t.Run("save session id by new id", func(t *testing.T) {
		resp.Header = make(map[string][]string)
		SaveSessionIDFromHTTP("127.0.0.1:9999", 0, resp, &http.Request{})
		for k, v := range Cache.Items() {
			assert.NotEqual(t, k, "go-chassisLB=07020b07-0b08-4e0e-410c-0908060a0e0e")
			assert.NotEqual(t, k, "07020b07-0b08-4e0e-410c-0908060a0e0e")
			assert.NotEqual(t, k, "go-chassisLB=09060c01-040a-4f03-490b-0307060e0008")
			assert.NotEqual(t, k, "09060c01-040a-4f03-490b-0307060e0008")
			assert.Equal(t, v.Object, "127.0.0.1:9999")
		}
	})
}

func TestDeletingKeySuccessiveFailure(t *testing.T) {
	resp := &http.Response{
		Header: make(map[string][]string),
	}
	cookie := &http.Cookie{
		Name:  common.LBSessionID,
		Value: common.LBSessionID + "=0d0a0503-0701-4c07-4108-060a0a070e02",
	}
	resp.Header.Add("Set-Cookie", cookie.String())
	Cache.Set("0d0a0503-0701-4c07-4108-060a0a070e02", "127.0.0.1:8888", 0)
	setLBCookie(common.LBSessionID, common.LBSessionID+"=0d0a0503-0701-4c07-4108-060a0a070e02")

	t.Run("delete session id by resp", func(t *testing.T) {
		DeletingKeySuccessiveFailure(resp)
		i, b := Cache.Get("0d0a0503-0701-4c07-4108-060a0a070e02")
		assert.False(t, b)
		assert.Nil(t, i)
		s := getLBCookie(common.LBSessionID)
		assert.NotEmpty(t, s)
		assert.Equal(t, common.LBSessionID+"=0d0a0503-0701-4c07-4108-060a0a070e02", s)
	})
	t.Run("delete session id for cookieMap", func(t *testing.T) {
		DeletingKeySuccessiveFailure(nil)
		assert.Empty(t, getLBCookie(common.LBSessionID))
	})

}
func TestGetSessionIDFromInv(t *testing.T) {
	ctx := context.WithValue(context.TODO(), common.ContextHeaderKey{}, map[string]string{
		common.LBSessionID: common.LBSessionID + "=09060c01-040a-4f03-490b-0307060e0008",
	})
	cookie := &http.Cookie{
		Name:  common.LBSessionID,
		Value: common.LBSessionID + "=0d0a0503-0701-4c07-4108-060a0a070e02",
	}
	resp := &http.Response{
		Header: make(map[string][]string),
	}
	resp.Header.Add("Set-Cookie", cookie.String())
	inv := invocation.Invocation{
		Reply: resp,
	}

	t.Run("get session id from inv.Reply", func(t *testing.T) {
		sessionID := GetSessionIDFromInv(inv, common.LBSessionID)
		assert.Equal(t, sessionID, common.LBSessionID+"=0d0a0503-0701-4c07-4108-060a0a070e02")
	})
	inv.Ctx = ctx
	inv.Reply = nil
	t.Run("get session id from inv.Ctx", func(t *testing.T) {
		sessionID := GetSessionIDFromInv(inv, common.LBSessionID)
		assert.Equal(t, sessionID, common.LBSessionID+"=09060c01-040a-4f03-490b-0307060e0008")
	})
}
