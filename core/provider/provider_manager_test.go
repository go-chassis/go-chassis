package provider_test

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/provider"
	"github.com/ServiceComb/go-chassis/examples/schemas"
	pb "github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddProvider(t *testing.T) {
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	p := provider.RegisterProvider("123", "service1")
	assert.Nil(t, p)
	p = provider.RegisterProvider("default", "service1")
	assert.NotNil(t, p)
	schema := "schema2"
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
	_, err = provider.GetOperation("service1", "schema2", "SayHello")
	assert.NoError(t, err)
}

func TestAddCustomProvider(t *testing.T) {
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	provider.RegisterCustomProvider("test", provider.NewProvider("test"))
	provider.RegisterCustomProvider("test", provider.NewProvider("test"))
}

func TestGetProvider(t *testing.T) {
	config.Init()
	_, err := provider.GetProvider("fake")
	assert.Error(t, err)
	_, err = provider.GetProvider("service1")
	assert.NoError(t, err)
}

func TestGetOperation(t *testing.T) {
	config.Init()
	_, err := provider.GetOperation("service1", "schema2", "SayHello")
	assert.NoError(t, err)
	_, err = provider.GetOperation("service1", "schemsa2", "SayHello")
	assert.Error(t, err)
	_, err = provider.GetOperation("notexistingservice", "schemsa2", "SayHello")
	assert.Error(t, err)

}

func TestRegisterSchema(t *testing.T) {
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	p := provider.RegisterProvider("default", "service1")
	assert.NotNil(t, p)
	err := provider.RegisterSchemaWithName("service1", "schema1", &schemas.HelloServer{})
	assert.NoError(t, err)
	err = provider.RegisterSchemaWithName("notexistingservice", "schema1", &schemas.HelloServer{})
	assert.Error(t, err)
	_, err = provider.RegisterSchema("HelloServer", &schemas.HelloServer{})
	assert.Error(t, err)
	_, err = provider.RegisterSchema("service1", &schemas.HelloServer{})
	assert.NoError(t, err)
}

func BenchmarkGetOperation(b *testing.B) {
	config.Init()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.GetOperation("service1", "schema1", "SayHello")
	}
}

func BenchmarkGetProvider(b *testing.B) {
	config.Init()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = provider.GetProvider("service1")
	}
}
