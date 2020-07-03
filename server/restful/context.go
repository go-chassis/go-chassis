package restful

import (
	"context"
	"github.com/emicklei/go-restful"
	"net/http"
)

//Context is a struct which has both request and response objects
// and request context
type Context struct {
	Ctx  context.Context
	Req  *restful.Request
	Resp *restful.Response
}

//NewBaseServer is a function which return context
func NewBaseServer(ctx context.Context) *Context {
	return &Context{
		Ctx: ctx,
	}
}

// write is the response writer.
func (bs *Context) Write(body []byte) error {
	_, err := bs.Resp.Write(body)
	if err != nil {
		return nil
	}
	return nil
}

//WriteHeader is the response head writer
func (bs *Context) WriteHeader(httpStatus int) {
	bs.Resp.WriteHeader(httpStatus)
}

//AddHeader is a function used to add header to a response
func (bs *Context) AddHeader(header string, value string) {
	bs.Resp.AddHeader(header, value)
}

//WriteError is a function used to write error into a response
func (bs *Context) WriteError(httpStatus int, err error) error {
	return bs.Resp.WriteError(httpStatus, err)
}

// WriteJSON used to write a JSON file into response
func (bs *Context) WriteJSON(value interface{}, contentType string) error {
	return bs.Resp.WriteJson(value, contentType)
}

// WriteHeaderAndJSON used to write head and JSON file in to response
func (bs *Context) WriteHeaderAndJSON(status int, value interface{}, contentType string) error {
	return bs.Resp.WriteHeaderAndJson(status, value, contentType)
}

//ReadEntity is request reader
func (bs *Context) ReadEntity(schema interface{}) (err error) {
	return bs.Req.ReadEntity(schema)
}

//ReadHeader is used to read header of request
func (bs *Context) ReadHeader(name string) string {
	return bs.Req.HeaderParameter(name)
}

//ReadPathParameter is used to read path parameter of a request
func (bs *Context) ReadPathParameter(name string) string {
	return bs.Req.PathParameter(name)
}

//ReadPathParameters used to read multiple path parameters of a request
func (bs *Context) ReadPathParameters() map[string]string {
	return bs.Req.PathParameters()
}

//ReadQueryParameter is used to read query parameter of a request
func (bs *Context) ReadQueryParameter(name string) string {
	return bs.Req.QueryParameter(name)
}

// ReadQueryEntity is used to read query parameters into a specified struct.
// The struct tag should be `form` like:
// type QueryRequest struct {
//     Name string `form:"name"`
//     Password string `form:"password"`
// }
func (bs *Context) ReadQueryEntity(schema interface{}) (err error) {
	return mapForm(schema, bs.Req.Request.URL.Query())
}

//ReadBodyParameter used to read body parameter of a request
func (bs *Context) ReadBodyParameter(name string) (string, error) {
	return bs.Req.BodyParameter(name)
}

//ReadRequest return a native net/http request
func (bs *Context) ReadRequest() *http.Request {
	return bs.Req.Request
}

//ReadRestfulRequest return a native  go-restful request
func (bs *Context) ReadRestfulRequest() *restful.Request {
	return bs.Req
}

//ReadResponseWriter return a native net/http ResponseWriter
func (bs *Context) ReadResponseWriter() http.ResponseWriter {
	return bs.Resp.ResponseWriter
}

//ReadRestfulResponse return a native go-restful Response
func (bs *Context) ReadRestfulResponse() *restful.Response {
	return bs.Resp
}
