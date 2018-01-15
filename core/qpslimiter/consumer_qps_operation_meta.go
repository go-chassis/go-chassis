package qpslimiter

import (
	"github.com/ServiceComb/go-chassis/core/invocation"
	"strings"
)

// OperationMeta operation meta struct
type OperationMeta struct {
	MicroServiceName       string
	SchemaQualifiedName    string
	OperationQualifiedName string
}

// GetSpecificKey get specific key
func GetSpecificKey(sourceName, serviceType, serviceName, schemaID, OperationID string) string {
	var cmd = "cse.flowcontrol"
	//for mesher to govern
	if sourceName != "" {
		cmd = strings.Join([]string{cmd, sourceName, serviceType, "qps.limit"}, ".")
	} else {
		cmd = strings.Join([]string{cmd, serviceType, "qps.limit"}, ".")
	}
	if serviceName != "" {
		cmd = strings.Join([]string{cmd, serviceName}, ".")
	}
	if schemaID != "" {
		cmd = strings.Join([]string{cmd, schemaID}, ".")
	}
	if OperationID != "" {
		cmd = strings.Join([]string{cmd, OperationID}, ".")
	}
	return cmd

}

// InitSchemaOperations initialize schema operations
func InitSchemaOperations(i *invocation.Invocation) *OperationMeta {
	opMeta := new(OperationMeta)

	opMeta.MicroServiceName = GetSpecificKey(i.SourceMicroService, "Consumer", i.MicroServiceName, "", "")
	opMeta.SchemaQualifiedName = GetSpecificKey(i.SourceMicroService, "Consumer", i.MicroServiceName, i.SchemaID, "")
	opMeta.OperationQualifiedName = GetSpecificKey(i.SourceMicroService, "Consumer", i.MicroServiceName, i.SchemaID, i.OperationID)

	//for mesher server side rate limit
	//as a proxy,mesher handler request from  instances that belong to different ms
	return opMeta
}

// GetSchemaQualifiedName get schema qualified name
func (op *OperationMeta) GetSchemaQualifiedName() string {
	return op.SchemaQualifiedName
}

// GetMicroServiceSchemaOpQualifiedName get micro-service schema operation qualified name
func (op *OperationMeta) GetMicroServiceSchemaOpQualifiedName() string {
	return op.OperationQualifiedName
}

// GetMicroServiceName get micro-service name
func (op *OperationMeta) GetMicroServiceName() string {
	return op.MicroServiceName
}
