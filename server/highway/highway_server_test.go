package highway_test

// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/client/highway"
	"github.com/ServiceComb/go-chassis/core/client"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/examples/schemas"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"

	"github.com/stretchr/testify/assert"
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

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain
	config.SelfServiceName = msName
	f, err := server.GetServerFunc("highway")
	assert.NoError(t, err)
	s := f(server.Options{
		Address:   addrHighway,
		ChainName: "default",
	})

	_, err = s.Register(&schemas.HelloServer{},
		server.WithSchemaID(schema))
	assert.NoError(t, err)

	err = s.Start()
	assert.NoError(t, err)

	name := s.String()
	log.Println("server protocol name:", name)
	c := highway.NewHighwayClient(client.Options{})

	if err != nil {
		t.Errorf("Unexpected dial err: %v", err)
	}

	arg := &helloworld.HelloRequest{
		Name: "peter",
	}
	reply := &helloworld.HelloReply{}

	name = c.String()
	log.Println("protocol name:", name)
	req := client.NewRequest(msName, schema, "SayHello", arg)
	log.Println("ms ", req.MicroServiceName, " send ", string(arg.Name))
	err = c.Call(context.TODO(), addrHighway, req, reply)
	log.Println("hello reply", reply)
	assert.NoError(t, err)

	//error scenario : Server2 microservice not exist
	req = client.NewRequest("Server2", schema, "SayHello", arg)
	log.Println("ms ", req.MicroServiceName, " send ", string(arg.Name))
	err = c.Call(context.TODO(), addrHighway, req, reply)
	log.Println("error is:", err)
	assert.Error(t, err)

	var ms = new(model.MicroserviceCfg)
	ms.Provider = "provider"
	config.MicroserviceDefinition = ms
	_, err = s.Register(&schemas.HelloServer{})
	assert.NoError(t, err)

	err = s.Stop()
	assert.NoError(t, err)
}
