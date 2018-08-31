package highway

import (
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"sync"

	"crypto/tls"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/provider"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/pkg/runtime"
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
	provider.RegisterProvider(pn, runtime.ServiceName)
	if opts.SchemaID != "" {
		err := provider.RegisterSchemaWithName(runtime.ServiceName, opts.SchemaID, schema)
		return opts.SchemaID, err
	}
	schemaID, err := provider.RegisterSchema(runtime.ServiceName, schema)
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
		lager.Logger.Error("listening failed, reason:" + lisErr.Error())
		return lisErr
	}
	go s.acceptLoop(listener)
	return nil
}

func (s *highwayServer) acceptLoop(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			lager.Logger.Errorf("Error accepting, err [%s]", err)
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
