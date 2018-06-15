package registry_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/util/iputil"
	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	var eps = []string{"https://127.0.0.1", "http://0.0.0.0"}
	mp, str := registry.GetProtocolMap(eps)
	assert.Equal(t, "0.0.0.0", mp["http"])
	assert.Equal(t, "127.0.0.1", mp["https"])
	assert.Equal(t, "http", str)

	var mapproto map[string]model.Protocol = make(map[string]model.Protocol)
	var mapprotoRest map[string]model.Protocol = make(map[string]model.Protocol)

	mapproto[common.ProtocolHighway] = model.Protocol{
		Listen:    "0.0.0.0:1",
		Advertise: "0.0.0.1:1",
	}
	strArr := registry.MakeEndpoints(mapproto)
	t.Log("making endpoints with listen and advertise addr, endpoint : ", strArr)
	assert.NotNil(t, strArr)
	assert.Equal(t, common.ProtocolHighway+"://"+mapproto[common.ProtocolHighway].Advertise, strArr[0])

	//Advertise address are given in the protocol map for highway
	protocolArr := registry.MakeEndpointMap(mapproto)
	t.Log("making endpoints with listen and advertise addr, endpoint : ", protocolArr)
	assert.NotNil(t, protocolArr)
	assert.Equal(t, common.ProtocolHighway+":"+mapproto[common.ProtocolHighway].Advertise, common.ProtocolHighway+":"+protocolArr[common.ProtocolHighway])

	//Advertise address are given in the protocol map for rest
	mapprotoRest[common.ProtocolRest] = model.Protocol{
		Listen:    "0.0.0.2:1",
		Advertise: "0.0.0.1:1",
	}

	protocolArrRest := registry.MakeEndpointMap(mapprotoRest)
	t.Log("making endpoints with listen and advertise addr, endpoint : ", protocolArrRest)
	assert.NotNil(t, protocolArrRest)
	assert.Equal(t, common.ProtocolRest+":"+mapprotoRest[common.ProtocolRest].Advertise, common.ProtocolRest+":"+protocolArrRest[common.ProtocolRest])

	// Advertise address are given in the protocol map for rest
	// and addr is loopback ip. so it should return empty response
	mapprotoRest[common.ProtocolRest] = model.Protocol{
		Listen:    "0.0.0.2:1",
		Advertise: "127.0.0.1:1",
	}

	protocolArrRest = registry.MakeEndpointMap(mapprotoRest)
	t.Log("making endpoints with listen and advertise addr, endpoint : ", protocolArrRest)
	assert.NotNil(t, protocolArrRest)
	assert.Equal(t, "127.0.0.1:1", protocolArrRest[common.ProtocolRest])

	// Advertise address are given in the protocol map for rest
	// and addr is IPV6 ip. so it should return empty response
	mapprotoRest[common.ProtocolRest] = model.Protocol{
		Listen:    "0.0.0.2:1",
		Advertise: "fe80::3436:b05c:350a:1ccd:1",
	}

	protocolArrRest = registry.MakeEndpointMap(mapprotoRest)
	t.Log("making endpoints with listen and advertise addr, endpoint : ", protocolArrRest)
	assert.NotNil(t, protocolArrRest)
	assert.Equal(t, "", protocolArrRest[common.ProtocolRest])

	// Advertise address is not given so based on the listen address it should choose the advertise addr.
	mapprotoRest[common.ProtocolRest] = model.Protocol{
		Listen: "0.0.0.2:1",
	}

	protocolArrRest = registry.MakeEndpointMap(mapprotoRest)
	t.Log("making endpoints with listen and advertise addr, endpoint : ", protocolArrRest)
	assert.NotNil(t, protocolArrRest)
	assert.Equal(t, common.ProtocolRest+":"+mapprotoRest[common.ProtocolRest].Listen, common.ProtocolRest+":"+protocolArrRest[common.ProtocolRest])

	// Advertise address is not given and listen addr is 0.0.0.0 so it should select the ip from IPV4 of eth.
	mapprotoRest[common.ProtocolRest] = model.Protocol{
		Listen: "0.0.0.0:1",
	}

	protocolArrRest = registry.MakeEndpointMap(mapprotoRest)
	t.Log("making endpoints with listen and advertise addr, endpoint : ", protocolArrRest)
	assert.NotNil(t, protocolArrRest)

	mapproto[common.ProtocolHighway] = model.Protocol{
		Listen: "0.0.0.0:1",
	}
	strArr = registry.MakeEndpoints(mapproto)
	t.Log("making endpoints with listen addr only, endpoint : ", strArr)
	assert.NotNil(t, strArr)
	assert.Equal(t, common.ProtocolHighway+"://"+mapproto[common.ProtocolHighway].Listen, strArr[0])

	mapproto[common.ProtocolHighway] = model.Protocol{
		Listen: "",
	}
	strArr = registry.MakeEndpoints(mapproto)
	t.Log("making endpoints without listen and advertise, endpoint : ", strArr)
	assert.NotNil(t, strArr)
	assert.Equal(t, common.ProtocolHighway+"://"+iputil.DefaultEndpoint4Protocol(common.ProtocolHighway), strArr[0])
}
func TestGetProtocolList(t *testing.T) {
	m := map[string]string{
		"rest": "1.1.1.1",
		"http": "1.1.1.1",
	}
	eps := registry.GetProtocolList(m)
	assert.Equal(t, 2, len(eps))
	t.Log(eps)
}

func TestURIs2Hosts(t *testing.T) {
	hosts, s, err := registry.URIs2Hosts([]string{"http://127.0.0.1:8080"})
	assert.NoError(t, err)
	assert.Equal(t, "http", s)
	assert.Equal(t, "127.0.0.1:8080", hosts[0])
}
