package core_test

import (
	//"fmt"
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core"
	"github.com/ServiceComb/go-chassis/core/config"
	_ "github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/util/metadata"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func initenv() {
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	//handler.Init()
}

func TestRPCInvoker_InvokeFailinChainInit(t *testing.T) {
	initenv()
	config.GlobalDefinition = &model.GlobalCfg{}
	invoker := core.NewRPCInvoker(core.ChainName(""))
	replyOne := &helloworld.HelloReply{}
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"X-User": "tianxiaoliang",
	})

	config.GlobalDefinition.Cse.References = make(map[string]model.ReferencesStruct)
	version := model.ReferencesStruct{Version: ""}
	config.GlobalDefinition.Cse.References["Server"] = version
	err := invoker.Invoke(ctx, "Server", "HelloServer", "SayHello", &helloworld.HelloRequest{Name: "Peter"}, replyOne, core.WithAppID("0.2"),
		core.WithMetadata(nil), core.WithStrategy(""), core.StreamingRequest())
	assert.Error(t, err)
}
func TestRestInvoker_ContextDo(t *testing.T) {
	initenv()
	restinvoker := core.NewRestInvoker()
	req, _ := rest.NewRequest("GET", "cse://Server/sayhello/myidtest")
	//use the invoker like http client.
	_, err := restinvoker.ContextDo(context.TODO(), req, core.WithContentType("application/json"), core.WithEndpoint("0.0.0.0"), core.WithProtocol("rest"), core.WithFilters(nil))
	assert.Error(t, err)
}

func TestOptions(t *testing.T) {
	opt := core.InvokeOptions{Version: "0.1"}
	option := core.DefaultCallOptions(opt)
	assert.NotEmpty(t, option)

	inv := core.StreamingRequest()
	assert.NotEmpty(t, inv)

	inv = core.WithEndpoint("0.0.0.0")
	assert.NotEmpty(t, inv)

	inv = core.WithVersion("0.0")
	assert.NotEmpty(t, inv)

	inv = core.WithProtocol("0.0")
	assert.NotEmpty(t, inv)

	inv = core.WithFilters(nil)
	assert.NotEmpty(t, inv)

	inv = core.WithStrategy("")
	assert.NotEmpty(t, inv)

	inv = core.WithMetadata(nil)
	assert.NotEmpty(t, inv)
}
