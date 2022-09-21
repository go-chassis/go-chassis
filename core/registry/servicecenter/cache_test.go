package servicecenter_test

import (
	"os"
	"testing"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/sc-client"

	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	utiltags "github.com/go-chassis/go-chassis/v2/pkg/util/tags"
	_ "github.com/go-chassis/go-chassis/v2/security/cipher/plugins/plain"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/rand"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "DEBUG",
	})
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.registry.address", "http://127.0.0.1:30100")
	archaius.Set("servicecomb.registry.autoSchemaIndex", true)
	archaius.Set("servicecomb.service.name", "Server")
	archaius.Set("servicecomb.service.hostname", "localhost")
	os.Setenv("HTTP_DEBUG", "1")
	defer os.Unsetenv("HTTP_DEBUG")
	config.ReadGlobalConfigFromArchaius()
}
func TestCacheManager_AutoSync(t *testing.T) {
	registry.Enable()
	registry.DoRegister()
	t.Log("持有id", runtime.ServiceID)
	t.Log("同步sc节点")
	time.Sleep(time.Second * 1)

	microservice := &registry.MicroService{
		AppID:       "default",
		ServiceName: "Server",
		Version:     "0.1",
		Status:      sc.MicorserviceUp,
		Level:       "FRONT",
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]*registry.Endpoint{
			"rest": {
				false,
				"10.146.207.197:5080",
			},
		},
		InstanceID: "event1",
		HostName:   "event_test",
		Status:     sc.MSInstanceUP,
	}
	sid, instanceID, err := registry.DefaultRegistrator.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)
	assert.Equal(t, "event1", instanceID)
	time.Sleep(time.Second * 1)
	tags := utiltags.NewDefaultTag("0.1", "default")

	t.Run("find instances, should has response", func(t *testing.T) {
		instances, err := registry.DefaultServiceDiscoveryService.FindMicroServiceInstances(sid, "Server", tags)
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
	})
	t.Run("find instances with service do not exists in service center", func(t *testing.T) {
		instances, err := registry.DefaultServiceDiscoveryService.FindMicroServiceInstances(sid, "NotExistServer", tags)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(instances))

		instances2, err2 := registry.DefaultServiceDiscoveryService.FindMicroServiceInstances(sid, "NotExistServer", tags)
		assert.NoError(t, err2)
		assert.Equal(t, 0, len(instances2))
	})
	var exist = false
	pro := make(map[string]string)
	pro["attr1"] = "b"
	err = registry.DefaultRegistrator.UpdateMicroServiceInstanceProperties(sid, "event1", pro)
	assert.NoError(t, err)
	if err != nil {
		exist = true
	}
	assert.False(t, exist)
	time.Sleep(time.Second * 1)
	t.Log("实例信息变化感知成功")
	t.Log("测试EVT_UPDATE操作")

	exist = false
	err = registry.DefaultRegistrator.UpdateMicroServiceInstanceStatus(sid, "event1", sc.MSIinstanceDown)
	assert.NoError(t, err)
	if err != nil {
		exist = true
	}
	assert.False(t, exist)
	time.Sleep(time.Second * 1)
	t.Log("实例状态变化感知成功")
	t.Log("测试EVT_DELETE操作")

	exist = false
	err = registry.DefaultRegistrator.UpdateMicroServiceInstanceStatus(sid, "event1", sc.MSInstanceUP)
	assert.NoError(t, err)
	if err != nil {
		exist = true
	}
	assert.False(t, exist)
	time.Sleep(time.Second * 1)
	t.Log("实例状态变化感知成功")
	t.Log("测试EVT_DELETE操作")

	exist = false
	err = registry.DefaultRegistrator.UnRegisterMicroServiceInstance(sid, "event1")
	assert.NoError(t, err)
	if err != nil {
		exist = true
	}
	assert.False(t, exist)
	time.Sleep(time.Second * 1)

	t.Log("删除实例感知成功")
	t.Log("测试EVT_DELETE操作")

	t.Log("持有id", runtime.ServiceID)
	t.Log("watch测试完成")

}

func TestServiceDiscovery_AutoSync(t *testing.T) {
	v1, _ := version.NewVersion("1.2.1")
	v2, _ := version.NewVersion("1.10.1")
	v3, _ := version.NewVersion("1.21.1")
	v4, err := version.NewVersion("0.0.0")
	v5, err := version.NewVersion("0.0.1")
	assert.NoError(t, err)
	assert.True(t, v1.LessThan(v2))
	assert.True(t, v1.LessThan(v3))
	assert.True(t, v4.LessThan(v1))
	assert.True(t, v4.LessThan(v5))
}

func TestCacheManager_MakeSchemaIndex(t *testing.T) {
	/*
		1. Init Config
		2. Add key for autoIndex
		3. Start a microservice
		4. Check the status of Cache
	*/

	config.GlobalDefinition.ServiceComb.Registry.RefreshInterval = "1"
	registry.Enable()
	registry.DoRegister()
	time.Sleep(time.Second * 1)

	microservice := &registry.MicroService{
		AppID:       "default",
		ServiceName: "AutoIndexServer",
		Version:     "0.1",
		Status:      sc.MicorserviceUp,
		Level:       "FRONT",
	}
	sid, _ := registry.DefaultRegistrator.RegisterService(microservice)
	schemaName := rand.String(10)
	schemaInfoString := "swagger: \"2.0\"\ninfo:\n  version: \"1.0.0\"\n  title: \"swagger definition for org.apache.servicecomb.samples.demo.sc.ClientApi\"\n  x-java-interface: \"cse.gen.huaweidemo.DemoClient2.hellosc." + schemaName + "\"\nbasePath: \"/\"\nconsumes:\n- \"application/json\"\nproduces:\n- \"application/json\"\npaths:\n  /sayhello:\n    get:\n      operationId: \"sayHello\"\n      parameters: []\n      responses:\n        200:\n          description: \"response of 200\"\n          schema:\n            type: \"string\"\n"
	registry.DefaultRegistrator.AddSchemas(sid, schemaName, schemaInfoString)
	registry.DefaultServiceDiscoveryService.AutoSync()
	time.Sleep(time.Second * 3)
	interfaceExistInCache := false
	if registry.SchemaInterfaceIndexedCache.ItemCount() >= 1 {
		interfaceExistInCache = true
	}
	assert.True(t, interfaceExistInCache)
}
