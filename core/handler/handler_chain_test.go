package handler_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
func TestCreateChain(t *testing.T) {
	t.Log("testing creation of chain with various service type,chain name and handlers")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	e := handler.RegisterHandler("fake", newProviderHandler)
	assert.NoError(t, e)
	c, err := handler.CreateChain("abc", "fake")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c, err = handler.CreateChain(common.Consumer, "fake")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c, err = handler.CreateChain(common.Provider, "fake")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c, err = handler.CreateChain(common.Consumer, "fake")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	chopt := handler.WithChainName("chainName")
	var ch *handler.ChainOptions = new(handler.ChainOptions)
	chopt(ch)
	assert.Equal(t, "chainName", ch.Name)
}
func BenchmarkChain_Next(b *testing.B) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "client"))
	defer os.Unsetenv("CHASSIS_HOME")
	config.GlobalDefinition = &model.GlobalCfg{}
	config.Init()
	iv := &invocation.Invocation{}
	handler.RegisterHandler("f1", createBizkeeperFakeHandler)
	handler.RegisterHandler("f2", createBizkeeperFakeHandler)
	handler.RegisterHandler("f3", createBizkeeperFakeHandler)
	if err := handler.CreateChains(common.Consumer, map[string]string{
		"default": "f1,f2,f3,f1,f2,f3,f1,f2,f3,f1,f2,f3,f1,f2",
	}); err != nil {
		b.Fatal(err)
	}

	c, err := handler.GetChain(common.Consumer, "default")
	if err != nil {
		b.Fatal(err)
	}
	log.Println("----------------------------------------------------")
	log.Println(c)

	for i := 0; i < b.N; i++ {
		c, _ = handler.GetChain(common.Consumer, "default")
		c.Next(iv, func(r *invocation.Response) {
		})
		iv.HandlerIndex = 0
	}
}
