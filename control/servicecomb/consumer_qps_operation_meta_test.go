package servicecomb_test

import (
	"testing"

	arhcaiusPanel "github.com/go-chassis/go-chassis/control/servicecomb"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
)

func TestGetConsumerKey(t *testing.T) {
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}

	opMeta := arhcaiusPanel.GetConsumerKey(i.SourceMicroService, i.MicroServiceName, i.SchemaID, i.OperationID)
	t.Log("initializing schemaoperation from invocation object, ConsumerKeys = ", *opMeta)
	sName := opMeta.MicroServiceName
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.service1", sName)

	schemaOperationName := opMeta.OperationQualifiedName
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.service1.schema1.SayHello", schemaOperationName)

	schemaName := opMeta.SchemaQualifiedName
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.service1.schema1", schemaName)

}
func TestGetConsumerKey2(t *testing.T) {
	i := &invocation.Invocation{
		SourceMicroService: "client:1.1:sock",
		MicroServiceName:   "service1",
		SchemaID:           "schema1",
		OperationID:        "SayHello",
		Args:               &helloworld.HelloRequest{Name: "peter"},
	}

	opMeta := arhcaiusPanel.GetConsumerKey(i.SourceMicroService, i.MicroServiceName, i.SchemaID, i.OperationID)
	t.Log("initializing schemaoperation from invocation object with sourceMicroserviceName, ConsumerKeys = ", *opMeta)
	sName := opMeta.GetMicroServiceName()
	assert.Equal(t, "cse.flowcontrol.client:1.1:sock.Consumer.qps.limit.service1", sName)

	schemaOpeartionName := opMeta.GetMicroServiceSchemaOpQualifiedName()
	assert.Equal(t, "cse.flowcontrol.client:1.1:sock.Consumer.qps.limit.service1.schema1.SayHello", schemaOpeartionName)

	schemaName := opMeta.GetSchemaQualifiedName()
	assert.Equal(t, "cse.flowcontrol.client:1.1:sock.Consumer.qps.limit.service1.schema1", schemaName)

}
func TestGetQpsRateWithPriority(t *testing.T) {
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}
	opMeta := arhcaiusPanel.GetConsumerKey(i.SourceMicroService, i.MicroServiceName, i.SchemaID, i.OperationID)

	rate, key := arhcaiusPanel.GetQPSRateWithPriority(opMeta.OperationQualifiedName, opMeta.SchemaQualifiedName, opMeta.MicroServiceName)
	t.Log("rate is :", rate)
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.service1", key)

	i = &invocation.Invocation{
		MicroServiceName: "service1",
	}
	keys := arhcaiusPanel.GetProviderKey(i.SourceMicroService)
	rate, key = arhcaiusPanel.GetQPSRateWithPriority(keys.ServiceOriented, keys.Global)
	assert.Equal(t, "cse.flowcontrol.Provider.qps.global.limit", key)
}
