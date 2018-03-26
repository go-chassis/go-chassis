package loadbalancer_test

import (
	"testing"

	_ "github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/loadbalancer"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/stretchr/testify/assert"
	"sort"
	"time"
)

func TestWeightedResponseStrategy_Pick(t *testing.T) {
	config.Init()
	config.GetLoadBalancing().Strategy["name"] = loadbalancer.StrategyLatency
	instances := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "127.0.0.1:9090"},
		},
		{
			EndpointsMap: map[string]string{"rest": "10.0.0.3:8080", "highway": "10.0.0.3:9090"},
		},
	}

	loadbalancer.SetLatency(2*time.Second, "127.0.0.1:8080", "Server", "", "", common.ProtocolRest)
	loadbalancer.SetLatency(3*time.Second, "127.0.0.1:8080", "Server", "", "", common.ProtocolRest)
	loadbalancer.SetLatency(4*time.Second, "127.0.0.1:8080", "Server", "", "", common.ProtocolRest)

	loadbalancer.SetLatency(1*time.Second, "10.0.0.3:8080", "Server", "", "", common.ProtocolRest)
	loadbalancer.SetLatency(1*time.Second, "10.0.0.3:8080", "Server", "", "", common.ProtocolRest)
	loadbalancer.SetLatency(1*time.Second, "10.0.0.3:8080", "Server", "", "", common.ProtocolRest)

	loadbalancer.SetLatency(1*time.Second, "127.0.0.1:9090", "Server", "", "", common.ProtocolHighway)
	loadbalancer.SetLatency(3*time.Second, "127.0.0.1:9090", "Server", "", "", common.ProtocolHighway)
	loadbalancer.SetLatency(5*time.Second, "127.0.0.1:9090", "Server", "", "", common.ProtocolHighway)
	loadbalancer.SetLatency(9*time.Second, "127.0.0.1:9090", "Server", "", "", common.ProtocolHighway)

	loadbalancer.SetLatency(1*time.Second, "10.0.0.3:9090", "Server", "", "", common.ProtocolHighway)
	loadbalancer.SetLatency(1*time.Second, "10.0.0.3:9090", "Server", "", "", common.ProtocolHighway)
	loadbalancer.SetLatency(1*time.Second, "10.0.0.3:9090", "Server", "", "", common.ProtocolHighway)
	loadbalancer.SetLatency(1*time.Second, "10.0.0.3:9090", "Server", "", "", common.ProtocolHighway)
	loadbalancer.Enable()
	f, _ := loadbalancer.GetStrategyPlugin(loadbalancer.StrategyLatency)
	s := f()
	s.ReceiveData(instances, "Server", common.ProtocolRest, "")
	time.Sleep(31 * time.Second)
	var count int
	for i := 0; i < 100; i++ {
		instance, err := s.Pick()
		assert.NoError(t, err)
		if "10.0.0.3:8080" == instance.EndpointsMap["rest"] {
			count++
		}
	}
	t.Log(count)
	if !(count < 100) {
		t.Error(count)
	}
	s = f()
	s.ReceiveData(instances, "Server", common.ProtocolHighway, "")
	count = 0
	for i := 0; i < 1000; i++ {
		instance, err := s.Pick()
		assert.NoError(t, err)
		if "10.0.0.3:9090" == instance.EndpointsMap["highway"] {
			count++
		}

	}
	t.Log(count)
	if !(count < 1000) {
		t.Error(count)
	}
}
func TestCalculateAvgLatency(t *testing.T) {
	loadbalancer.SetLatency(2*time.Second, "127.0.0.1:3000", "Server1", "", "", common.ProtocolRest)
	loadbalancer.SetLatency(3*time.Second, "10.1.1.1.1:3000", "Server1", "", "", common.ProtocolRest)
	loadbalancer.SetLatency(1*time.Second, "10.0.0.1:3000", "Server1", "", "", common.ProtocolRest)
	loadbalancer.SetLatency(1*time.Second, "127.0.0.1:5000", "Server1", "", "", common.ProtocolHighway)
	loadbalancer.SetLatency(5*time.Second, "127.0.0.1:5000", "Server1", "", "", common.ProtocolHighway)
	loadbalancer.SetLatency(1*time.Second, "10.0.0.1:5000", "Server1", "", "", common.ProtocolHighway)
	loadbalancer.SetLatency(9*time.Second, "10.0.0.1:5000", "Server1", "", "", common.ProtocolHighway)
	loadbalancer.CalculateAvgLatency()
	for k, v := range loadbalancer.ProtocolStatsMap {
		if k == loadbalancer.BuildKey("Server1", "", "", common.ProtocolRest) {
			for _, s := range v {
				if s.Addr == "127.0.0.1:3000" {
					assert.Equal(t, time.Duration(2*time.Second), s.AvgLatency)
				}
				if s.Addr == "10.1.1.1:3000" {
					assert.Equal(t, time.Duration(3*time.Second), s.AvgLatency)
				}
			}

		}
		if k == loadbalancer.BuildKey("Server1", "", "", common.ProtocolHighway) {
			for _, s := range v {
				if s.Addr == "127.0.0.1:5000" {
					assert.Equal(t, time.Duration(3*time.Second), s.AvgLatency)
				}
				if s.Addr == "10.0.0.1:5000" {
					assert.Equal(t, time.Duration(5*time.Second), s.AvgLatency)
				}
			}

		}
	}
	loadbalancer.SortLatency()
	t.Log(loadbalancer.ProtocolStatsMap)
	assert.Equal(t, "127.0.0.1:5000", loadbalancer.ProtocolStatsMap[loadbalancer.BuildKey("Server1", "", "", common.ProtocolHighway)][0].Addr)
	assert.Equal(t, "10.0.0.1:3000", loadbalancer.ProtocolStatsMap[loadbalancer.BuildKey("Server1", "", "", common.ProtocolRest)][0].Addr)
}
func TestSortLatency(t *testing.T) {
	p1 := &loadbalancer.ProtocolStats{
		AvgLatency: 1 * time.Second,
	}
	p2 := &loadbalancer.ProtocolStats{
		AvgLatency: 4 * time.Second,
	}
	p3 := &loadbalancer.ProtocolStats{
		AvgLatency: 2 * time.Second,
	}
	p4 := &loadbalancer.ProtocolStats{
		AvgLatency: 6 * time.Second,
	}
	s := []*loadbalancer.ProtocolStats{p1, p2, p3, p4}
	sort.Sort(loadbalancer.ByDuration(s))
	t.Log(s[0].AvgLatency)
	t.Log(s[1].AvgLatency)
	t.Log(s[2].AvgLatency)
	t.Log(s[3].AvgLatency)
}
