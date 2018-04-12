package servicecenter

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-sc-client"
	"github.com/ServiceComb/go-sc-client/model"
	"gopkg.in/yaml.v2"
)

// parseSchemaContent parse schema content into SchemaContent structure
func parseSchemaContent(content []byte) (registry.SchemaContent, error) {
	var (
		err           error
		schema        = &registry.Schema{}
		schemaContent = &registry.SchemaContent{}
	)

	if err = yaml.Unmarshal(content, schema); err != nil {
		return *schemaContent, err
	}

	if err = yaml.Unmarshal([]byte(schema.Schema), schemaContent); err != nil {
		return *schemaContent, err
	}

	return *schemaContent, nil
}

// parseSchemaContent parse schema content into SchemaContent structure
func unmarshalSchemaContent(content []byte) (*registry.SchemaContent, error) {
	var (
		err           error
		schema        = &registry.Schema{}
		schemaContent = &registry.SchemaContent{}
	)

	if err = yaml.Unmarshal(content, schema); err != nil {
		return schemaContent, err
	}

	if err = yaml.Unmarshal([]byte(schema.Schema), schemaContent); err != nil {
		return schemaContent, err
	}

	return schemaContent, nil
}

// filterInstances filter instances
func filterInstances(providerInstances []*model.MicroServiceInstance) []*registry.MicroServiceInstance {
	instances := make([]*registry.MicroServiceInstance, 0)
	for _, ins := range providerInstances {
		if ins.Status != model.MSInstanceUP {
			continue
		}
		msi := ToMicroServiceInstance(ins)
		instances = append(instances, msi)
	}
	return instances
}

func closeClient(r *client.RegistryClient) error {
	err := r.Close()
	if err != nil {
		lager.Logger.Errorf(err, "Conn close failed.")
		return err
	}
	lager.Logger.Debugf("Conn close success.")
	return nil
}
