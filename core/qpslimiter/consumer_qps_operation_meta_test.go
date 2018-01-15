package qpslimiter_test

import (
	//"fmt"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/qpslimiter"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitSchemaOperations(t *testing.T) {
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}

	opMeta := qpslimiter.InitSchemaOperations(i)
	t.Log("initializing schemaoperation from invocation object, OperationMeta = ", *opMeta)
	sName := opMeta.GetMicroServiceName()
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.service1", sName)

	schemaOpeartionName := opMeta.GetMicroServiceSchemaOpQualifiedName()
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.service1.schema1.SayHello", schemaOpeartionName)

	schemaName := opMeta.GetSchemaQualifiedName()
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.service1.schema1", schemaName)

}
func TestInitSchemaOperations4Mesher(t *testing.T) {
	i := &invocation.Invocation{
		SourceMicroService: "client:1.1:sock",
		MicroServiceName:   "service1",
		SchemaID:           "schema1",
		OperationID:        "SayHello",
		Args:               &helloworld.HelloRequest{Name: "peter"},
	}

	opMeta := qpslimiter.InitSchemaOperations(i)
	t.Log("initializing schemaoperation from invocation object with sourceMicroserviceName, OperationMeta = ", *opMeta)
	sName := opMeta.GetMicroServiceName()
	assert.Equal(t, "cse.flowcontrol.client:1.1:sock.Consumer.qps.limit.service1", sName)

	schemaOpeartionName := opMeta.GetMicroServiceSchemaOpQualifiedName()
	assert.Equal(t, "cse.flowcontrol.client:1.1:sock.Consumer.qps.limit.service1.schema1.SayHello", schemaOpeartionName)

	schemaName := opMeta.GetSchemaQualifiedName()
	assert.Equal(t, "cse.flowcontrol.client:1.1:sock.Consumer.qps.limit.service1.schema1", schemaName)

}
