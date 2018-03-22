package highway

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"sync"

	"crypto/tls"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/provider"
	"github.com/ServiceComb/go-chassis/core/server"
	"net"
	"time"
)

const (
	//Name is a variable of type string which says about the protocol used
	Name = "highway"
)

// constants for request and login
const (
	Request = 0
	Login   = 1
)

var remoteLogin = true

type highwayServer struct {
	connMgr *ConnectionMgr
	opts    server.Options
	sync.RWMutex
}

func (s *highwayServer) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	opts := server.RegisterOptions{}
	var pn string
	for _, o := range options {
		o(&opts)
	}
	mc := config.MicroserviceDefinition
	if mc == nil {
		pn = common.DefaultProvider
	}
	if mc == nil || mc.Provider == "" {
		pn = common.DefaultProvider
	} else {
		if mc.Provider == "" {
			pn = common.DefaultProvider
		} else {
			pn = mc.Provider
		}

	}
	provider.RegisterProvider(pn, config.SelfServiceName)
	if opts.SchemaID != "" {
		err := provider.RegisterSchemaWithName(config.SelfServiceName, opts.SchemaID, schema)
		return opts.SchemaID, err
	}
	schemaID, err := provider.RegisterSchema(config.SelfServiceName, schema)
	return schemaID, err
}

func (s *highwayServer) Start() error {
	opts := s.opts
	//TODO lot of options
	var listener net.Listener
	var lisErr error
	if s.opts.TLSConfig == nil {
		listener, lisErr = net.Listen("tcp", opts.Address)
	} else {
		listener, lisErr = tls.Listen("tcp", opts.Address, s.opts.TLSConfig)
	}

	if lisErr != nil {
		lager.Logger.Error("listening falied, reason:", lisErr)
		return lisErr
	}
	go s.AcceptLoop(listener)
	return nil
}

func (s *highwayServer) AcceptLoop(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			lager.Logger.Errorf(err, "Error accepting")
			select {
			case <-time.After(time.Second * 3):
				lager.Logger.Info("Sleep three second")
			}
		}
		highwayConn := s.connMgr.createConn(conn, s.opts.ChainName)
		highwayConn.Open()
	}
}

func (s *highwayServer) Stop() error {
	s.connMgr.DeactiveAllConn()
	return nil
}

func newHighwayServer(opts server.Options) server.ProtocolServer {
	return &highwayServer{
		connMgr: newConnectMgr(),
		opts:    opts,
	}
}
func (s *highwayServer) String() string {
	return Name
}
func init() {
	server.InstallPlugin(Name, newHighwayServer)
}
