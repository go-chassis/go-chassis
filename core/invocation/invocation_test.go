package invocation_test

import (
	"context"
	"testing"

	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	var inv = new(invocation.Invocation)
	inv.Endpoint = "1.2.3.4"
}
func TestNew(t *testing.T) {
	i := invocation.New(nil)
	assert.NotNil(t, i.Ctx)
}
func TestInvocation_Headers(t *testing.T) {
	inv := invocation.New(context.TODO())
	inv.SetMetadata("a", "1")
	inv.SetHeader("asd", "123")
	assert.Equal(t, "123", inv.Headers()["asd"])
	assert.Equal(t, "1", inv.Metadata["a"])
}

/*
type handler1 struct {

}

	func (h *handler1) Handle(c *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
		//log.Println("执行1")
		c.Next(i, func(r *invocation.Response) error {
			//log.Println("回调到起始")
			//log.Println(r)
			r.Err = errors.New("wrong")
			return cb(r)
		})
	}

	func (h *handler1) Name() string {
		return "test"
	}

type handler2 struct {

}

	func (h *handler2) Name() string {
		return "test"
	}

	func (h *handler2) Handle(c *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
		//log.Println("执行2")
		i.Endpoint = "test"
		c.Next(i, func(r *invocation.Response) error {
			//log.Println(r)
			r.Status = 2
			//log.Println("回调到1")

			return cb(r)
		})

}

type transportHandler struct {

}

	func (h *transportHandler) Name() string {
		return "test"
	}

	func (h *transportHandler) Handle(c *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
		//log.Println("fake transport handler")
		r := &invocation.Response{Status:200, }
		cb(r)
	}

	func TestChain(t *testing.T) {
		p := os.Getenv("GOPATH")
		os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "communication", "client"))
		defer os.Unsetenv("CHASSIS_HOME")
		lager.Initialize()
		config.Init()
		c := &handler.Chain{}
		i := &invocation.Invocation{}
		c.Handlers = append(c.Handlers, &handler1{}, &handler2{}, &transportHandler{})
		var err error
		c.Next(i, func(r *invocation.Response) error {
			if r != nil {
				err = r.Err
				assert.Equal(t, "test", i.Endpoint)
				return r.Err
			}
			return nil
		})
		log.Println("err " + err.Error())
		//assert.Equal(t, 123, 123, "they should be equal")
		//
		//// assert inequality
		//assert.NotEqual(t, 123, 456, "they should not be equal")
		//
		//// assert for nil (good for errors)
		//assert.Nil(t, object)
		//
		//// assert for not nil (good when you expect something)
		//if assert.NotNil(t, object) {
		//
		//	// now we know that object isn't nil, we are safe to make
		//	// further assertions without causing any errors
		//	assert.Equal(t, "Something", object.Value)
		//
		//}

}

	func BenchmarkChainNext(b *testing.B) {
		p := os.Getenv("GOPATH")
		os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "communication", "client"))
		defer os.Unsetenv("CHASSIS_HOME")
		lager.Initialize()
		config.Init()
		c := &handler.Chain{}
		i := &invocation.Invocation{}
		c.Handlers = append(c.Handlers, &handler1{}, &handler2{}, &transportHandler{})
		var err error
		for j := 0; j < b.N; j++ {
			c.Next(i, func(r *invocation.Response) error {
				if r != nil {
					err = r.Err
					return r.Err
				}
				return nil
			})
			c.HandlerIndex=0
		}

}
*/
func BenchmarkConextMap(b *testing.B) {
	ctx := context.WithValue(context.Background(), "test", map[string]string{
		"test": "1",
	})
	for j := 0; j < b.N; j++ {
		m := ctx.Value("test").(map[string]string)
		_, _ = m["test"]

	}

}
func BenchmarkMap(b *testing.B) {
	m := make(map[string]string)
	m["test"] = "1"
	for j := 0; j < b.N; j++ {
		_, _ = m["test"]

	}

}
