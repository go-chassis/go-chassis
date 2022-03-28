package handler_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/go-chassis/v2/control"
	_ "github.com/go-chassis/go-chassis/v2/control/servicecomb"
	"github.com/go-chassis/go-chassis/v2/core/config"
	chassisModel "github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/loadbalancer"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	mk "github.com/go-chassis/go-chassis/v2/core/registry/mock"
	_ "github.com/go-chassis/go-chassis/v2/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/v2/examples/schemas/helloworld"
	_ "github.com/go-chassis/go-chassis/v2/initiator"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/go-chassis/v2/pkg/util/fileutil"
	utiltags "github.com/go-chassis/go-chassis/v2/pkg/util/tags"
	"github.com/stretchr/testify/assert"
)

var callTimes = 0

type handler1 struct {
}

func (th *handler1) Name() string {
	return "loadbalancer"
}

func (th *handler1) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	callTimes++
	cb(&invocation.Response{})
}

type handler2 struct {
}

func (h *handler2) Name() string {
	return "handler2"
}

func (h *handler2) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	callTimes++
	r := &invocation.Response{
		Err: fmt.Errorf("A fake error from handler2"),
	}
	if callTimes < 10 {
		cb(r)
		return
	}
	cb(&invocation.Response{})
}

func TestTLSEndpointLBHandlerWithRetry(t *testing.T) {
	microContent := `---
servicecomb:
  service:
    name: Client
    version: 0.1`
	var yamlContent = `---
region:
  name: us-east
  availableZone: us-east-1
cse:
  loadbalance:
    strategy:
      name: RoundRobin
      sessionTimeoutInSeconds: 30
    retryEnabled: true
    retryOnNext: 0
    retryOnSame: 3
    backoff:
      kind: constant
      minMs: 200
      maxMs: 400
 `
	wd, _ := fileutil.GetWorkDir()
	os.Setenv("CHASSIS_HOME", wd)
	defer os.Unsetenv("CHASSIS_HOME")
	chassisConf := filepath.Join(wd, "conf")
	logConf := filepath.Join(wd, "log")
	err := os.MkdirAll(chassisConf, 0700)
	assert.NoError(t, err)
	err = os.MkdirAll(logConf, 0700)
	assert.NoError(t, err)
	chassisyaml := filepath.Join(chassisConf, "chassis.yaml")
	microserviceyaml := filepath.Join(chassisConf, "microservice.yaml")
	f1, err := os.Create(chassisyaml)
	assert.NoError(t, err)
	f2, err := os.Create(microserviceyaml)
	assert.NoError(t, err)
	_, err = io.WriteString(f1, yamlContent)
	assert.NoError(t, err)
	_, err = io.WriteString(f2, microContent)

	t.Log("testing load balance handler with retry")
	runtime.ServiceID = "selfServiceID"
	err = config.Init()
	assert.NoError(t, err)
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	err = control.Init(opts)
	assert.NoError(t, err)

	c := handler.Chain{}
	c.AddHandler(&handler.LBHandler{})
	c.AddHandler(&handler2{})

	var mss []*registry.MicroServiceInstance
	var ms1 = &registry.MicroServiceInstance{
		InstanceID: "instanceID",
		EndpointsMap: map[string]*registry.Endpoint{
			"rest": {
				Address:    "127.0.0.1:1234",
				SSLEnabled: true,
			},
		},
	}
	var ms2 = new(registry.MicroServiceInstance)
	ms2.EndpointsMap = map[string]*registry.Endpoint{
		"rest": {
			Address:    "127.0.0.1:1234",
			SSLEnabled: true,
		},
	}
	ms2.InstanceID = "ins2"
	mss = append(mss, ms1)
	mss = append(mss, ms2)

	testRegistryObj := new(mk.DiscoveryMock)
	registry.DefaultServiceDiscoveryService = testRegistryObj
	testRegistryObj.On("FindMicroServiceInstances", "selfServiceID", "appID", "service1", "1.0", "").
		Return(mss, nil)

	config.GlobalDefinition = &chassisModel.GlobalCfg{}
	config.GetLoadBalancing().Strategy = make(map[string]string)
	config.GetLoadBalancing().RetryEnabled = true
	config.GetLoadBalancing().RetryOnSame = 2
	loadbalancer.Enable(loadbalancer.StrategyRoundRobin)
	req, _ := rest.NewRequest("GET", "127.0.0.1", nil)
	req.Header.Set("Set-Cookie", "sessionid=100")
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             req,
		Strategy:         loadbalancer.StrategyRoundRobin,
		RouteTags:        utiltags.NewDefaultTag("1.0", "appID"),
		SourceServiceID:  runtime.ServiceID,
	}
	c.Next(i, func(r *invocation.Response) {
		assert.Error(t, r.Err)
	})

	var lbh = new(handler.LBHandler)
	str := lbh.Name()
	assert.Equal(t, "loadbalancer", str)
	assert.Equal(t, "rest", i.Protocol)
	assert.Equal(t, "127.0.0.1:1234", i.Endpoint)
	assert.True(t, i.SSLEnable)

}
func TestLBHandlerWithRetry(t *testing.T) {
	microContent := `---
servicecomb:
  service:
    name: Client
    version: 0.1`
	var yamlContent = `---
region:
  name: us-east
  availableZone: us-east-1
cse:
  loadbalance:
    strategy:
      name: RoundRobin
      sessionTimeoutInSeconds: 30
    retryEnabled: true
    retryOnNext: 0
    retryOnSame: 3
    backoff:
      kind: constant
      minMs: 200
      maxMs: 400
 `
	wd, _ := fileutil.GetWorkDir()
	os.Setenv("CHASSIS_HOME", wd)
	defer os.Unsetenv("CHASSIS_HOME")
	chassisConf := filepath.Join(wd, "conf")
	logConf := filepath.Join(wd, "log")
	err := os.MkdirAll(chassisConf, 0700)
	assert.NoError(t, err)
	err = os.MkdirAll(logConf, 0700)
	assert.NoError(t, err)
	chassisyaml := filepath.Join(chassisConf, "chassis.yaml")
	microserviceyaml := filepath.Join(chassisConf, "microservice.yaml")
	f1, err := os.Create(chassisyaml)
	assert.NoError(t, err)
	f2, err := os.Create(microserviceyaml)
	assert.NoError(t, err)
	_, err = io.WriteString(f1, yamlContent)
	assert.NoError(t, err)
	_, err = io.WriteString(f2, microContent)

	t.Log("testing load balance handler with retry")
	runtime.ServiceID = "selfServiceID"
	err = config.Init()
	assert.NoError(t, err)
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	err = control.Init(opts)
	assert.NoError(t, err)

	c := handler.Chain{}
	c.AddHandler(&handler.LBHandler{})
	c.AddHandler(&handler2{})

	var mss []*registry.MicroServiceInstance
	var ms1 = &registry.MicroServiceInstance{
		InstanceID: "instanceID",
		EndpointsMap: map[string]*registry.Endpoint{
			"rest": {
				Address: "127.0.0.1",
			},
		},
	}
	var ms2 = new(registry.MicroServiceInstance)
	ms2.EndpointsMap = map[string]*registry.Endpoint{
		"rest": {
			Address: "127.0.0.1",
		},
	}
	ms2.InstanceID = "ins2"
	mss = append(mss, ms1)
	mss = append(mss, ms2)

	testRegistryObj := new(mk.DiscoveryMock)
	registry.DefaultServiceDiscoveryService = testRegistryObj
	testRegistryObj.On("FindMicroServiceInstances", "selfServiceID", "appID", "service1", "1.0", "").
		Return(mss, nil)

	config.GlobalDefinition = &chassisModel.GlobalCfg{}
	config.GetLoadBalancing().Strategy = make(map[string]string)
	config.GetLoadBalancing().RetryEnabled = true
	config.GetLoadBalancing().RetryOnSame = 2
	loadbalancer.Enable(loadbalancer.StrategyRoundRobin)
	req, _ := rest.NewRequest("GET", "127.0.0.1", nil)
	req.Header.Set("Set-Cookie", "sessionid=100")
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             req,
		Strategy:         loadbalancer.StrategyRoundRobin,
		RouteTags:        utiltags.NewDefaultTag("1.0", "appID"),
		SourceServiceID:  runtime.ServiceID,
	}
	c.Next(i, func(r *invocation.Response) {
		assert.Error(t, r.Err)
	})

	var lbh = new(handler.LBHandler)
	str := lbh.Name()
	assert.Equal(t, "loadbalancer", str)
	assert.Equal(t, "rest", i.Protocol)
	assert.Equal(t, "127.0.0.1", i.Endpoint)
	assert.False(t, i.SSLEnable)
}
func TestTLSLBHandlerWithNoRetry(t *testing.T) {
	microContent := `---
#微服务的私有属性
servicecomb:
  service:
    name: Client
    version: 0.1`
	var yamlContent = `---
region:
  name: us-east
  availableZone: us-east-1
cse:
  loadbalance:
    strategy:
      name: RoundRobin
    retryEnabled: false
 `
	wd, _ := fileutil.GetWorkDir()
	os.Setenv("CHASSIS_HOME", wd)
	defer os.Unsetenv("CHASSIS_HOME")
	chassisConf := filepath.Join(wd, "conf")
	logConf := filepath.Join(wd, "log")
	err := os.MkdirAll(chassisConf, 0700)
	assert.NoError(t, err)
	err = os.MkdirAll(logConf, 0700)
	assert.NoError(t, err)
	chassisyaml := filepath.Join(chassisConf, "load_balancing.yaml")
	microserviceyaml := filepath.Join(chassisConf, "microservice.yaml")
	f1, err := os.Create(chassisyaml)
	assert.NoError(t, err)
	f2, err := os.Create(microserviceyaml)
	assert.NoError(t, err)
	_, err = io.WriteString(f1, yamlContent)
	assert.NoError(t, err)
	_, err = io.WriteString(f2, microContent)

	t.Log("testing load balance handler with retry")
	runtime.ServiceID = "selfServiceID"
	err = config.Init()
	assert.NoError(t, err)
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	err = control.Init(opts)
	assert.NoError(t, err)
	c := handler.Chain{}
	c.AddHandler(&handler.LBHandler{})
	c.AddHandler(&handler1{})
	var mss []*registry.MicroServiceInstance
	var ms1 = new(registry.MicroServiceInstance)
	var ms2 = new(registry.MicroServiceInstance)
	var mp = make(map[string]*registry.Endpoint)
	mp["any-tls"] = &registry.Endpoint{Address: "127.0.0.1", SSLEnabled: true}
	ms1.EndpointsMap = mp
	ms1.InstanceID = "ins1"
	ms2.EndpointsMap = mp
	ms2.InstanceID = "ins2"
	testRegistryObj := new(mk.DiscoveryMock)

	loadbalancer.Enable(loadbalancer.StrategyRoundRobin)
	registry.DefaultServiceDiscoveryService = testRegistryObj
	mss = append(mss, ms1)
	mss = append(mss, ms2)

	testRegistryObj.On("FindMicroServiceInstances",
		"selfServiceID", "appID", "service2", "1.0", "").Return(mss, nil)

	i := &invocation.Invocation{
		MicroServiceName: "service2",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
		SourceServiceID:  "selfServiceID",
		RouteTags:        utiltags.NewDefaultTag("1.0", "appID"),
	}
	t.Run("invocation without strategy", func(t *testing.T) {
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		i.HandlerIndex = 0
	})

	i.Strategy = loadbalancer.StrategyRoundRobin
	c.Next(i, func(r *invocation.Response) {
		assert.NoError(t, r.Err)
	})

	var lbh = new(handler.LBHandler)
	str := lbh.Name()
	assert.Equal(t, "loadbalancer", str)
	assert.Equal(t, "any-tls", i.Protocol)
	assert.Equal(t, "127.0.0.1", i.Endpoint)
	assert.True(t, i.SSLEnable)
}
func TestLBHandlerWithNoRetry(t *testing.T) {
	microContent := `---
#微服务的私有属性
servicecomb:
  service:
	  name: Client
	  version: 0.1`
	var yamlContent = `---
region:
  name: us-east
  availableZone: us-east-1
cse:
  loadbalance:
    strategy:
      name: RoundRobin
    retryEnabled: false
 `
	wd, _ := fileutil.GetWorkDir()
	os.Setenv("CHASSIS_HOME", wd)
	defer os.Unsetenv("CHASSIS_HOME")
	chassisConf := filepath.Join(wd, "conf")
	logConf := filepath.Join(wd, "log")
	err := os.MkdirAll(chassisConf, 0700)
	assert.NoError(t, err)
	err = os.MkdirAll(logConf, 0700)
	assert.NoError(t, err)
	chassisyaml := filepath.Join(chassisConf, "chassis.yaml")
	microserviceyaml := filepath.Join(chassisConf, "microservice.yaml")
	f1, err := os.Create(chassisyaml)
	assert.NoError(t, err)
	f2, err := os.Create(microserviceyaml)
	assert.NoError(t, err)
	_, err = io.WriteString(f1, yamlContent)
	assert.NoError(t, err)
	_, err = io.WriteString(f2, microContent)

	t.Log("testing load balance handler with retry")
	runtime.ServiceID = "selfServiceID"
	err = config.Init()
	assert.NoError(t, err)
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	err = control.Init(opts)
	assert.NoError(t, err)
	c := handler.Chain{}
	c.AddHandler(&handler.LBHandler{})
	c.AddHandler(&handler1{})
	var mss []*registry.MicroServiceInstance
	var ms1 = new(registry.MicroServiceInstance)
	var ms2 = new(registry.MicroServiceInstance)
	var mp = make(map[string]*registry.Endpoint)
	mp["any"] = &registry.Endpoint{Address: "127.0.0.1:123", SSLEnabled: false}
	ms1.EndpointsMap = mp
	ms1.InstanceID = "ins1"
	ms2.EndpointsMap = mp
	ms2.InstanceID = "ins2"
	testRegistryObj := new(mk.DiscoveryMock)

	loadbalancer.Enable(loadbalancer.StrategyRoundRobin)
	registry.DefaultServiceDiscoveryService = testRegistryObj
	t.Run("select service1 without instances", func(t *testing.T) {
		testRegistryObj.On("FindMicroServiceInstances",
			"selfServiceID", "appID", "service1", "1.0", "").
			Return(make([]*registry.MicroServiceInstance, 0), nil)

		i := &invocation.Invocation{
			MicroServiceName: "service1",
			SourceServiceID:  "selfServiceID",
			RouteTags:        utiltags.NewDefaultTag("1.0", "appID"),
		}
		c.Next(i, func(r *invocation.Response) {
			assert.Error(t, r.Err)
		})
		i.HandlerIndex = 0
	})
	t.Run("select service3 without instances eps", func(t *testing.T) {
		testRegistryObj.On("FindMicroServiceInstances",
			"selfServiceID", "appID", "service3", "1.0", "").
			Return(
				[]*registry.MicroServiceInstance{{EndpointsMap: make(map[string]*registry.Endpoint, 0)}}, nil)

		i := &invocation.Invocation{
			MicroServiceName: "service3",
			SourceServiceID:  "selfServiceID",
			RouteTags:        utiltags.NewDefaultTag("1.0", "appID"),
		}
		c.Next(i, func(r *invocation.Response) {
			assert.Error(t, r.Err)
		})
		i.HandlerIndex = 0
	})
	mss = append(mss, ms1)
	mss = append(mss, ms2)

	testRegistryObj.On("FindMicroServiceInstances",
		"selfServiceID", "appID", "service2", "1.0", "").Return(mss, nil)

	i := &invocation.Invocation{
		MicroServiceName: "service2",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
		SourceServiceID:  "selfServiceID",
		RouteTags:        utiltags.NewDefaultTag("1.0", "appID"),
	}
	t.Run("invocation without strategy", func(t *testing.T) {
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		i.HandlerIndex = 0
	})

	i.Strategy = loadbalancer.StrategyRoundRobin
	c.Next(i, func(r *invocation.Response) {
		assert.NoError(t, r.Err)
	})

	var lbh = new(handler.LBHandler)
	str := lbh.Name()
	assert.Equal(t, "loadbalancer", str)
	t.Log(i.Protocol)
	t.Log(i.Endpoint)
}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
func BenchmarkLBHandler_Handle(b *testing.B) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "client"))
	defer os.Unsetenv("CHASSIS_HOME")
	config.Init()
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	control.Init(opts)
	registry.Enable()
	registry.DoRegister()
	loadbalancer.Enable(archaius.GetString("cse.loadbalance.strategy.name", ""))
	testData1 := []*registry.MicroService{
		{
			ServiceName: "test2",
			AppID:       "default",
			Version:     "1.0",
			Status:      "UP",
		},
	}
	testData2 := []*registry.MicroServiceInstance{
		{
			HostName: "test1",
			Status:   "UP",
			EndpointsMap: map[string]*registry.Endpoint{
				"highway": {
					Address: "10.0.0.4:1234",
				},
			},
		},
		{
			HostName: "test2",
			Status:   "UP",
			EndpointsMap: map[string]*registry.Endpoint{
				"highway": {
					Address: "10.0.0.3:1234",
				},
			},
		},
	}
	sid, _, _ := registry.DefaultRegistrator.RegisterServiceAndInstance(testData1[0], testData2[0])
	_, _, _ = registry.DefaultRegistrator.RegisterServiceAndInstance(testData1[0], testData2[1])
	c := handler.Chain{}
	c.AddHandler(&handler.LBHandler{})
	c.AddHandler(&handler1{})
	runtime.ServiceID = sid
	iv := &invocation.Invocation{
		MicroServiceName: "test2",
		Protocol:         "highway",
		Strategy:         loadbalancer.StrategyRoundRobin,
		SourceServiceID:  runtime.ServiceID,
		RouteTags:        utiltags.NewDefaultTag("1.0", "default"),
	}

	b.Log(runtime.ServiceID)
	time.Sleep(1 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Next(iv, func(r *invocation.Response) {
		})
		iv.HandlerIndex = 0
	}

}
