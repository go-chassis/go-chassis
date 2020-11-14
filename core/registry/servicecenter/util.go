package servicecenter

import (
	"fmt"
	scregistry "github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/go-chassis/v2/pkg/scclient"
	"github.com/go-chassis/go-chassis/v2/pkg/util/tags"
	"github.com/go-chassis/openlog"
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
func filterInstances(providerInstances []*scregistry.MicroServiceInstance) []*registry.MicroServiceInstance {
	instances := make([]*registry.MicroServiceInstance, 0)
	for _, ins := range providerInstances {
		if ins.Status != client.MSInstanceUP {
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
		openlog.Error(fmt.Sprintf("Conn close failed. err %s", err))
		return err
	}
	openlog.Debug("Conn close success.")
	return nil
}

func wrapTagsForServiceCenter(t utiltags.Tags) utiltags.Tags {
	if t.KV != nil {
		if v, ok := t.KV[common.BuildinTagVersion]; !ok || v == "" {
			t.KV[common.BuildinTagVersion] = common.LatestVersion
			t.Label += "|" + common.BuildinLabelVersion
		}
		if v, ok := t.KV[common.BuildinTagApp]; !ok || v == "" {
			t.KV[common.BuildinTagApp] = runtime.App
			t.Label += "|" + common.BuildinTagApp + ":" + runtime.App
		}
		return t
	}
	//if app and version is empty, need to find with latest version in same app
	return utiltags.NewDefaultTag(common.LatestVersion, runtime.App)
}

//GetCriteria generate batch find criteria from provider cache
func GetCriteria() []*scregistry.FindService {
	services := make([]*scregistry.FindService, 0)
	for _, service := range registry.GetProvidersFromCache() {
		services = append(services, &scregistry.FindService{
			Service: &scregistry.MicroServiceKey{
				ServiceName: service.ServiceName,
				Version:     service.Version,
				AppId:       service.AppID,
			},
		})
	}
	return services
}

//GetCriteriaByService generate batch find criteria from provider cache with same service name and different app
func GetCriteriaByService(sn string) []*scregistry.FindService {
	services := make([]*scregistry.FindService, 0)
	for _, service := range registry.GetProvidersFromCache() {
		if sn != service.ServiceName {
			continue
		}
		services = append(services, &scregistry.FindService{
			Service: &scregistry.MicroServiceKey{
				ServiceName: service.ServiceName,
				Version:     service.Version,
				AppId:       service.AppID,
			},
		})
	}
	return services
}
