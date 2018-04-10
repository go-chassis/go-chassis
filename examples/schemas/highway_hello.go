package schemas

import (
	"context"

	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"time"
)

//HelloServer is a struct
type HelloServer struct {
}

//SayHello is a method used to reply message
func (s *HelloServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if in.Name == "a" {
		<-time.After(1000 * time.Millisecond)
	}
	return &helloworld.HelloReply{Message: "Go Hello  " + in.Name}, nil
}
