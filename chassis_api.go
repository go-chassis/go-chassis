package chassis

import (
	"fmt"
	"github.com/go-chassis/openlog"
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
	_ "github.com/go-chassis/go-chassis/control/servicecomb"
	// registry
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/core/server"
	// prometheus reporter for circuit breaker metrics
	_ "github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix/reporter"
	// aes package handles security related plugins
	_ "github.com/go-chassis/go-chassis/security/cipher/plugins/aes"
	_ "github.com/go-chassis/go-chassis/security/cipher/plugins/plain"
	//config servers
	_ "github.com/go-chassis/go-archaius/source/remote"
	_ "github.com/go-chassis/go-archaius/source/remote/kie"
	"github.com/go-chassis/go-chassis/core/metadata"
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

//HajackSignal set signals that want to hajack.
func HajackSignal(sigs ...os.Signal) {
	goChassis.sigs = sigs
}

//InstalPreShutdown instal what you want to achieve before graceful shutdown
func InstalPreShutdown(name string, f func(os.Signal)) {
	// lazy init
	if goChassis.preShutDownFuncs == nil {
		goChassis.preShutDownFuncs = make(map[string]func(os.Signal))
	}
	goChassis.preShutDownFuncs[name] = f
}

//InstalPostShutdown instal what you want to achieve after graceful shutdown
func InstalPostShutdown(name string, f func(os.Signal)) {
	// lazy init
	if goChassis.postShutDownFuncs == nil {
		goChassis.postShutDownFuncs = make(map[string]func(os.Signal))
	}
	goChassis.postShutDownFuncs[name] = f
}

//HajackGracefulShutdown reset GracefulShutdown
func HajackGracefulShutdown(f func(os.Signal)) {
	goChassis.hajackGracefulShutdown = f
}

//Run bring up the service,it waits for os signal,and shutdown gracefully
//before all protocol server start successfully, it may return error.
func Run() error {
	err := goChassis.start()
	if err != nil {
		openlog.Error("run chassis failed:" + err.Error())
		return err
	}
	if !config.GetRegistratorDisable() {
		//Register instance after Server started
		if err := registry.DoRegister(); err != nil {
			openlog.Error("register instance failed:" + err.Error())
			return err
		}
	}

	waitingSignal()
	return nil
}

func waitingSignal() {
	c := make(chan os.Signal)
	if len(goChassis.sigs) > 0 {
		signal.Notify(c, goChassis.sigs...)
	} else {
		signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	}

	var s os.Signal
	select {
	case s = <-c:
		openlog.Info("got os signal " + s.String())
	case err := <-server.ErrRuntime:
		openlog.Info("got server error " + err.Error())
	}

	if goChassis.preShutDownFuncs != nil {
		for k, v := range goChassis.preShutDownFuncs {
			openlog.Info(fmt.Sprintf("exec pre shutdown funcs %s", k))
			v(s)
		}
	}
	goChassis.hajackGracefulShutdown(s)
	if goChassis.postShutDownFuncs != nil {
		for k, v := range goChassis.postShutDownFuncs {
			openlog.Info(fmt.Sprintf("exec post shutdown funcs %s", k))
			v(s)
		}
	}
}

//GracefulShutdown graceful shut down api
func GracefulShutdown(s os.Signal) {
	if !config.GetRegistratorDisable() {
		registry.HBService.Stop()
		openlog.Info("unregister servers ...")
		if err := server.UnRegistrySelfInstances(); err != nil {
			openlog.Warn("servers failed to unregister: " + err.Error())
		}
	}

	for name, s := range server.GetServers() {
		openlog.Info("stopping server " + name + "...")
		err := s.Stop()
		if err != nil {
			openlog.Warn("servers failed to stop: " + err.Error())
		}
		openlog.Info(name + " server stop success")
	}

	openlog.Info("go chassis server gracefully shutdown")
}

//Init prepare the chassis framework runtime
func Init() error {
	if goChassis.DefaultConsumerChainNames == nil {
		defaultChain := strings.Join([]string{
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
			handler.TracingProvider,
		}, ",")
		goChassis.DefaultProviderChainNames = map[string]string{
			common.DefaultKey: defaultChain,
		}
	}
	goChassis.hajackGracefulShutdown = GracefulShutdown
	if err := goChassis.initialize(); err != nil {
		log.Println("init chassis fail:", err)
		return err
	}
	openlog.Info("init chassis success, version is " + metadata.SdkVersion)
	return nil
}
