package chassis

import (
	"fmt"
	"github.com/go-chassis/go-chassis/pkg/metrics"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	//init logger first
	_ "github.com/go-chassis/go-chassis/initiator"

	"github.com/go-chassis/go-chassis/bootstrap"

	//load balancing
	_ "github.com/go-chassis/go-chassis/pkg/loadbalancing"

	//protocols
	_ "github.com/go-chassis/go-chassis/client/rest"
	_ "github.com/go-chassis/go-chassis/server/restful"

	//routers
	_ "github.com/go-chassis/go-chassis/core/router/servicecomb"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/registry"

	//control panel
	_ "github.com/go-chassis/go-chassis/control/archaius"

	// registry
	_ "github.com/go-chassis/go-chassis/core/registry/file"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"

	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/core/tracing"
	"github.com/go-chassis/go-chassis/eventlistener"

	// prometheus reporter for circuit breaker metrics
	_ "github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix/reporter"

	// aes package handles security related plugins
	_ "github.com/go-chassis/go-chassis/security/plugins/aes"
	_ "github.com/go-chassis/go-chassis/security/plugins/plain"

	//config centers
	_ "github.com/go-chassis/go-chassis-config/configcenter"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/configcenter"
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/metadata"
	"github.com/go-chassis/go-chassis/pkg/circuit"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
)

var goChassis *chassis

func init() {
	goChassis = &chassis{}
}

type chassis struct {
	version     string
	schemas     []*Schema
	mu          sync.Mutex
	Initialized bool

	DefaultConsumerChainNames map[string]string
	DefaultProviderChainNames map[string]string
}

// Schema struct for to represent schema info
type Schema struct {
	serverName string
	schema     interface{}
	opts       []server.RegisterOption
}

func (c *chassis) initChains(chainType string) error {
	var defaultChainName = "default"
	var handlerNameMap = map[string]string{defaultChainName: ""}
	switch chainType {
	case common.Provider:
		if providerChainMap := config.GlobalDefinition.Cse.Handler.Chain.Provider; len(providerChainMap) != 0 {
			if _, ok := providerChainMap[defaultChainName]; !ok {
				providerChainMap[defaultChainName] = c.DefaultProviderChainNames[defaultChainName]
			}
			handlerNameMap = providerChainMap
		} else {
			handlerNameMap = c.DefaultProviderChainNames
		}
	case common.Consumer:
		if consumerChainMap := config.GlobalDefinition.Cse.Handler.Chain.Consumer; len(consumerChainMap) != 0 {
			if _, ok := consumerChainMap[defaultChainName]; !ok {
				consumerChainMap[defaultChainName] = c.DefaultConsumerChainNames[defaultChainName]
			}
			handlerNameMap = consumerChainMap
		} else {
			handlerNameMap = c.DefaultConsumerChainNames
		}
	}
	openlogging.GetLogger().Debugf("Init %s's handlermap", chainType)
	return handler.CreateChains(chainType, handlerNameMap)
}
func (c *chassis) initHandler() error {
	if err := c.initChains(common.Provider); err != nil {
		openlogging.GetLogger().Errorf("chain int failed: %s", err)
		return err
	}
	if err := c.initChains(common.Consumer); err != nil {
		openlogging.GetLogger().Errorf("chain int failed: %s", err)
		return err
	}
	openlogging.Info("chain init success")
	return nil
}

//Init
func (c *chassis) initialize() error {
	if c.Initialized {
		return nil
	}
	if err := config.Init(); err != nil {
		openlogging.Error("failed to initialize conf: " + err.Error())
		return err
	}
	if err := runtime.Init(); err != nil {
		return err
	}

	err := c.initHandler()
	if err != nil {
		openlogging.GetLogger().Errorf("handler init failed: %s", err)
		return err
	}

	err = server.Init()
	if err != nil {
		return err
	}
	bootstrap.Bootstrap()
	if !archaius.GetBool("cse.service.registry.disabled", false) {
		err := registry.Enable()
		if err != nil {
			return err
		}
		strategyName := archaius.GetString("cse.loadbalance.strategy.name", "")
		if err := loadbalancer.Enable(strategyName); err != nil {
			return err
		}
	}

	err = configcenter.Init()
	if err != nil {
		openlogging.Warn("lost config server: " + err.Error())
	}
	// router needs get configs from config-center when init
	// so it must init after bootstrap
	if err = router.Init(); err != nil {
		return err
	}
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	if err := control.Init(opts); err != nil {
		return err
	}

	if err = tracing.Init(); err != nil {
		return err
	}
	if err = metrics.Init(); err != nil {
		return err
	}
	go hystrix.StartReporter()
	circuit.Init()
	eventlistener.Init()
	c.Initialized = true
	return nil
}

func (c *chassis) registerSchema(serverName string, structPtr interface{}, opts ...server.RegisterOption) {
	schema := &Schema{
		serverName: serverName,
		schema:     structPtr,
		opts:       opts,
	}
	c.mu.Lock()
	c.schemas = append(c.schemas, schema)
	c.mu.Unlock()
}

func (c *chassis) start() error {
	if !c.Initialized {
		return fmt.Errorf("the chassis do not init. please run chassis.Init() first")
	}

	for _, v := range c.schemas {
		if v == nil {
			continue
		}
		s, err := server.GetServer(v.serverName)
		if err != nil {
			return err
		}
		_, err = s.Register(v.schema, v.opts...)
		if err != nil {
			return err
		}
	}
	err := server.StartServer()
	if err != nil {
		return err
	}
	return nil
}

//RegisterSchema Register a API service to specific server by name
//You must register API first before Call Init
func RegisterSchema(serverName string, structPtr interface{}, opts ...server.RegisterOption) {
	goChassis.registerSchema(serverName, structPtr, opts...)
}

//SetDefaultConsumerChains your custom chain map for Consumer,if there is no config, this default chain will take affect
func SetDefaultConsumerChains(c map[string]string) {
	goChassis.DefaultConsumerChainNames = c
}

//SetDefaultProviderChains set your custom chain map for Provider,if there is no config, this default chain will take affect
func SetDefaultProviderChains(c map[string]string) {
	goChassis.DefaultProviderChainNames = c
}

//Run bring up the service,it waits for os signal,and shutdown gracefully
//before all protocol server start successfully, it may return error.
func Run() error {
	err := goChassis.start()
	if err != nil {
		openlogging.Error("run chassis failed:" + err.Error())
		return err
	}
	if !config.GetRegistratorDisable() {
		//Register instance after Server started
		if err := registry.DoRegister(); err != nil {
			openlogging.Error("register instance failed:" + err.Error())
			return err
		}
	}
	//Graceful shutdown
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	select {
	case s := <-c:
		openlogging.Info("got os signal " + s.String())
	case err := <-server.ErrRuntime:
		openlogging.Info("got server error " + err.Error())
	}
	for name, s := range server.GetServers() {
		openlogging.Info("stopping server " + name + "...")
		err := s.Stop()
		if err != nil {
			openlogging.GetLogger().Warnf("servers failed to stop: %s", err)
		}
		openlogging.Info(name + " server stop success")
	}
	if !config.GetRegistratorDisable() {
		if err = server.UnRegistrySelfInstances(); err != nil {
			openlogging.GetLogger().Warnf("servers failed to unregister: %s", err)
		}
	}
	openlogging.Info("go chassis server gracefully shutdown")
	return nil
}

//Init prepare the chassis framework runtime
func Init() error {
	if goChassis.DefaultConsumerChainNames == nil {
		defaultChain := strings.Join([]string{
			handler.RatelimiterConsumer,
			handler.Router,
			handler.Loadbalance,
			handler.TracingConsumer,
			handler.Transport,
		}, ",")
		goChassis.DefaultConsumerChainNames = map[string]string{
			common.DefaultKey: defaultChain,
		}
	}
	if goChassis.DefaultProviderChainNames == nil {
		defaultChain := strings.Join([]string{
			handler.RatelimiterProvider,
			handler.TracingProvider,
		}, ",")
		goChassis.DefaultProviderChainNames = map[string]string{
			common.DefaultKey: defaultChain,
		}
	}
	if err := goChassis.initialize(); err != nil {
		log.Println("Init chassis fail:", err)
		return err
	}
	openlogging.GetLogger().Infof("Init chassis success, Version is %s", metadata.SdkVersion)
	return nil
}
