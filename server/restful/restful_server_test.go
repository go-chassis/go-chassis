package restful_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/examples/schemas"
	serverOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport/tcp"
	"github.com/stretchr/testify/assert"
)

var addrHighway = "127.0.0.1:2399"
var addrHighway1 = "127.0.0.1:2330"

func initEnv() {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	log.Println(os.Getenv("CHASSIS_HOME"))
	os.Setenv("GO_CHASSIS_SWAGGERFILEPATH", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	log.Println(os.Getenv("GO_CHASSIS_SWAGGERFILEPATH"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition = &model.GlobalCfg{}
}

func TestRestStart(t *testing.T) {
	t.Log("Testing restful server start function")
	initEnv()
	msName := "Server1"
	schema := "schema1"

	trServer := tcp.NewTransport()
	//trClient := tcp.NewTransport()

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(
		serverOption.Transport(trServer),
		serverOption.Address(addrHighway),
		serverOption.ChainName("default"))

	_, err = s.Register(&schemas.RestFulHello{},
		serverOption.WithMicroServiceName(msName),
		serverOption.WithSchemaID(schema))
	assert.NoError(t, err)

	err = s.Init(serverOption.Name("Server1"))
	assert.NoError(t, err)

	err = s.Start()
	assert.NoError(t, err)

	name := s.String()
	assert.Equal(t, "rest", name)

	err = s.Stop()
	assert.NoError(t, err)
}

func TestRestStartFailure(t *testing.T) {
	t.Log("Testing restful server for start function failure")
	initEnv()
	msName := "Server2"
	schema := "schema2"

	trServer := tcp.NewTransport()
	//trClient := tcp.NewTransport()

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(
		serverOption.Transport(trServer),
		serverOption.Address(addrHighway),
		serverOption.ChainName("default"))

	_, err = s.Register(&schemas.HelloServer{},
		serverOption.WithMicroServiceName(msName),
		serverOption.WithSchemaID(schema))
	assert.Error(t, err)

	err = s.Init(serverOption.Name("Server2"))
	assert.NoError(t, err)

	err = s.Start()
	assert.NoError(t, err)

	name := s.String()
	assert.Equal(t, "rest", name)

	err = s.Stop()
	assert.NoError(t, err)
}
