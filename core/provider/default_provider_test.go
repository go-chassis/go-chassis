package provider_test

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/provider"
	"github.com/ServiceComb/go-chassis/examples/schemas"
	pb "github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"testing"
)

type HelloServerfk1 struct{}
type HelloServerfk2 struct{}

func (s *HelloServerfk1) sayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Go Hello  " + in.Name}, nil
}

func TestRegister(t *testing.T) {
	t.Log("testing registeration of a schema")
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	config.Init()
	p := &provider.DefaultProvider{}
	err := p.RegisterName("schema1", &schemas.HelloServer{})
	assert.NoError(t, err)
}

func TestRegister2(t *testing.T) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	config.Init()
	provider := &provider.DefaultProvider{}
	sname, err := provider.Register(&schemas.HelloServer{})
	assert.NoError(t, err)
	assert.NotEqual(t, "", sname)
	assert.NotEqual(t, nil, sname)
	op, _ := provider.GetOperation(sname, "SayHello")
	assert.NotNil(t, op)
	method := op.Method()
	assert.NotEqual(t, nil, method)
	args := op.Args()
	assert.Equal(t, 2, len(args))
	replay := op.Reply()
	assert.Equal(t, 2, len(replay))
	schema1 := "fkschema"
	err = provider.RegisterName(schema1, &HelloServerfk1{})
	assert.Error(t, err)
	schema2 := "fkschema"
	err = provider.RegisterName(schema2, &HelloServerfk2{})
	assert.Error(t, err)

}

func TestProvider_Invoke(t *testing.T) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	p := &provider.DefaultProvider{}
	schema := "schema1"
	err := p.RegisterName(schema, &schemas.HelloServer{})
	assert.NoError(t, err)
	inv := &invocation.Invocation{
		SchemaID:    schema,
		OperationID: "SayHello",
		Args: &pb.HelloRequest{
			Name: "1",
		},
	}
	_, err = p.Invoke(inv)
	assert.NoError(t, err)

}

func TestDefaultProvider_GetOperation(t *testing.T) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	p := &provider.DefaultProvider{}
	schema := "schema1"
	err := p.RegisterName(schema, &schemas.HelloServer{})
	assert.NoError(t, err)
	op, _ := p.GetOperation(schema, "SayHello")
	assert.NotNil(t, op)
	op, _ = p.GetOperation(schema, "SayHelloasd")
	assert.Nil(t, op)
}

func TestDefaultProvider_Exist(t *testing.T) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	p := &provider.DefaultProvider{}
	schema := "schema1"
	err := p.RegisterName(schema, &schemas.HelloServer{})
	assert.NoError(t, err)
	e := p.Exist(schema, "SayHello")
	assert.Equal(t, true, e)
	e = p.Exist(schema, "SayHelloaasda")
	assert.Equal(t, false, e)
	e = p.Exist(schema+"123", "SayHelloaasda")
	assert.Equal(t, false, e)
}

func BenchmarkDefaultProvider_GetOperation(b *testing.B) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	p := &provider.DefaultProvider{}
	schema := "schema1"
	_ = p.RegisterName(schema, &schemas.HelloServer{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.GetOperation(schema, "SayHello")
	}
}
