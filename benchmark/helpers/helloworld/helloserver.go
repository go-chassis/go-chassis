package helloworld

import (
	pb "github.com/go-chassis/go-chassis/benchmark/helpers/helloworld/protobuf"
	"github.com/go-chassis/go-chassis/server/restful"
	"golang.org/x/net/context"
	"math"
	"net/http"
)

type HelloServer struct {
}

var message string

func createMessage(str string, count int64) {
	message = str
	if count > 0 {
		var i int64
		count := int64(math.Log2(float64(count)))
		for i = 0; i < count; i++ {
			message = message + message
		}
	}
}

func (h *HelloServer) CreateMessage(ctx context.Context, rq *pb.MessageRequest) (*pb.HelloReply, error) {
	createMessage(rq.GetStr(), rq.GetCount())
	re := &pb.HelloReply{Message: message}
	return re, nil
}
func (h *HelloServer) GetMessage(ctx context.Context, rq *pb.NullMessageRequest) (*pb.HelloReply, error) {
	re := &pb.HelloReply{Message: message}
	return re, nil
}

type RestHelloServer struct {
}

func (r *RestHelloServer) CreateMessage(ctx *restful.Context) {
	requestBody := pb.MessageRequest{}
	err := ctx.ReadEntity(&requestBody)
	if err != nil {
		ctx.Write([]byte(err.Error()))
		return
	}
	createMessage(requestBody.Str, requestBody.Count)
	helloReply := pb.HelloReply{Message: message}
	ctx.WriteJSON(helloReply, "application/json")
	return
}

func (r *RestHelloServer) GetMessage(ctx *restful.Context) {
	helloReply := pb.HelloReply{Message: message}
	ctx.WriteJSON(helloReply, "application/json")
	return
}

func (r *RestHelloServer) URLPatterns() []restful.Route {
	return []restful.Route{
		{http.MethodPost, "/createmessage", "CreateMessage"},
		{http.MethodGet, "/getmessage", "GetMessage"},
	}
}
