package handler_test

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/provider"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

type ProviderHandler struct {
}

func (ph *ProviderHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	p, err := provider.GetProvider(i.MicroServiceName)
	if err != nil {
		lager.Logger.Error("GetProvider failed.", err)
	}
	p.Invoke(i)
}

func (ph *ProviderHandler) Name() string {
	return "test"
}
func newProviderHandler() handler.Handler {
	return &ProviderHandler{}
}
func TestRegisterHandlerFunc(t *testing.T) {
	t.Log("testing registration of a custom  handler")
	config.Init()
	e := handler.RegisterHandler("fake", newProviderHandler)
	assert.NoError(t, e)
	t.Log("testing registration of a custom handler against a name which is already registered")
	e = handler.RegisterHandler(handler.Transport, newProviderHandler)
	assert.Error(t, e)
}

func TestCreateHandler(t *testing.T) {
	t.Log("testing creation of handler")
	config.Init()
	e := handler.RegisterHandler("fake", newProviderHandler)
	assert.NoError(t, e)
	e = handler.RegisterHandler(handler.Transport, newProviderHandler)
	assert.Error(t, e)
	_, err := handler.CreateHandler("123")
	assert.Error(t, err)
	handler, err := handler.CreateHandler("fake")
	assert.NoError(t, err)
	t.Log(handler)
}

var BIZKEEPERFAKE = "bizkeeper-fake"

type BizkeeperFakeHandler struct{}

func (bizkeeperfhandler *BizkeeperFakeHandler) Name() string {
	return BIZKEEPERFAKE
}

func (bizkeeperfhandler *BizkeeperFakeHandler) Handle(*handler.Chain, *invocation.Invocation, invocation.ResponseCallBack) {
	return
}

func createBizkeeperFakeHandler() handler.Handler {
	return &BizkeeperFakeHandler{}
}

func TestGetChain(t *testing.T) {
	t.Log("getting chain of various service type and chain name")
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = map[string]string{
		"default": "bizkeeper-fake,loadbalance-fake",
		"custom":  "bizkeeper-fake",
	}
	config.GlobalDefinition.Cse.Handler.Chain.Provider = map[string]string{
		"default": "bizkeeper-fake,loadbalance-fake",
	}
	handler.RegisterHandler(BIZKEEPERFAKE, createBizkeeperFakeHandler)
	handler.RegisterHandler("loadbalance-fake", createBizkeeperFakeHandler)
	handler.CreateChains(common.Provider, config.GlobalDefinition.Cse.Handler.Chain.Provider)
	handler.CreateChains(common.Consumer, config.GlobalDefinition.Cse.Handler.Chain.Consumer)
	c, err := handler.GetChain(common.Consumer, "custom")
	assert.NoError(t, err)
	assert.Equal(t, "custom", c.Name)
	assert.Equal(t, common.Consumer, c.ServiceType)
	t.Log(c.Handlers[0])

	c, err = handler.GetChain(common.Provider, "")
	assert.NoError(t, err)
	assert.Equal(t, "default", c.Name)
	assert.Equal(t, common.Provider, c.ServiceType)
	t.Log(c.Handlers[0])

	t.Log("handler name为空")
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = map[string]string{
		"default": "",
		"custom":  ",",
	}
	assert.NoError(t, err)
	c, err = handler.GetChain(common.Consumer, "custom")
	assert.NoError(t, err)
	assert.Equal(t, "custom", c.Name)
	assert.Equal(t, common.Consumer, c.ServiceType)
	c, err = handler.GetChain(common.Consumer, "default")
	assert.NoError(t, err)
	assert.Equal(t, "default", c.Name)
	assert.Equal(t, common.Consumer, c.ServiceType)
	t.Log(c.Handlers[0])

	ch, err := handler.GetChain("serviceType", "name")
	assert.Nil(t, ch)
	assert.Error(t, err)

}
func BenchmarkPool_GetChain(b *testing.B) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "client"))
	config.GlobalDefinition = &model.GlobalCfg{}
	config.Init()
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = map[string]string{
		"default": "bizkeeper-fake,loadbalance-fake",
		"custom":  "bizkeeper-fake",
	}
	config.GlobalDefinition.Cse.Handler.Chain.Provider = map[string]string{
		"default": "bizkeeper-fake,loadbalance-fake",
	}
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	for i := 0; i < b.N; i++ {
		_, _ = handler.GetChain(common.Consumer, "default")
	}

}
