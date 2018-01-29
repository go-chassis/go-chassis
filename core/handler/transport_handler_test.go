package handler_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/core/util/metadata"
	"github.com/ServiceComb/go-chassis/examples/schemas"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	serverOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport/tcp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var tr transport.Transport = tcp.NewTransport()

func TestTransportHandler_Handle(t *testing.T) {
	chassis.RegisterSchema("highway", &schemas.HelloServer{}, serverOption.WithSchemaID("HelloServer"))
	chassis.RegisterSchema("highway", &schemas.EmployServer{}, serverOption.WithSchemaID("EmployServer"), serverOption.WithMicroServiceName("Server"))
	t.Log("testing transport handler with highway protocol")
	p := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "communication/client")
	os.Setenv("CHASSIS_HOME", p)
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	var addrHighway = "127.0.0.1:4567"
	msName := "Server"
	schema := "schema2"
	config.Init()
	trServer := tcp.NewTransport()
	f, err := server.GetServerFunc("highway")
	assert.NoError(t, err)
	s := f(
		serverOption.Transport(trServer),
		serverOption.Address(addrHighway),
		serverOption.ChainName("default"))
	_, err = s.Register(&schemas.HelloServer{},
		serverOption.WithMicroServiceName(msName),
		serverOption.WithSchemaID(schema))

	assert.NoError(t, err)
	err = s.Start()
	assert.NoError(t, err)
	//dial
	c := &handler.Chain{}
	i := &invocation.Invocation{}
	i.Reply = &helloworld.HelloReply{}
	i.Ctx = metadata.NewContext(context.Background(), map[string]string{
		"X-User": "tianxiaoliang",
	})
	i.Protocol = "highway"
	i.Args = &helloworld.HelloRequest{Name: "peter"}

	i.Endpoint = addrHighway
	i.Protocol = "highway"
	i.SchemaID = schema
	i.MicroServiceName = msName
	h := &handler.TransportHandler{}
	c.Handlers = append(c.Handlers, h)

	var err2 error
	c.Next(i, func(r *invocation.InvocationResponse) error {
		log.Println("chain start")
		log.Println(r.Result)
		log.Println(r.Err)
		//assert.Empty(t, r.Err)
		//assert.NoError(t, r.Err)
		assert.Equal(t, nil, r.Result)
		err2 = r.Err
		return r.Err
	})

	var th *handler.TransportHandler = new(handler.TransportHandler)
	str := th.Name()
	assert.Equal(t, "transport", str)

}

func TestTransportHandler_HandleRest(t *testing.T) {
	t.Log("testing transport handler with rest protocol")
	p := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "communication/client")
	os.Setenv("CHASSIS_HOME", p)
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	config.Init()
	config.GlobalDefinition.Cse.Protocols = map[string]model.Protocol{
		"rest": {Listen: "0.0.0.0:2678", Advertise: "0.0.0.0:8888", WorkerNumber: 1},
	}

	l, err := tr.Listen("127.0.0.1:9992")
	if err != nil {
		t.Errorf("Unexpected listen err: %v", err)
	}
	defer l.Close()

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

	i.Protocol = "highway"
	i.Args = &helloworld.HelloRequest{Name: "peter"}

	i.Endpoint = "127.0.0.1:9992"
	i.Protocol = "rest"
	h := &handler.TransportHandler{}
	c.Handlers = append(c.Handlers, h)

	var err2 error
	c.Next(i, func(r *invocation.InvocationResponse) error {
		log.Println("chain start")
		log.Println(r.Result)
		log.Println(r.Err)
		//assert.Empty(t, r.Err)
		//assert.NoError(t, r.Err)
		assert.Equal(t, nil, r.Result)
		err2 = r.Err
		return r.Err
	})
	l.Close()

}
