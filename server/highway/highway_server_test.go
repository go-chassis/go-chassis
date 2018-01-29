package tcp_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ServiceComb/go-chassis"
	tcpClient "github.com/ServiceComb/go-chassis/client/highway"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/examples/schemas"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	clientOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	serverOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport/tcp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var addrHighway = "127.0.0.1:2399"

func initEnv() {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	log.Println(os.Getenv("CHASSIS_HOME"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition = &model.GlobalCfg{}
	chassis.Init()
}

func TestStart(t *testing.T) {
	t.Log("Testing highway server start function")
	initEnv()
	msName := "Server1"
	schema := "schema2"

	trServer := tcp.NewTransport()
	trClient := tcp.NewTransport()

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

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

	err = s.Init(serverOption.Name("Server"))
	assert.NoError(t, err)

	err = s.Start()
	assert.NoError(t, err)
	opt := s.Options()
	assert.NotEmpty(t, opt)

	name := s.String()
	log.Println("server protocol name:", name)
	c := tcpClient.NewHighwayClient(
		clientOption.Transport(trClient),
		clientOption.ContentType("application/protobuf"))

	if err != nil {
		t.Errorf("Unexpected dial err: %v", err)
	}

	arg := &helloworld.HelloRequest{
		Name: "peter",
	}
	reply := &helloworld.HelloReply{}

	name = c.String()
	log.Println("protocol name:", name)
	options := c.Options()
	log.Println("options are :", options)
	req := c.NewRequest(msName, schema, "SayHello", arg)
	log.Println("ms ", req.MicroServiceName, " send ", string(arg.Name))
	err = c.Call(context.TODO(), addrHighway, req, reply)
	log.Println("hello reply", reply)
	assert.NoError(t, err)

	//error scenario : Server2 microservice not exist
	req = c.NewRequest("Server2", schema, "SayHello", arg)
	log.Println("ms ", req.MicroServiceName, " send ", string(arg.Name))
	err = c.Call(context.TODO(), addrHighway, req, reply)
	log.Println("error is:", err)
	assert.Error(t, err)

	var ms *model.MicroserviceCfg = new(model.MicroserviceCfg)
	ms.Provider = "provider"
	config.MicroserviceDefinition = ms
	_, err = s.Register(&schemas.HelloServer{}, serverOption.WithMicroServiceName(msName))
	assert.NoError(t, err)

	err = s.Stop()
	assert.NoError(t, err)
}
