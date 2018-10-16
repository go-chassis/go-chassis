package main

import (
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	pb "github.com/go-chassis/go-chassis/examples/grpc/helloworld"
	_ "github.com/go-chassis/go-chassis/server/grpc"
	"golang.org/x/net/context"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rpc/server/
// Server is used to implement helloworld.GreeterServer.
type Server struct{}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
func main() {
	chassis.RegisterSchema("grpc", &Server{}, server.WithGRPCServiceDesc(&pb.Greeter_serviceDesc))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.")
		return
	}
	chassis.Run()
}
