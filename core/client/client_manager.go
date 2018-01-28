package client

import (
	"strings"

	"fmt"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	chassisTLS "github.com/ServiceComb/go-chassis/core/tls"
	"github.com/ServiceComb/go-chassis/core/transport"
	microClient "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	microTransport "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	"sync"
)

var clients = make(map[string]map[string]microClient.Client)
var pl sync.RWMutex
var sl sync.RWMutex

// GetProtocolSpec is to get protocol specifications
func GetProtocolSpec(p string) model.Protocol {
	return config.GlobalDefinition.Cse.Protocols[p]
}

// CreateClient is for to create client based on protocol and the service name
func CreateClient(protocol, service string) (microClient.Client, error) {
	f, err := GetClientNewFunc(protocol)
	if err != nil {
		err = fmt.Errorf("don not Support [%s] client", protocol)
		lager.Logger.Error("", err)
		return nil, err
	}
	tlsConfig, sslConfig, err := chassisTLS.GetTLSConfigByService(service, protocol, common.Consumer)
	if err != nil {
		if !chassisTLS.IsSSLConfigNotExist(err) {
			return nil, err
		}
	} else {
		lager.Logger.Warnf(nil, "%s %s TLS mode, verify peer: %t, cipher plugin: %s.",
			protocol, service, sslConfig.VerifyPeer, sslConfig.CipherPlugin)
	}
	p := GetProtocolSpec(protocol)
	var tr microTransport.Transport
	if p.Transport == "" {
		p.Transport = common.TransportTCP
	}
	trF, err := transport.GetTransportFunc(p.Transport)
	if err != nil {
		return nil, err
	}
	tr = trF(microTransport.TLSConfig(tlsConfig))

	poolSize := microClient.DefaultPoolSize

	failureList := strings.Split(p.Failure, ",")
	failureMap := make(map[string]bool)
	for _, v := range failureList {
		if v == "" {
			continue
		}
		failureMap[v] = true
	}

	c := f(
		microClient.Transport(tr),
		microClient.ContentType("application/json"),
		microClient.TLSConfig(tlsConfig),
		microClient.WithConnectionPoolSize(poolSize),
		microClient.WithFailure(failureMap))

	if err = c.Init(); err != nil {
		return nil, err
	}
	return c, nil
}

// GetClient is to get the client based on protocol and service name
func GetClient(protocol, service string) (microClient.Client, error) {
	var c microClient.Client
	var err error
	pl.RLock()
	clientMap, ok := clients[protocol]
	pl.RUnlock()
	if !ok {
		lager.Logger.Info("Create client map for " + protocol)
		clientMap = make(map[string]microClient.Client)
		pl.Lock()
		clients[protocol] = clientMap
		pl.Unlock()
	}
	sl.RLock()
	c, ok = clientMap[service]
	sl.RUnlock()
	if !ok {
		lager.Logger.Info("Create client for " + service)
		c, err = CreateClient(protocol, service)
		if err != nil {
			return nil, err
		}
		sl.Lock()
		clientMap[service] = c
		sl.Unlock()
	}
	return c, nil
}
