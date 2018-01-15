package restful

import (
	"context"
	"github.com/emicklei/go-restful"
	"net/http"
)

//Context is a struct which has both request and response objects
type Context struct {
	ctx  context.Context
	req  *restful.Request
	resp *restful.Response
}

//NewBaseServer is a function which return context
func NewBaseServer(ctx context.Context) *Context {
	return &Context{
		ctx: ctx,
	}
}

// write is the response writer.
func (bs *Context) Write(body []byte) {
	bs.resp.Write(body)
}

//WriteHeader is the response head writer
func (bs *Context) WriteHeader(httpStatus int) {
	bs.resp.WriteHeader(httpStatus)
}

//AddHeader is a function used to add header to a response
func (bs *Context) AddHeader(header string, value string) {
	bs.resp.AddHeader(header, value)
}

//WriteError is a function used to write error into a response
func (bs *Context) WriteError(httpStatus int, err error) error {
	return bs.resp.WriteError(httpStatus, err)
}

// WriteJSON used to write a JSON file into response
func (bs *Context) WriteJSON(value interface{}, contentType string) error {
	return bs.resp.WriteJson(value, contentType)
}

// WriteHeaderAndJSON used to write head and JSON file in to response
func (bs *Context) WriteHeaderAndJSON(status int, value interface{}, contentType string) error {
	return bs.resp.WriteHeaderAndJson(status, value, contentType)
}

//ReadEntity is request reader
func (bs *Context) ReadEntity(schema interface{}) (err error) {
	return bs.req.ReadEntity(schema)
}

//ReadHeader is used to read header of request
func (bs *Context) ReadHeader(name string) string {
	return bs.req.HeaderParameter(name)
}

//ReadPathParameter is used to read path parameter of a request
func (bs *Context) ReadPathParameter(name string) string {
	return bs.req.PathParameter(name)
}

//ReadPathParameters used to read multiple path parameters of a request
func (bs *Context) ReadPathParameters() map[string]string {
	return bs.req.PathParameters()
}

//ReadQueryParameter is used to read query parameter of a request
func (bs *Context) ReadQueryParameter(name string) string {
	return bs.req.QueryParameter(name)
}

//ReadBodyParameter used to read body parameter of a request
func (bs *Context) ReadBodyParameter(name string) (string, error) {
	return bs.req.BodyParameter(name)
}

//ReadRequest used to read the request
func (bs *Context) ReadRequest() *http.Request {
	return bs.req.Request
}
