package loadbalance_test

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestFilterAvailableZoneAffinity(t *testing.T) {
	t.Log("testing filter with specified region and zone ")
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	config.Init()
	t.Log(config.GlobalDefinition)
	var datacenter *registry.DataCenterInfo = new(registry.DataCenterInfo)
	if *config.GlobalDefinition.DataCenter == (model.DataCenterInfo{}) {
		config.GlobalDefinition.DataCenter.AvailableZone = "default-df-1"
		config.GlobalDefinition.DataCenter.Name = "default-df"
		config.GlobalDefinition.DataCenter.Region = "default-df"
	}

	//same zone and same region case
	datacenter.AvailableZone = config.GlobalDefinition.DataCenter.AvailableZone
	datacenter.Region = config.GlobalDefinition.DataCenter.Name
	datacenter.Name = config.GlobalDefinition.DataCenter.Name
	testData := []*registry.MicroServiceInstance{
		{
			EndpointsMap:   map[string]string{"rest": "127.0.0.1"},
			Metadata:       map[string]string{"key": "1"},
			DataCenterInfo: datacenter,
		},
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1"},
			Metadata:     map[string]string{"key": "1"},
		},
	}
	loadbalance.InstallFilter(loadbalance.ZoneAware, loadbalance.FilterAvailableZoneAffinity)
	instances := loadbalance.Filters[loadbalance.ZoneAware](testData)
	assert.NotEqual(t, 0, len(instances))

	//out of region case
	datacenter.Region = "default"
	testData = []*registry.MicroServiceInstance{
		{
			EndpointsMap:   map[string]string{"rest": "127.0.0.1"},
			Metadata:       map[string]string{"key": "1"},
			DataCenterInfo: datacenter,
		},
	}

	instances = loadbalance.FilterAvailableZoneAffinity(testData)
	assert.NotEqual(t, 0, len(instances))

	//Same region but any available zone
	datacenter.Region = config.GlobalDefinition.DataCenter.Name
	datacenter.AvailableZone = "default-df-2"
	testData = []*registry.MicroServiceInstance{
		{
			EndpointsMap:   map[string]string{"rest": "127.0.0.1"},
			Metadata:       map[string]string{"key": "1"},
			DataCenterInfo: datacenter,
		},
	}

	instances = loadbalance.FilterAvailableZoneAffinity(testData)
	assert.NotEqual(t, 0, len(instances))

}
func TestInstallFilter(t *testing.T) {

}
