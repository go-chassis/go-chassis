package client_test

import (
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"testing"

	"github.com/go-chassis/go-chassis/v2/client/rest"
	_ "github.com/go-chassis/go-chassis/v2/initiator"

	"time"

	"github.com/go-chassis/go-chassis/v2/core/client"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/stretchr/testify/assert"
)

func init() {

	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
	config.HystrixConfig = &model.HystrixConfigWrapper{
		HystrixConfig: &model.HystrixConfig{
			IsolationProperties: &model.IsolationWrapper{
				Consumer: &model.IsolationSpec{},
			},
		},
	}
}
func TestGetFailureMap(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.Transport.Failure = map[string]string{
		"rest": "http_500,http:502",
	}

	t.Run("get failed map about protocol ",
		func(t *testing.T) {
			m := client.GetFailureMap("rest")
			assert.NotEmpty(t, m)
			assert.True(t, m["http_500"])
			assert.False(t, m["http_540"])
			m = client.GetFailureMap("rpc")
			assert.Empty(t, m)

		})
}
func TestGetMaxIdleCon(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.Transport.MaxIdlCons = map[string]int{
		"rest": 1,
	}

	t.Run("get max idle by user customize or default",
		func(t *testing.T) {
			n := client.GetMaxIdleCon("rest")
			assert.Equal(t, 1, n)
			config.GlobalDefinition.ServiceComb.Transport.MaxIdlCons = map[string]int{}
			n = client.GetMaxIdleCon("rest")
			assert.Equal(t, 512, n)
		})
}
func TestInitError(t *testing.T) {
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	m := make(map[string]model.Protocol)
	var pro model.Protocol
	pro.WorkerNumber = 1
	m["fake"] = pro

	config.GlobalDefinition.ServiceComb.Protocols = m

	t.Run("get client func without installing its plugin",
		func(t *testing.T) {
			f, err := client.GetClientNewFunc("fake1")
			assert.Nil(t, f)
			assert.Error(t, err)
		})
	t.Run("get Client without initializing",
		func(t *testing.T) {
			i := &invocation.Invocation{Protocol: "fake1"}
			c, err := client.GetClient(i)
			assert.Error(t, err)
			assert.Nil(t, c)
		})
	t.Run("get Client for rest",
		func(t *testing.T) {
			i := &invocation.Invocation{Protocol: "rest"}
			c, err := client.GetClient(i)
			assert.Nil(t, err)
			assert.NotNil(t, c)
		})
}
func TestInit(t *testing.T) {
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	m := make(map[string]model.Protocol)
	var pro model.Protocol
	pro.WorkerNumber = 1
	m["fake"] = pro

	config.GlobalDefinition.ServiceComb.Protocols = m
	t.Run("get client func after installing its plugin",
		func(t *testing.T) {
			client.InstallPlugin("fake", rest.NewRestClient)
			f, err := client.GetClientNewFunc("fake")
			assert.NotNil(t, f)
			assert.NoError(t, err)
		})
	t.Run("get Client after initializing",
		func(t *testing.T) {
			i := &invocation.Invocation{Protocol: "fake", MicroServiceName: "service1", Endpoint: "127.0.0.1:9090"}
			c, err := client.GetClient(i)
			assert.NoError(t, err)
			assert.NotNil(t, c)
		})
	t.Run("close client , client exist ot not exist",
		func(t *testing.T) {
			err := client.Close("fake", "service1", "127.0.0.1:9090")
			assert.Nil(t, err)
			err = client.Close("notExist", "service1", "127.0.0.1:9090")
			assert.NotNil(t, err)
			assert.Equal(t, err, client.ErrClientNotExist)

		})
}

func BenchmarkGetClient(b *testing.B) {
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	m := make(map[string]model.Protocol)
	m["highway"] = model.Protocol{}
	m["rest"] = model.Protocol{}
	m["grpc"] = model.Protocol{}
	config.GlobalDefinition.ServiceComb.Protocols = m

	i := &invocation.Invocation{Protocol: "highway"}
	c, err := client.GetClient(i)
	b.Log(c)
	if err != nil {
		b.Error(b, err)
	}

	b.Run("benchmark get highway client , no support by default",
		func(b *testing.B) {
			i := &invocation.Invocation{Protocol: "highway"}
			c, err := client.GetClient(i)
			assert.Nil(b, c)
			assert.NotNil(b, err)
		})
	b.Run("benchmark get grpc client , no support by default",
		func(b *testing.B) {
			i := &invocation.Invocation{Protocol: "grpc"}
			c, err := client.GetClient(i)
			assert.Nil(b, c)
			assert.NotNil(b, err)
		})
	b.Run("benchmark get rest client",
		func(b *testing.B) {
			i := &invocation.Invocation{Protocol: "rest"}
			c, err := client.GetClient(i)
			assert.NotNil(b, c)
			assert.Nil(b, err)
		})
}
func TestSetTimeoutToClientCache(t *testing.T) {
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	m := make(map[string]model.Protocol)
	m["rest"] = model.Protocol{}
	config.GlobalDefinition.ServiceComb.Protocols = m
	i := &invocation.Invocation{Protocol: "rest", MicroServiceName: "rest_server"}
	c, err := client.GetClient(i)
	assert.NotEmpty(t, c)
	assert.Nil(t, err)
	i2 := &invocation.Invocation{Protocol: "rest", MicroServiceName: "rest_server1"}
	c, err = client.GetClient(i2)
	assert.NotEmpty(t, c)
	assert.Nil(t, err)

	spec := &model.IsolationWrapper{
		Consumer: &model.IsolationSpec{
			TimeoutInMilliseconds: config.DefaultTimeout,
		},
	}
	t.Run("set global timeout will set all service",
		func(t *testing.T) {
			spec.Consumer.TimeoutInMilliseconds = 20
			client.SetTimeoutToClientCache(spec)
			i := &invocation.Invocation{Protocol: "rest", MicroServiceName: "rest_server"}
			c, err := client.GetClient(i)
			assert.Nil(t, err)
			assert.Equal(t, time.Duration(20)*time.Millisecond, c.GetOptions().Timeout)
			i2 := &invocation.Invocation{Protocol: "rest", MicroServiceName: "rest_server1"}
			c, err = client.GetClient(i2)
			assert.Nil(t, err)
			assert.Equal(t, time.Duration(20)*time.Millisecond, c.GetOptions().Timeout)
		})

	t.Run("set service timeout will set one or more for service",
		func(t *testing.T) {
			spec.Consumer.AnyService = map[string]model.IsolationSpec{
				"rest_server": {
					TimeoutInMilliseconds: 10,
				},
			}
			client.SetTimeoutToClientCache(spec)
			i := &invocation.Invocation{Protocol: "rest", MicroServiceName: "rest_server"}
			c, err := client.GetClient(i)
			assert.Nil(t, err)
			assert.Equal(t, time.Duration(10)*time.Millisecond, c.GetOptions().Timeout)
			i2 := &invocation.Invocation{Protocol: "rest", MicroServiceName: "rest_server1"}
			c, err = client.GetClient(i2)
			assert.Nil(t, err)
			assert.Equal(t, time.Duration(20)*time.Millisecond, c.GetOptions().Timeout)
		})

}
