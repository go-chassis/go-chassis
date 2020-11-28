package httputil_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/pkg/util/httputil"
	"github.com/stretchr/testify/assert"
)

func TestHttpRequest(t *testing.T) {

	t.Run("convert invocation with ctx to http request,should success",
		func(t *testing.T) {
			r, err := rest.NewRequest("GET", "http://hello", nil)
			assert.NoError(t, err)
			inv := &invocation.Invocation{
				Args: r,
				Ctx: common.NewContext(map[string]string{
					"os":   "mac",
					"user": "peter",
				})}
			r2, err := httputil.HTTPRequest(inv)
			assert.NoError(t, err)
			assert.Equal(t, "mac", r2.Header.Get("os"))
			httputil.SetURI(r, "http://example.com")
			assert.Equal(t, "http://example.com", r.URL.String())

		})
	t.Run("set wrong type to invocation,should fail",
		func(t *testing.T) {
			inv := &invocation.Invocation{
				Ctx: common.NewContext(map[string]string{
					"os":   "mac",
					"user": "peter",
				})}
			_, err := httputil.HTTPRequest(inv)
			assert.Error(t, err)
		})
}

func TestSetURI(t *testing.T) {
	req := &http.Request{}

	t.Run("set wrong url to request", func(t *testing.T) {
		httputil.SetURI(req, "127.0.0.1:8080")
		assert.Nil(t, req.URL)
	})

	t.Run("set right url to request", func(t *testing.T) {
		httputil.SetURI(req, "http://127.0.0.1:8080")
		assert.NotNil(t, req.URL)
		assert.Equal(t, req.URL.Host, "127.0.0.1:8080")
		assert.Equal(t, req.URL.Scheme, "http")
	})

}
func TestSetBody(t *testing.T) {
	req := &http.Request{}
	t.Run("set data of body to request", func(t *testing.T) {
		httputil.SetBody(req, nil)
		body, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.NotNil(t, body)
		assert.Zero(t, len(body))

		data := map[string]string{
			"Test1": "test1",
			"Test2": "test2",
		}
		b, err := json.Marshal(data)
		assert.Nil(t, err)
		httputil.SetBody(req, b)
		body, err = ioutil.ReadAll(req.Body)
		assert.NotNil(t, body)
		assert.NotZero(t, len(body))
		assert.Equal(t, b, body)
	})
}
func TestSetGetCookie(t *testing.T) {
	req := &http.Request{
		Header: make(map[string][]string),
	}
	t.Run("set data to req cookie", func(t *testing.T) {
		ck := "cookie_key"
		cv := "cookie_value"
		httputil.SetCookie(req, ck, cv)

		v, err := req.Cookie(ck)
		assert.Nil(t, err)
		assert.Equal(t, v.Value, cv)

		gv := httputil.GetCookie(req, ck)
		assert.NotEmpty(t, gv)
		assert.Equal(t, gv, cv)
	})

}

func TestSetContentType(t *testing.T) {
	req := &http.Request{
		Header: make(map[string][]string),
	}
	t.Run("set value to req.header of ContentType", func(t *testing.T) {
		ct := "application/json"
		httputil.SetContentType(req, ct)
		cv := httputil.GetContentType(req)
		assert.Equal(t, cv, ct)
	})
}
func TestHTTPRequest(t *testing.T) {

	t.Run("get request for inv args nil will reply error , not nil reply request", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), common.ContextHeaderKey{}, map[string]string{
			"user":    "peter",
			"address": "beijing",
		})
		inv := &invocation.Invocation{
			Args: nil,
			Ctx:  ctx,
		}
		req, err := httputil.HTTPRequest(inv)
		assert.NotNil(t, err)
		assert.Equal(t, err, httputil.ErrInvalidReq)
		assert.Nil(t, req)

		inv.Args = &http.Request{
			Header: make(map[string][]string),
		}
		req, err = httputil.HTTPRequest(inv)
		assert.Nil(t, err)
		assert.NotNil(t, req)
		rv := req.Header.Get("user")
		assert.Equal(t, rv, "peter")
		rv = req.Header.Get("address")
		assert.Equal(t, rv, "beijing")
	})
}

func TestRespBody(t *testing.T) {
	resp := &http.Response{
		Body: nil,
	}
	t.Run("resp or body is nil , did not reply any data ", func(t *testing.T) {
		b := httputil.ReadBody(nil)
		assert.Nil(t, b)
		b = httputil.ReadBody(resp)
		assert.Nil(t, b)
	})
	bodies := []byte("test read resp bodies")
	var bb io.Reader
	t.Run("get data of resp body", func(t *testing.T) {
		bb = bytes.NewReader(bodies)
		rc, ok := bb.(io.ReadCloser)
		if !ok && bodies != nil {
			rc = ioutil.NopCloser(bb)
		}
		resp.Body = rc
		b := httputil.ReadBody(resp)
		assert.NotNil(t, b)
		assert.Equal(t, bodies, b)
	})
}

func TestGetRespCookie(t *testing.T) {
	resp := &http.Response{
		Header: make(map[string][]string),
	}
	cookies := []*http.Cookie{
		{
			Name:  "k1",
			Value: "v1",
		},
		{
			Name:  "k2",
			Value: "v2",
		},
	}
	for _, v := range cookies {
		httputil.SetRespCookie(resp, v)
	}
	t.Run("get exist key for cookie", func(t *testing.T) {
		b := httputil.GetRespCookie(resp, "k1")
		assert.NotNil(t, b)
		assert.Equal(t, b, []byte("v1"))
	})
	t.Run("get not exist key for cookie", func(t *testing.T) {
		b := httputil.GetRespCookie(resp, "k3")
		assert.Nil(t, b)
	})
}
