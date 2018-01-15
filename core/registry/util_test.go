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

	mapproto[common.ProtocolHighway] = model.Protocol{
		Listen:    "0.0.0.0:1",
		Advertise: "0.0.0.1:1",
	}
	strArr := registry.MakeEndpoints(mapproto)
	t.Log("making endpoints with listen and advertise addr, endpoint : ", strArr)
	assert.NotNil(t, strArr)
	assert.Equal(t, common.ProtocolHighway+"://"+mapproto[common.ProtocolHighway].Advertise, strArr[0])

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
