package handler_test

import (
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/provider"
	"github.com/go-chassis/go-chassis/v2/pkg/util/fileutil"
	"github.com/go-chassis/openlog"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func prepareConfDir(t *testing.T) string {
	wd, _ := fileutil.GetWorkDir()
	os.Setenv("CHASSIS_HOME", wd)
	defer os.Unsetenv("CHASSIS_HOME")
	chassisConf := filepath.Join(wd, "conf")
	logConf := filepath.Join(wd, "log")
	err := os.MkdirAll(chassisConf, 0700)
	assert.NoError(t, err)
	err = os.MkdirAll(logConf, 0700)
	assert.NoError(t, err)
	return chassisConf
}
func prepareTestFile(t *testing.T, confDir, file, content string) {
	fullPath := filepath.Join(confDir, file)
	err := os.Remove(fullPath)
	f, err := os.Create(fullPath)
	assert.NoError(t, err)
	_, err = io.WriteString(f, content)
	assert.NoError(t, err)
}

type ProviderHandler struct {
}

func (ph *ProviderHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	p, err := provider.GetProvider(i.MicroServiceName)
	if err != nil {
		openlog.Error("GetProvider failed." + err.Error())
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
	e := handler.RegisterHandler("fake3", newProviderHandler)
	assert.NoError(t, e)
	t.Log("testing registration of a custom handler against a name which is already registered")
	e = handler.RegisterHandler("fake3", newProviderHandler)
	assert.Equal(t, handler.ErrDuplicatedHandler, e)

	e = handler.RegisterHandler(handler.Transport, newProviderHandler)
	assert.Error(t, e)
}

func TestCreateHandler(t *testing.T) {
	t.Log("testing creation of handler")
	config.Init()
	e := handler.RegisterHandler("fake2", newProviderHandler)
	assert.NoError(t, e)
	e = handler.RegisterHandler(handler.Transport, newProviderHandler)
	assert.Error(t, e)
	_, err := handler.CreateHandler("123")
	assert.Error(t, err)
	handler, err := handler.CreateHandler("fake3")
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
func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
func TestGetChain(t *testing.T) {
	t.Log("getting chain of various service type and chain name")

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = map[string]string{
		"default": "bizkeeper-fake,loadbalancer-fake",
		"custom":  "bizkeeper-fake",
	}
	config.GlobalDefinition.ServiceComb.Handler.Chain.Provider = map[string]string{
		"default": "bizkeeper-fake,loadbalancer-fake",
	}
	handler.RegisterHandler(BIZKEEPERFAKE, createBizkeeperFakeHandler)
	handler.RegisterHandler("loadbalancer-fake", createBizkeeperFakeHandler)
	handler.CreateChains(common.Provider, config.GlobalDefinition.ServiceComb.Handler.Chain.Provider)
	handler.CreateChains(common.Consumer, config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer)
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
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = map[string]string{
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
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "client"))
	config.GlobalDefinition = &model.GlobalCfg{}
	config.Init()
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = map[string]string{
		"default": "bizkeeper-fake,loadbalancer-fake",
		"custom":  "bizkeeper-fake",
	}
	config.GlobalDefinition.ServiceComb.Handler.Chain.Provider = map[string]string{
		"default": "bizkeeper-fake,loadbalancer-fake",
	}

	for i := 0; i < b.N; i++ {
		_, _ = handler.GetChain(common.Consumer, "default")
	}

}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
