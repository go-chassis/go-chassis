package client

import (
	"fmt"
	"strings"
	"sync"

	"crypto/tls"
	"errors"
	"time"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	chassisTLS "github.com/go-chassis/go-chassis/core/tls"
)

var clients = make(map[string]ProtocolClient)
var sl sync.RWMutex

//ErrClientNotExist happens if client do not exist
var ErrClientNotExist = errors.New("client not exists")

//DefaultPoolSize is 500
const DefaultPoolSize = 50

//Options is configs for client creation
type Options struct {
	Service   string
	PoolSize  int
	Endpoint  string
	PoolTTL   time.Duration
	TLSConfig *tls.Config
	Failure   map[string]bool
}

// GetFailureMap return failure map
func GetFailureMap(p string) map[string]bool {
	failureList := strings.Split(config.GlobalDefinition.Cse.Transport.Failure[p], ",")
	failureMap := make(map[string]bool)
	for _, v := range failureList {
		if v == "" {
			continue
		}
		failureMap[v] = true
	}
	return failureMap
}

// CreateClient is for to create client based on protocol and the service name
func CreateClient(protocol, service, endpoint string) (ProtocolClient, error) {
	f, err := GetClientNewFunc(protocol)
	if err != nil {
		lager.Logger.Error(fmt.Sprintf("do not Support [%s] client", protocol))
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

	poolSize := DefaultPoolSize

	return f(Options{
		Service:   service,
		TLSConfig: tlsConfig,
		PoolSize:  poolSize,
		Failure:   GetFailureMap(protocol),
		Endpoint:  endpoint,
	})
}
func generateKey(protocol, service, endpoint string) string {
	return protocol + service + endpoint
}

// GetClient is to get the client based on protocol, service,endpoint name
func GetClient(protocol, service, endpoint string) (ProtocolClient, error) {
	var c ProtocolClient
	var err error
	key := generateKey(protocol, service, endpoint)
	sl.RLock()
	c, ok := clients[key]
	sl.RUnlock()
	if !ok {
		lager.Logger.Info("Create client for " + protocol + ":" + service + ":" + endpoint)
		c, err = CreateClient(protocol, service, endpoint)
		if err != nil {
			return nil, err
		}
		sl.Lock()
		clients[key] = c
		sl.Unlock()
	}
	return c, nil
}

//Close close a client conn
func Close(protocol, service, endpoint string) error {
	key := generateKey(protocol, service, endpoint)
	sl.RLock()
	c, ok := clients[key]
	sl.RUnlock()
	if !ok {
		return ErrClientNotExist
	}
	if err := c.Close(); err != nil {
		lager.Logger.Errorf("can not close client %s:%s%:s, err [%s]", protocol, service, endpoint, err.Error())
		return err
	}
	sl.Lock()
	delete(clients, key)
	sl.Unlock()
	return nil
}

// GetAllClient get all client cache for cache
func GetAllClient() map[string]ProtocolClient {
	return clients
}
