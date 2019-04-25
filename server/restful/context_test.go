package restful

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

//type Header map[string][]string

type contextTest struct {
}

func (c contextTest) Header() http.Header {
	h := http.Header{}
	h.Set("Content-Type", "applicarion/json")
	return h
}

func (c contextTest) Write([]byte) (int, error) {
	return 0, nil
}

func (c contextTest) WriteHeader(int) {
	return
}

func TestContextFuncs(t *testing.T) {
	t.Log("Testing all the restful server functions")
	ctx := NewBaseServer(context.TODO())
	ctx.req = &restful.Request{Request: &http.Request{Method: "Get"}}
	rw := httptest.NewRecorder()
	resp := restful.NewResponse(rw)
	ctx.resp = resp
	ctx.AddHeader("Content-Type", "application/json")

	_, er := ctx.ReadBodyParameter("hello")
	assert.NoError(t, er)

	paramVal := ctx.ReadPathParameter("abc")
	assert.Empty(t, paramVal)

	param := ctx.ReadPathParameters()
	assert.Empty(t, param)

	ctx.WriteHeader(200)

	val := ctx.ReadHeader("Content-Type")
	assert.Empty(t, val)

	req := ctx.ReadRequest()
	assert.NotEmpty(t, req)

	ctx.Write([]byte("success"))

	err := ctx.ReadEntity("hhhh")
	assert.Error(t, err)

	err = ctx.WriteError(200, errors.New("error"))
	assert.NoError(t, err)

	err = ctx.WriteHeaderAndJSON(204, "deleted", "success")
	assert.NoError(t, err)

	err = ctx.WriteJSON("json", "application")
	assert.NoError(t, err)

	query := ctx.ReadQueryParameter("hhh")
	assert.Empty(t, query)

	type queryRequest struct {
		Name     string `form:"name"`
		Password string `form:"password"`
	}
	var queryReq queryRequest
	expectReq := queryRequest{Name: "admin", Password: "admin"}
	url, _ := url.Parse("http://127.0.0.1/test?name=admin&password=admin")
	ctx.req.Request.URL = url
	err = ctx.ReadQueryEntity(&queryReq)
	assert.NoError(t, err)
	assert.Equal(t, expectReq, queryReq)

	rreq := ctx.ReadRestfulRequest()
	assert.Equal(t, "127.0.0.1", rreq.Request.URL.Host)

	rresp := ctx.ReadRestfulResponse()
	assert.NotNil(t, rresp)

	wr := ctx.ReadResponseWriter()
	assert.NotNil(t, wr)
}
