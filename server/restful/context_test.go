package restful

import (
	"testing"
	//"context"
	"github.com/emicklei/go-restful"
	"golang.org/x/net/context"
	"net/http"

	"errors"
	"github.com/stretchr/testify/assert"
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
	contxt := NewBaseServer(context.TODO())
	contxt.req = &restful.Request{Request: &http.Request{Method: "Get"}}
	rw := contextTest{}
	resp := restful.NewResponse(rw)
	contxt.resp = resp
	contxt.AddHeader("Content-Type", "application/json")

	_, er := contxt.ReadBodyParameter("hello")
	assert.NoError(t, er)

	paramVal := contxt.ReadPathParameter("abc")
	assert.Empty(t, paramVal)

	param := contxt.ReadPathParameters()
	assert.Empty(t, param)

	contxt.WriteHeader(200)

	val := contxt.ReadHeader("Content-Type")
	assert.Empty(t, val)

	req := contxt.ReadRequest()
	assert.NotEmpty(t, req)

	contxt.Write([]byte("success"))

	err := contxt.ReadEntity("hhhh")
	assert.Error(t, err)

	err = contxt.WriteError(200, errors.New("error"))
	assert.NoError(t, err)

	err = contxt.WriteHeaderAndJSON(204, "deleted", "success")
	assert.NoError(t, err)

	err = contxt.WriteJSON("json", "application")
	assert.NoError(t, err)

	query := contxt.ReadQueryParameter("hhh")
	assert.Empty(t, query)

}
