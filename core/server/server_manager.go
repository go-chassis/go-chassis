package server

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/registry"
	chassisTLS "github.com/go-chassis/go-chassis/core/tls"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/util"
	"github.com/go-chassis/go-chassis/pkg/util/iputil"
	"github.com/go-mesh/openlogging"
)

//NewFunc returns a ProtocolServer
type NewFunc func(Options) ProtocolServer

var serverPlugins = make(map[string]NewFunc)
var servers = make(map[string]ProtocolServer)

//InstallPlugin For developer
func InstallPlugin(protocol string, newFunc NewFunc) {
	serverPlugins[protocol] = newFunc
	openlogging.Info("Installed Server Plugin, protocol:" + protocol)
}

//GetServerFunc returns the server function
func GetServerFunc(protocol string) (NewFunc, error) {
	f, ok := serverPlugins[protocol]
	if !ok {
		return nil, fmt.Errorf("unknown protocol server [%s]", protocol)
	}
	return f, nil
}

//GetServer return the server based on protocol
func GetServer(protocol string) (ProtocolServer, error) {
	s, ok := servers[protocol]
	if !ok {
		return nil, fmt.Errorf("[%s] server isn't running ", protocol)
	}
	return s, nil
}

//GetServers returns the map of servers
func GetServers() map[string]ProtocolServer {
	return servers
}

//ErrRuntime is an error channel, if it receive any signal will cause graceful shutdown of go chassis, process will exit
var ErrRuntime = make(chan error)

//StartServer starting the server
func StartServer() error {
	for name, server := range servers {
		openlogging.GetLogger().Info("starting server " + name + "...")
		err := server.Start()
		if err != nil {
			openlogging.GetLogger().Errorf("servers failed to start, err %s", err)
			return fmt.Errorf("can not start [%s] server,%s", name, err.Error())
		}
		openlogging.GetLogger().Debug(name + " server start success")
	}
	openlogging.GetLogger().Info("all server start completed")

	return nil
}

//UnRegistrySelfInstances this function removes the self instance
func UnRegistrySelfInstances() error {
	if err := registry.DefaultRegistrator.UnRegisterMicroServiceInstance(runtime.ServiceID, runtime.InstanceID); err != nil {
		openlogging.GetLogger().Errorf("StartServer() UnregisterMicroServiceInstance failed, sid/iid: %s/%s: %s",
			runtime.ServiceID, runtime.InstanceID, err)
		return err
	}
	return nil
}

//Init initializes
func Init() error {
	var err error
	for k, v := range config.GlobalDefinition.Cse.Protocols {
		if err = initialServer(config.GlobalDefinition.Cse.Handler.Chain.Provider, v, k); err != nil {
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
	openlogging.GetLogger().Debugf("Init server [%s], protocol is [%s]", name, protocolName)
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
		openlogging.GetLogger().Warnf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
			sslTag, sslConfig.VerifyPeer, sslConfig.CipherPlugin)
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
		BodyLimit:          config.GlobalDefinition.Cse.Transport.MaxBodyBytes[protocolName],
		HeaderLimit:        config.GlobalDefinition.Cse.Transport.MaxHeaderBytes[protocolName],
	}
	if t := config.GlobalDefinition.Cse.Transport.Timeout[protocolName]; len(t) > 0 {
		timeout, err := time.ParseDuration(t)
		if err != nil {
			openlogging.GetLogger().Errorf("parse timeout failed: %s", err)
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
