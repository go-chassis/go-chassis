package servicecenter_test

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	_ "github.com/ServiceComb/go-chassis/security/plugins/plain"
	"github.com/ServiceComb/go-sc-client/model"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCacheManager_AutoSync(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	t.Log("Test cache.go")
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	registry.Enable()
	registry.DoRegister()
	t.Log("持有id", config.SelfServiceID)
	t.Log("同步sc节点")
	time.Sleep(time.Second * 1)

	microservice := &registry.MicroService{
		AppID:       "default",
		ServiceName: "Server",
		Version:     "0.1",
		Status:      model.MicorserviceUp,
		Level:       "FRONT",
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]string{"rest": "10.146.207.197:8080"},
		InstanceID:   "event1",
		HostName:     "event_test",
		Status:       model.MSInstanceUP,
		Environment:  common.EnvValueProd,
	}
	sid, instanceID, err := registry.RegistryService.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)
	assert.Equal(t, "event1", instanceID)
	time.Sleep(time.Second * 1)
	instances, err := registry.RegistryService.FindMicroServiceInstances(sid, "default", "Server", "0.1")
	assert.NotZero(t, len(instances))
	assert.NoError(t, err)
	var ok = false
	for _, ins := range instances {
		t.Log(ins.InstanceID)
		if ins.InstanceID == "event1" {
			ok = true
			break
		}
	}
	assert.True(t, ok)
	t.Log("新增实例感知成功")
	t.Log("测试EVT_CREATE操作")

	var exist = false
	pro := make(map[string]string)
	pro["attr1"] = "b"
	err = registry.RegistryService.UpdateMicroServiceInstanceProperties(sid, "event1", pro)
	assert.NoError(t, err)
	if err != nil {
		exist = true
	}
	assert.False(t, exist)
	time.Sleep(time.Second * 1)
	t.Log("实例信息变化感知成功")
	t.Log("测试EVT_UPDATE操作")

	exist = false
	err = registry.RegistryService.UpdateMicroServiceInstanceStatus(sid, "event1", model.MSIinstanceDown)
	assert.NoError(t, err)
	if err != nil {
		exist = true
	}
	assert.False(t, exist)
	time.Sleep(time.Second * 1)
	t.Log("实例状态变化感知成功")
	t.Log("测试EVT_DELETE操作")

	exist = false
	err = registry.RegistryService.UpdateMicroServiceInstanceStatus(sid, "event1", model.MSInstanceUP)
	assert.NoError(t, err)
	if err != nil {
		exist = true
	}
	assert.False(t, exist)
	time.Sleep(time.Second * 1)
	t.Log("实例状态变化感知成功")
	t.Log("测试EVT_DELETE操作")

	exist = false
	err = registry.RegistryService.UnregisterMicroServiceInstance(sid, "event1")
	assert.NoError(t, err)
	if err != nil {
		exist = true
	}
	assert.False(t, exist)
	time.Sleep(time.Second * 1)

	t.Log("删除实例感知成功")
	t.Log("测试EVT_DELETE操作")

	t.Log("持有id", config.SelfServiceID)
	t.Log("watch测试完成")

}
