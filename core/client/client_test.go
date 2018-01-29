package client_test

import (
	"testing"
	"time"

	"github.com/ServiceComb/go-chassis/client/highway"
	"github.com/ServiceComb/go-chassis/core/client"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	clientOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/codec"
	_ "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport/tcp"
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
	c, err := client.GetClient("fake1", "")
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
	client.InstallPlugin("fake", tcp.NewHighwayClient)
	f, err := client.GetClientNewFunc("fake")
	assert.NotNil(t, f)
	assert.NoError(t, err)
	t.Log("get Client after initializing")
	c, err := client.GetClient("fake", "")
	assert.NoError(t, err)
	assert.NotNil(t, c)
}
func TestOptions(t *testing.T) {
	t.Log("sets various parameter to option")
	tduration := time.Second * 2
	c := make(map[string]codec.Codec)
	c["fake"] = codec.NewJSONCodec()

	cp := clientOption.WithCodecs(c)

	var cstruct *clientOption.Options = new(clientOption.Options)
	var copstruct *clientOption.CallOptions = new(clientOption.CallOptions)

	cstruct.ContentType = "fakectype"
	cp(cstruct)
	assert.Equal(t, cstruct.ContentType, "fakectype")

	psize := clientOption.PoolSize(2)
	psize(cstruct)
	assert.Equal(t, cstruct.PoolSize, 2)

	op1 := clientOption.PoolTTL(tduration)
	op1(cstruct)
	assert.Equal(t, cstruct.PoolTTL, tduration)

	op2 := clientOption.Retries(3)
	op2(cstruct)
	assert.Equal(t, cstruct.CallOptions.Retries, 3)

	op3 := clientOption.RequestTimeout(tduration)
	op3(cstruct)
	assert.Equal(t, cstruct.CallOptions.RequestTimeout, tduration)

	op4 := clientOption.DialTimeout(tduration)
	op4(cstruct)
	assert.Equal(t, cstruct.CallOptions.DialTimeout, tduration)

	cop := clientOption.WithRetries(3)
	cop(copstruct)
	assert.Equal(t, copstruct.Retries, 3)

	cop1 := clientOption.WithRequestTimeout(tduration)
	cop1(copstruct)
	assert.Equal(t, copstruct.RequestTimeout, tduration)

	cop2 := clientOption.WithDialTimeout(tduration)
	cop2(copstruct)
	assert.Equal(t, copstruct.DialTimeout, tduration)

	cop3 := clientOption.WithContentType("fakect")
	cop3(copstruct)
	assert.Equal(t, copstruct.ContentType, "fakect")

	cop4 := clientOption.WithUrlPath("fakeurl")
	cop4(copstruct)
	assert.Equal(t, copstruct.UrlPath, "fakeurl")

	cop5 := clientOption.WithMethodType("fakemt")
	cop5(copstruct)
	assert.Equal(t, copstruct.MethodType, "fakemt")

}
