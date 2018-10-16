package grpc

import (
	"crypto/tls"
	"errors"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-mesh/openlogging"
	"google.golang.org/grpc"
	"net"
)

//err define
var (
	ErrGRPCSvcDescMissing = errors.New("must use server.WithGRPCServiceDesc to set desc")
)

//const
const (
	Name = "grpc"
)

//Server is grpc server holder
type Server struct {
	s    *grpc.Server
	opts server.Options
}

//New create grpc server
func New(opts server.Options) server.ProtocolServer {
	return &Server{
		opts: opts,
		s:    grpc.NewServer(),
	}
}

//Register register grpc services
func (s *Server) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	opts := server.RegisterOptions{}
	for _, o := range options {
		o(&opts)
	}
	if opts.GRPCSvcDesc == nil {
		return "", ErrGRPCSvcDescMissing
	}
	s.s.RegisterService(opts.GRPCSvcDesc, schema)
	return "", nil
}

//Start launch the server
func (s *Server) Start() error {
	var listener net.Listener
	var lisErr error
	if s.opts.TLSConfig == nil {
		listener, lisErr = net.Listen("tcp", s.opts.Address)
	} else {
		listener, lisErr = tls.Listen("tcp", s.opts.Address, s.opts.TLSConfig)
	}
	if lisErr != nil {
		openlogging.GetLogger().Error("listening failed, reason:" + lisErr.Error())
		return lisErr
	}
	go func() {
		if err := s.s.Serve(listener); err != nil {
			server.ErrRuntime <- err
		}
	}()
	return nil
}

//Stop gracfully shutdown grpc server
func (s *Server) Stop() error {
	s.s.GracefulStop()
	return nil
}

//String return server name
func (s *Server) String() string {
	return Name
}
func init() {
	server.InstallPlugin(Name, New)
}
