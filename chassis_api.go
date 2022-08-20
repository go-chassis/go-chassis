package chassis

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chassis/openlog"
	//init logger first
	_ "github.com/go-chassis/go-chassis/v2/initiator"

	// transport handler
	_ "github.com/go-chassis/go-chassis/v2/core/client"

	//load balancing
	_ "github.com/go-chassis/go-chassis/v2/pkg/loadbalancing"

	//protocols
	_ "github.com/go-chassis/go-chassis/v2/client/rest"
	_ "github.com/go-chassis/go-chassis/v2/server/restful"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	//router
	_ "github.com/go-chassis/go-chassis/v2/core/router/servicecomb"
	//control panel
	_ "github.com/go-chassis/go-chassis/v2/control/servicecomb"
	// registry
	_ "github.com/go-chassis/go-chassis/v2/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/v2/core/server"
	// prometheus reporter for circuit breaker metrics
	_ "github.com/go-chassis/go-chassis/v2/third_party/forked/afex/hystrix-go/hystrix/reporter"
	// aes package handles security related plugins
	_ "github.com/go-chassis/go-chassis/v2/security/cipher/plugins/aes"
	_ "github.com/go-chassis/go-chassis/v2/security/cipher/plugins/plain"
	//config servers
	_ "github.com/go-chassis/go-archaius/source/remote"
	_ "github.com/go-chassis/go-archaius/source/remote/kie"
	"github.com/go-chassis/go-chassis/v2/core/metadata"
)

var goChassis *chassis

func init() {
	goChassis = &chassis{}
}

// RegisterSchema Register a API service to specific server by name
// You must register API first before Call Init
func RegisterSchema(serverName string, structPtr interface{}, opts ...server.RegisterOption) {
	goChassis.registerSchema(serverName, structPtr, opts...)
}

// SetDefaultConsumerChains your custom chain map for Consumer,if there is no config, this default chain will take affect
func SetDefaultConsumerChains(c map[string]string) {
	goChassis.DefaultConsumerChainNames = c
}

// SetDefaultProviderChains set your custom chain map for Provider,if there is no config, this default chain will take affect
func SetDefaultProviderChains(c map[string]string) {
	goChassis.DefaultProviderChainNames = c
}

// HijackSignal set signals that want to hijack.
func HijackSignal(sigs ...os.Signal) {
	goChassis.sigs = sigs
}

// InstallPreShutdown instal what you want to achieve before graceful shutdown
func InstallPreShutdown(name string, f func(os.Signal)) {
	// lazy init
	if goChassis.preShutDownFuncs == nil {
		goChassis.preShutDownFuncs = make(map[string]func(os.Signal))
	}
	goChassis.preShutDownFuncs[name] = f
}

// InstallPostShutdown instal what you want to achieve after graceful shutdown
func InstallPostShutdown(name string, f func(os.Signal)) {
	// lazy init
	if goChassis.postShutDownFuncs == nil {
		goChassis.postShutDownFuncs = make(map[string]func(os.Signal))
	}
	goChassis.postShutDownFuncs[name] = f
}

// HijackGracefulShutdown reset GracefulShutdown
func HijackGracefulShutdown(f func(os.Signal)) {
	goChassis.hijackGracefulShutdown = f
}

// Run bring up the service,it waits for os signal,and shutdown gracefully
// before all protocol server start successfully, it may return error.
func Run(options ...server.RunOption) error {
	err := goChassis.start(options...)
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
	c := make(chan os.Signal, 1)
	if len(goChassis.sigs) > 0 {
		signal.Notify(c, goChassis.sigs...)
	} else {
		signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM,
			syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
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
	goChassis.hijackGracefulShutdown(s)
	if goChassis.postShutDownFuncs != nil {
		for k, v := range goChassis.postShutDownFuncs {
			openlog.Info(fmt.Sprintf("exec post shutdown funcs %s", k))
			v(s)
		}
	}
}

// GracefulShutdown graceful shut down api
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

// Init prepare the chassis framework runtime
func Init() error {
	if goChassis.DefaultConsumerChainNames == nil {
		defaultChain := strings.Join([]string{
			handler.Router,
			handler.LoadBalancing,
			handler.TracingConsumer,
			handler.Transport,
		}, ",")
		goChassis.DefaultConsumerChainNames = map[string]string{
			common.DefaultKey: defaultChain,
		}
	}
	if goChassis.DefaultProviderChainNames == nil {
		defaultChain := strings.Join([]string{handler.TracingProvider}, ",")
		goChassis.DefaultProviderChainNames = map[string]string{
			common.DefaultKey: defaultChain,
		}
	}
	goChassis.hijackGracefulShutdown = GracefulShutdown
	if err := goChassis.initialize(); err != nil {
		log.Println("init chassis fail:", err)
		return err
	}
	openlog.Info("init chassis success, version is " + metadata.SdkVersion)
	return nil
}
