package handler_test

import (
	"errors"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var yamlContent = `---
region:
  name: us-east
  availableZone: us-east-1
#APPLICATION_ID: CSE optional

cse:
#  credentials:
#    accessKey:
#    secretKey:
  governance:
    Consumer:
      _global:
        policy:
          fault:
            protocols:
              rest:
                delay:
                  fixedDelay: 5
                  percent: 100
                abort:
                  httpStatus: 451
                  percent: 100
  flowcontrol:
    Consumer:
      qps:
        enabled: true
        limit:
          Server.HelloServer: 10
  loadbalance:
    strategy:
      name: RoundRobin
      sessionTimeoutInSeconds: 30
    retryEnabled: false
    retryOnNext: 2
    retryOnSame: 3
    backoff:
      kind: constant
      minMs: 200
      maxMs: 400
  service:
    registry:
      #disabled: false           optional:禁用注册发现选项，默认开始注册发现
      type: servicecenter           #optional:可选zookeeper/servicecenter，zookeeper供中软使用，不配置的情况下默认为servicecenter
      scope: full                   #optional:scope不为full时，只允许在本app间访问，不允许跨app访问；为full就是注册时允许跨app，并且发现本租户全部微服务
      address: https://cse.cn-north-1.myhwclouds.com:443
      #register: manual          optional：register不配置时默认为自动注册，可选参数有自动注册auto和手动注册manual
      refeshInterval : 1
      watch: true
  config:
    client:
      #serverUri: 10.74.175.142:30103 #ip of config center
      serverUri: https://cse.cn-north-1.myhwclouds.com:443
      tenantName:  default #This configuration is for local environment, for paas platform there is a auth plugin for authentication. If dont provide the tenant name it will take default values.
      refreshMode: 1  # 配置动态刷新模式，0为configcenter在发生变化时主动推送，1为client端周期拉取，其他值均为非法，不会去连配置中心
      refreshInterval: 1 #refreshMode配置为1时，client端主动从配置中心拉取配置的周期，单位毫秒
      autodiscovery: false
      api:
        version: v3
  protocols:
    highway:
      listenAddress: 127.0.0.1:8080
      advertiseAddress: 127.0.0.1:8080
      workerNumber: 10
    rest:
      listenAddress: 127.0.0.1:8888
      advertiseAddress: 127.0.0.1:8888
      workerNumber: 10
      failure: http_500,http_502 # Defines what is considered an unsuccessful attempt of communication with a server.
  handler:
    chain:
    #  Consumer:
    #    default: bizkeeper-consumer, loadbalancer,ratelimiter-consumer
  references:    #optional：配置客户端依赖的微服务信息，协议信息
    ServerChassis:
      version: 0.1
      transport: rest
    LoadBServer:
      version: 0.1
      transport: rest
    HelloServer:
      version: 0.1
      transport: highway
ssl:
  registry.consumer.cipherPlugin: default
  registry.consumer.verifyPeer: false
  registry.consumer.cipherSuits: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  registry.consumer.protocol: TLSv1.2
  registry.consumer.caFile:
  registry.consumer.certFile:
  registry.consumer.keyFile:
  registry.consumer.certPwdFile:`

func TestRestFaultHandler_Names(t *testing.T) {
	restCon := handler.FaultHandle()
	conName := restCon.Name()
	assert.Equal(t, "fault-inject", conName)

	t.Log("testing fault-inject handler")
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	microContent := `---
#微服务的私有属性
service_description:
  name: Client
  level: FRONT
  version: 0.1`

	os.Setenv("CHASSIS_HOME", "/tmp")
	defer os.Unsetenv("CHASSIS_HOME")
	chassisConf := filepath.Join("/tmp/", "conf")
	os.MkdirAll(chassisConf, 0600)
	chassisyaml := filepath.Join(chassisConf, "chassis.yaml")
	microserviceyaml := filepath.Join(chassisConf, "microservice.yaml")
	f1, _ := os.Create(chassisyaml)
	f2, _ := os.Create(microserviceyaml)
	io.WriteString(f1, yamlContent)
	io.WriteString(f2, microContent)
	config.Init()
	archaius.Init()

	c := handler.Chain{}
	handler.RegisterHandler("fault-inject", handler.FaultHandle)
	c.AddHandler(&handler.FaultHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Consumer["fault-inject"] = "fault-inject"

	inv := &invocation.Invocation{
		MicroServiceName: "ShoppingCart",
	}

	c.Next(inv, func(r *invocation.InvocationResponse) error {
		assert.Error(t, errors.New("injecting abort and delay"), r.Err)
		log.Println(r.Result)
		return r.Err
	})
}
