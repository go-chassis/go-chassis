package handler_test

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateChain(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
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
