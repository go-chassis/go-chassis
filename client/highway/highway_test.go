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
	_ "github.com/ServiceComb/go-chassis/server/highway"
	clientOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	serverOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport/tcp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var addrHighway = "127.0.0.1:8969"

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

func TestTcpClient_Call(t *testing.T) {
	t.Log("Testing Highway client functions")
	initEnv()
	msName := "Server"
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
	err = s.Start()
	assert.NoError(t, err)
	c := tcpClient.NewHighwayClient(clientOption.Transport(trClient))

	if err != nil {
		t.Errorf("Unexpected dial err: %v", err)
	}

	c.Init(clientOption.ContentType("application/protobuf"))
	arg := &helloworld.HelloRequest{
		Name: "peter",
	}
	reply := &helloworld.HelloReply{}

	name := c.String()
	log.Println("protocol name:", name)
	options := c.Options()
	log.Println("options are :", options)
	req := c.NewRequest(msName, schema, "SayHello", arg)
	log.Println("ms ", req.MicroServiceName, " send ", string(arg.Name))
	err = c.Call(context.TODO(), addrHighway, req, reply)
	log.Println("hellp reply", reply)
	assert.NoError(t, err)
	err = c.Call(context.TODO(), "127.0.0.1:5000", req, reply)
	assert.Error(t, err)
	err = c.Call(context.TODO(), "127.0.0.1:5000%err", req, reply)
	assert.Error(t, err)
}

func BenchmarkHighwayClient_Call(b *testing.B) {
	initEnv()
	msName := "Server"
	schema := "schema2"

	trServer := tcp.NewTransport()
	trClient := tcp.NewTransport()

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("highway")
	log.Println(err)
	s := f(
		serverOption.Transport(trServer),
		serverOption.Address(addrHighway),
		serverOption.ChainName("default"))

	_, err = s.Register(&schemas.HelloServer{},
		serverOption.WithMicroServiceName(msName),
		serverOption.WithSchemaID(schema))
	err = s.Start()
	c := tcpClient.NewHighwayClient(
		clientOption.Transport(trClient),
		clientOption.ContentType("application/protobuf"))

	arg := &helloworld.HelloRequest{
		Name: "peter",
	}
	reply := &helloworld.HelloReply{}

	req := c.NewRequest(msName, schema, "SayHello", arg)
	log.Println("ms ", req.MicroServiceName, " send ", string(arg.Name))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = c.Call(context.TODO(), addrHighway, req, reply)
		if err != nil {
			panic(err)
		}
	}
}

//请求Id溢出测试
//func TestInt(t *testing.T) {
//	var i int32
//
//	i = (1 << 31) - 2
//	fmt.Printf("int i:%d \n", i)
//	i=1+i
//	fmt.Printf("int i+1:%d \n", i)
//	if i == ((1 << 31) - 1) {
//		i=i%i
//	}
//	fmt.Printf("int i+2:%d \n", i)
//	i=1+i
//	fmt.Printf("int i+3:%d \n", i)
//}

//func Test(t *testing.T) {
//	//s := "highway://100.101.57.37:7070?login=true"
//	s := "100.101.57.37:7070?login=true"
//	fmt.Println("xxxx:", s)
//	u, err := url.Parse(s)
//
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(u.Scheme)
//	fmt.Println("xxxxxxxx",u.Host)
//
//}
