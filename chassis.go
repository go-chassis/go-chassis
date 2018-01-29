package chassis

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/ServiceComb/go-chassis/auth"
	"github.com/ServiceComb/go-chassis/bootstrap"
	// highway package handles remote procedure calls
	_ "github.com/ServiceComb/go-chassis/client/highway"
	// rest package handle rest apis
	_ "github.com/ServiceComb/go-chassis/client/rest"
	// archaius package to get the conguration info fron diffent configuration sources
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/core/registry"
	// file package for file based registration
	_ "github.com/ServiceComb/go-chassis/core/registry/file"
	// servicecenter package handles service center api calls
	_ "github.com/ServiceComb/go-chassis/core/registry/servicecenter"
	"github.com/ServiceComb/go-chassis/core/route"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/core/tracing"
	"github.com/ServiceComb/go-chassis/eventlistener"
	// aes package handles security related plugins
	_ "github.com/ServiceComb/go-chassis/security/plugins/aes"
	_ "github.com/ServiceComb/go-chassis/security/plugins/plain"
	_ "github.com/ServiceComb/go-chassis/server/highway"
	_ "github.com/ServiceComb/go-chassis/server/restful"
	serverOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
	// tcp package handles transport related things
	_ "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport/tcp"
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
	protocol string
	schema   interface{}
	opts     []serverOption.RegisterOption
}

func (c *chassis) initChains(chainType string) error {
	var defaultChainName = "default"
	var handlerNameMap = map[string]string{defaultChainName: ""}
	switch chainType {
	case common.Provider:
		if len(config.GlobalDefinition.Cse.Handler.Chain.Provider) != 0 {
			handlerNameMap = config.GlobalDefinition.Cse.Handler.Chain.Provider
		} else {
			handlerNameMap = c.DefaultProviderChainNames
		}
	case common.Consumer:
		if len(config.GlobalDefinition.Cse.Handler.Chain.Consumer) != 0 {
			handlerNameMap = config.GlobalDefinition.Cse.Handler.Chain.Consumer
		} else {
			handlerNameMap = c.DefaultConsumerChainNames
		}
	}
	lager.Logger.Debugf("Init %s's handlermap", chainType)
	return handler.CreateChains(chainType, handlerNameMap)
}
func (c *chassis) initHandler() error {
	if err := c.initChains(common.Provider); err != nil {
		lager.Logger.Errorf(err, "chain int failed")
		return err
	}
	if err := c.initChains(common.Consumer); err != nil {
		lager.Logger.Errorf(err, "chain int failed")
		return err
	}
	lager.Logger.Info("Chain init success")
	return nil
}

//Init
func (c *chassis) initialize() error {
	if c.Initialized {
		return nil
	}
	err := config.Init()
	if err != nil {
		lager.Logger.Error("Failed to initialize conf,", err)
		return err
	}
	router.Init(config.GetRouterConfig().Destinations, config.GetRouterConfig().SourceTemplates)

	auth.Init()

	err = c.initHandler()
	if err != nil {
		lager.Logger.Errorf(err, "Handler init failed")
		return err
	}

	bootstrap.Bootstrap()

	err = server.Init()
	if err != nil {
		return err
	}

	if archaius.GetBool("cse.service.registry.disabled", false) != true {
		err := registry.Enable()
		if err != nil {
			return err
		}
		if err := loadbalance.Enable(); err != nil {
			return err
		}
	}
	err = tracing.Init()
	if err != nil {
		return err
	}

	eventlistener.Init()
	c.Initialized = true
	return nil
}

func (c *chassis) registerSchema(protocol string, structPtr interface{}, opts ...serverOption.RegisterOption) {
	schema := &Schema{
		protocol: protocol,
		schema:   structPtr,
		opts:     opts,
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
		s, err := server.GetServer(v.protocol)
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

//RegisterSchema Register a API service to specific protocol
//You must register API first before Call Init
func RegisterSchema(protocol string, structPtr interface{}, opts ...serverOption.RegisterOption) {
	goChassis.registerSchema(protocol, structPtr, opts...)
}

//SetDefaultConsumerChains your custom chain map for Consumer,if there is no config, this default chain will take affect
func SetDefaultConsumerChains(c map[string]string) {
	goChassis.DefaultConsumerChainNames = c
}

//SetDefaultProviderChains set your custom chain map for Provider,if there is no config, this default chain will take affect
func SetDefaultProviderChains(c map[string]string) {
	goChassis.DefaultProviderChainNames = c
}

//Run bring up the service,it will not return error,instead just waiting for os signal,and shutdown gracefully
func Run() {
	err := goChassis.start()
	if err != nil {
		lager.Logger.Error("run chassis fail:", err)
	}
	//Register instance after Server started
	if err := registry.DoRegister(); err != nil {
		lager.Logger.Error("register instance fail:", err)
	}
	//Graceful shutdown
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	select {
	case s := <-c:
		lager.Logger.Info("got os signal " + s.String())
	case err := <-server.ServerErr:
		lager.Logger.Info("got Server Error " + err.Error())
	}
	for name, s := range server.GetServers() {
		lager.Logger.Info("stopping server " + name + "...")
		err := s.Stop()
		if err != nil {
			lager.Logger.Errorf(err, "servers failed to stop")
		}
		lager.Logger.Info(name + " server stop success")
	}
	if err = server.UnRegistrySelfInstances(); err != nil {
		lager.Logger.Errorf(err, "servers failed to unregister")
	}
}

//Init prepare the chassis framework runtime
func Init() error {
	if goChassis.DefaultConsumerChainNames == nil {
		defaultChain := strings.Join([]string{
			handler.RatelimiterConsumer,
			handler.BizkeeperConsumer,
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
			handler.BizkeeperProvider,
		}, ",")
		goChassis.DefaultProviderChainNames = map[string]string{
			common.DefaultKey: defaultChain,
		}
	}
	if err := goChassis.initialize(); err != nil {
		log.Println("Init chassis fail:", err)
		return err
	}
	lager.Logger.Infof("Init chassis success")
	return nil
}
