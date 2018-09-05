package client_test

import (
	"testing"

	_ "github.com/go-chassis/go-chassis/initiator"

	"github.com/go-chassis/go-chassis/client/highway"
	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
)

func TestInitError(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	m := make(map[string]model.Protocol)
	var pro model.Protocol
	pro.WorkerNumber = 1
	m["fake"] = pro

	config.GlobalDefinition.Cse.Protocols = m
	t.Log("get client func without installing its plugin")
	f, err := client.GetClientNewFunc("fake1")
	assert.Nil(t, f)
	assert.Error(t, err)
	t.Log("get Client without initializing")
	c, err := client.GetClient("fake1", "", "")
	assert.Error(t, err)
	assert.Nil(t, c)
}
func TestInit(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	m := make(map[string]model.Protocol)
	var pro model.Protocol
	pro.WorkerNumber = 1
	m["fake"] = pro

	config.GlobalDefinition.Cse.Protocols = m
	t.Log("get client func after installing its plugin")
	client.InstallPlugin("fake", highway.NewHighwayClient)
	f, err := client.GetClientNewFunc("fake")
	assert.NotNil(t, f)
	assert.NoError(t, err)
	t.Log("get Client after initializing")
	c, err := client.GetClient("fake", "service1", "127.0.0.1:9090")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	client.Close("fake", "service1", "127.0.0.1:9090")
	client.Close("notExist", "service1", "127.0.0.1:9090")

}

func BenchmarkGetClient(b *testing.B) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	m := make(map[string]model.Protocol)
	m["highway"] = model.Protocol{}
	config.GlobalDefinition.Cse.Protocols = m
	c, err := client.GetClient("highway", "", "")
	b.Log(c)
	if err != nil {
		b.Error(b, err)
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.GetClient("highway", "", "")
	}
}
