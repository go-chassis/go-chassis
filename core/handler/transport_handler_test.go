package handler_test

import (
	"github.com/go-chassis/go-chassis/core/lager"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"

	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
func TestTransportHandler_HandleRest(t *testing.T) {
	t.Log("testing transport handler with rest protocol")
	p := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-chassis", "go-chassis", "examples", "communication/client")
	os.Setenv("CHASSIS_HOME", p)

	config.Init()
	config.GlobalDefinition.Cse.Protocols = map[string]model.Protocol{
		"rest": {Listen: "0.0.0.0:2678", Advertise: "0.0.0.0:8888", WorkerNumber: 1},
	}

	/*fn := func(sock transport.Socket) {
		defer sock.Close()

		for {
			//metadata := make(map[string]string)
			//metadata["requestID"] = "0"
			responseHeader, responseBody, _, ID, err := sock.Recv()
			if err != nil {
				return
			}
			log.Println("server receive", string(responseBody))
			if err := sock.Send(responseHeader, responseBody, nil, ID); err != nil {
				return
			}
			log.Println("server send", string(responseBody))
		}
	}

	done := make(chan bool)

	go func() {
		if err := l.Accept(fn); err != nil {
			log.Println(err)
			select {
			case <-done:
			default:
				t.Errorf("Unexpected accept err: %v", err)
			}
		}
	}()*/

	//dial
	c := &handler.Chain{}
	i := &invocation.Invocation{}
	i.Reply = &helloworld.HelloReply{}

	i.Endpoint = "127.0.0.1:9992"
	i.Protocol = "rest"
	h := &handler.TransportHandler{}
	c.Handlers = append(c.Handlers, h)

	c.Next(i, func(r *invocation.Response) error {
		log.Println("chain start")
		log.Println(r.Result)
		log.Println(r.Err)
		assert.Equal(t, nil, r.Result)
		return r.Err
	})

}
