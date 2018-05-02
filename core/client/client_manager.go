package client

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	chassisTLS "github.com/ServiceComb/go-chassis/core/tls"
)

var clients = make(map[string]map[string]ProtocolClient)
var pl sync.RWMutex
var sl sync.RWMutex

//DefaultPoolSize is 500
const DefaultPoolSize = 50

// GetProtocolSpec is to get protocol specifications
func GetProtocolSpec(p string) model.Protocol {
	return config.GlobalDefinition.Cse.Protocols[p]
}

// CreateClient is for to create client based on protocol and the service name
func CreateClient(protocol, service string) (ProtocolClient, error) {
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
		lager.Logger.Warnf("%s %s TLS mode, verify peer: %t, cipher plugin: %s.",
			protocol, service, sslConfig.VerifyPeer, sslConfig.CipherPlugin)
	}
	p := GetProtocolSpec(protocol)

	poolSize := DefaultPoolSize

	failureList := strings.Split(p.Failure, ",")
	failureMap := make(map[string]bool)
	for _, v := range failureList {
		if v == "" {
			continue
		}
		failureMap[v] = true
	}

	c := f(Options{
		TLSConfig: tlsConfig,
		PoolSize:  poolSize,
		Failure:   failureMap,
	})

	return c, nil
}

// GetClient is to get the client based on protocol and service name
func GetClient(protocol, service string) (ProtocolClient, error) {
	var c ProtocolClient
	var err error
	pl.RLock()
	clientMap, ok := clients[protocol]
	pl.RUnlock()
	if !ok {
		lager.Logger.Info("Create client map for " + protocol)
		clientMap = make(map[string]ProtocolClient)
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
