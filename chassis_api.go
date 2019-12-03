package chassis

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	//init logger first
	_ "github.com/go-chassis/go-chassis/initiator"
	//load balancing
	_ "github.com/go-chassis/go-chassis/pkg/loadbalancing"
	//protocols
	_ "github.com/go-chassis/go-chassis/client/rest"
	_ "github.com/go-chassis/go-chassis/server/restful"
	//routers
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/registry"
	//router
	_ "github.com/go-chassis/go-chassis/core/router/servicecomb"
	//control panel
	_ "github.com/go-chassis/go-chassis/control/archaius"
	// registry
	_ "github.com/go-chassis/go-chassis/core/registry/file"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/core/server"
	// prometheus reporter for circuit breaker metrics
	_ "github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix/reporter"
	// aes package handles security related plugins
	_ "github.com/go-chassis/go-chassis/security/plugins/aes"
	_ "github.com/go-chassis/go-chassis/security/plugins/plain"
	//config centers
	_ "github.com/go-chassis/go-archaius/source/remote/configcenter"
	"github.com/go-chassis/go-chassis/core/metadata"
	"github.com/go-mesh/openlogging"
)

var goChassis *chassis

func init() {
	goChassis = &chassis{}
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
	waitingSignal()
	return nil
}

func waitingSignal() {
	//Graceful shutdown
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	select {
	case s := <-c:
		openlogging.Info("got os signal " + s.String())
	case err := <-server.ErrRuntime:
		openlogging.Info("got server error " + err.Error())
	}

	if !config.GetRegistratorDisable() {
		registry.HBService.Stop()
		openlogging.Info("unregister servers ...")
		if err := server.UnRegistrySelfInstances(); err != nil {
			openlogging.GetLogger().Warnf("servers failed to unregister: %s", err)
		}
	}

	for name, s := range server.GetServers() {
		openlogging.Info("stopping server " + name + "...")
		err := s.Stop()
		if err != nil {
			openlogging.GetLogger().Warnf("servers failed to stop: %s", err)
		}
		openlogging.Info(name + " server stop success")
	}

	openlogging.Info("go chassis server gracefully shutdown")
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
		log.Println("init chassis fail:", err)
		return err
	}
	openlogging.GetLogger().Infof("init chassis success, version is %s", metadata.SdkVersion)
	return nil
}
