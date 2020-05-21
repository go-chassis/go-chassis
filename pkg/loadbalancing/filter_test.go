package loadbalancing_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/loadbalancing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	archaius.Init(archaius.WithMemorySource())
	config.ReadGlobalConfigFromArchaius()
}
func TestFilterAvailableZoneAffinity(t *testing.T) {

	var dc = new(registry.DataCenterInfo)
	config.GlobalDefinition.DataCenter.AvailableZone = "default-df-1"
	config.GlobalDefinition.DataCenter.Name = "default-df"
	config.GlobalDefinition.DataCenter.Region = "default-df"

	//same zone and same region case
	dc.AvailableZone = config.GlobalDefinition.DataCenter.AvailableZone
	dc.Region = config.GlobalDefinition.DataCenter.Name
	dc.Name = config.GlobalDefinition.DataCenter.Name
	testData := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]*registry.Endpoint{"rest": {
				false,
				"127.0.0.1:80",
			}},
			Metadata:       map[string]string{"key": "1"},
			DataCenterInfo: dc,
		},
		{
			EndpointsMap: map[string]*registry.Endpoint{"rest": {
				false,
				"127.0.0.1:80",
			}},
			Metadata: map[string]string{"key": "1"},
		},
	}
	loadbalancer.InstallFilter(loadbalancer.ZoneAware, loadbalancing.FilterAvailableZoneAffinity)
	instances := loadbalancer.Filters[loadbalancer.ZoneAware](testData, nil)
	assert.NotEqual(t, 0, len(instances))

	//out of region case
	dc.Region = "default"
	testData = []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]*registry.Endpoint{"rest": {
				false,
				"127.0.0.1:80",
			}},
			Metadata:       map[string]string{"key": "1"},
			DataCenterInfo: dc,
		},
	}

	instances = loadbalancing.FilterAvailableZoneAffinity(testData, nil)
	assert.NotEqual(t, 0, len(instances))

	//Same region but any available zone
	dc.Region = config.GlobalDefinition.DataCenter.Name
	dc.AvailableZone = "default-df-2"
	testData = []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]*registry.Endpoint{"rest": {
				false,
				"127.0.0.1:80",
			}},
			Metadata:       map[string]string{"key": "1"},
			DataCenterInfo: dc,
		},
	}

	instances = loadbalancing.FilterAvailableZoneAffinity(testData, nil)
	assert.NotEqual(t, 0, len(instances))

}
