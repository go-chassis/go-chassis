package server

import (
	"errors"
	"fmt"
	"github.com/go-chassis/go-archaius"
	"log"
	"time"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	chassisTLS "github.com/go-chassis/go-chassis/v2/core/tls"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/go-chassis/v2/pkg/util"
	"github.com/go-chassis/go-chassis/v2/pkg/util/iputil"
	"github.com/go-chassis/openlog"
)

// constants for server
const (
	DefaultMetricPath  = "metrics"
	DefaultProfilePath = "profile"
)

// NewFunc returns a ProtocolServer
type NewFunc func(Options) ProtocolServer

var serverPlugins = make(map[string]NewFunc)
var servers = make(map[string]ProtocolServer)

// InstallPlugin For developer
func InstallPlugin(protocol string, newFunc NewFunc) {
	serverPlugins[protocol] = newFunc
	openlog.Info("installed server plugin: " + protocol)
}

// GetServerFunc returns the server function
func GetServerFunc(protocol string) (NewFunc, error) {
	f, ok := serverPlugins[protocol]
	if !ok {
		return nil, fmt.Errorf("unknown protocol server [%s]", protocol)
	}
	return f, nil
}

// GetServer return the server based on protocol
func GetServer(protocol string) (ProtocolServer, error) {
	s, ok := servers[protocol]
	if !ok {
		return nil, fmt.Errorf("[%s] server isn't running ", protocol)
	}
	return s, nil
}

// GetServers returns the map of servers
func GetServers() map[string]ProtocolServer {
	return servers
}

// ErrRuntime is an error channel, if it receive any signal will cause graceful shutdown of go chassis, process will exit
var ErrRuntime = make(chan error)

// StartServer starting the server
func StartServer(options ...RunOption) error {
	opts := RunOptions{}
	for _, o := range options {
		o(&opts)
	}
	for name, server := range servers {
		openlog.Info("starting server " + name + "...")
		if opts.serverMasks.Has(name) {
			openlog.Warn("server " + name + " is masked, and will not start.")
			continue
		}
		err := server.Start()
		if err != nil {
			openlog.Error(fmt.Sprintf("servers failed to start, err %s", err))
			return fmt.Errorf("can not start [%s] server,%w", name, err)
		}
		openlog.Debug(name + " server start success")
	}
	openlog.Info("all server start completed")

	return nil
}

// UnRegistrySelfInstances this function removes the self instance
func UnRegistrySelfInstances() error {
	if err := registry.DefaultRegistrator.UnRegisterMicroServiceInstance(runtime.ServiceID, runtime.InstanceID); err != nil {
		openlog.Error(fmt.Sprintf("unregister instance failed, sid/iid: %s/%s: %s",
			runtime.ServiceID, runtime.InstanceID, err))
		return err
	}
	return nil
}

// Init initializes
func Init() error {
	var err error
	for k, v := range config.GlobalDefinition.ServiceComb.Protocols {
		if err = initialServer(config.GlobalDefinition.ServiceComb.Handler.Chain.Provider, v, k); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil

}

func initialServer(providerMap map[string]string, p model.Protocol, name string) error {
	protocolName, _, err := util.ParsePortName(name)
	if err != nil {
		return err
	}
	openlog.Debug(fmt.Sprintf("init server [%s], protocol is [%s]", name, protocolName))
	f, err := GetServerFunc(protocolName)
	if err != nil {
		return fmt.Errorf("do not support [%s] server", name)
	}

	sslTag := name + "." + common.Provider
	tlsConfig, sslConfig, err := chassisTLS.GetTLSConfigByService("", name, common.Provider)
	if err != nil {
		if !chassisTLS.IsSSLConfigNotExist(err) {
			return err
		}
	} else {
		openlog.Warn(fmt.Sprintf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
			sslTag, sslConfig.VerifyPeer, sslConfig.CipherPlugin))
	}

	if p.Listen == "" {
		if p.Advertise != "" {
			p.Listen = p.Advertise
		} else {
			p.Listen = iputil.DefaultEndpoint4Protocol(name)
		}
	}

	chainName := common.DefaultChainName
	if _, ok := providerMap[name]; ok {
		chainName = name
	}

	var s ProtocolServer
	o := Options{
		Address:            p.Listen,
		ProtocolServerName: name,
		ChainName:          chainName,
		TLSConfig:          tlsConfig,
		BodyLimit:          config.GlobalDefinition.ServiceComb.Transport.MaxBodyBytes[protocolName],
		HeaderLimit:        config.GlobalDefinition.ServiceComb.Transport.MaxHeaderBytes[protocolName],
		ProfilingAPI:       archaius.GetString("servicecomb.profile.apiPath", DefaultProfilePath),
		ProfilingEnable:    archaius.GetBool("servicecomb.profile.enable", false),
		MetricsAPI:         archaius.GetString("servicecomb.metrics.apiPath", DefaultMetricPath),
		MetricsEnable:      archaius.GetBool("servicecomb.metrics.enable", false),
	}
	if t := config.GlobalDefinition.ServiceComb.Transport.Timeout[protocolName]; len(t) > 0 {
		timeout, err := time.ParseDuration(t)
		if err != nil {
			openlog.Error(fmt.Sprintf("parse timeout failed: %s", err))
			return err
		}
		if timeout < 0 {
			return errors.New("timeout should be positive, but get: " + t)
		}
		o.Timeout = timeout
	}
	s = f(o)
	servers[name] = s
	return nil
}
