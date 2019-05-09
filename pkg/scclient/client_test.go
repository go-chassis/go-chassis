package client_test

import (
	"github.com/go-chassis/go-chassis/pkg/scclient"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/scclient/proto"
	"github.com/go-chassis/paas-lager"
	"github.com/go-mesh/openlogging"
	"os"
	"time"
)

func init() {
	log.Init(log.Config{
		LoggerLevel:   "DEBUG",
		EnableRsyslog: false,
		LogFormatText: true,
		Writers:       []string{"stdout"},
	})
	l := log.NewLogger("test")
	openlogging.SetLogger(l)
}
func TestLoadbalance(t *testing.T) {

	t.Log("Testing Round robin function")
	var sArr []string

	sArr = append(sArr, "s1")
	sArr = append(sArr, "s2")

	next := client.RoundRobin(sArr)
	_, err := next()
	assert.NoError(t, err)
}

func TestLoadbalanceEmpty(t *testing.T) {
	t.Log("Testing Round robin with empty endpoint arrays")
	var sArrEmpty []string

	next := client.RoundRobin(sArrEmpty)
	_, err := next()
	assert.Error(t, err)

}

func TestClientInitializeHttpErr(t *testing.T) {
	t.Log("Testing for HTTPDo function with errors")

	hostname, err := os.Hostname()
	if err != nil {
		openlogging.GetLogger().Error("Get hostname failed.")
		return
	}
	microServiceInstance := &proto.MicroServiceInstance{
		Endpoints: []string{"rest://127.0.0.1:3000"},
		HostName:  hostname,
		Status:    client.MSInstanceUP,
	}

	registryClient := &client.RegistryClient{}

	err = registryClient.Initialize(
		client.Options{
			Addrs: []string{"127.0.0.1:30100"},
		})
	assert.NoError(t, err)

	err = registryClient.SyncEndpoints()
	assert.NoError(t, err)

	httpHeader := registryClient.GetDefaultHeaders()
	assert.NotEmpty(t, httpHeader)

	resp, err := registryClient.HTTPDo("GET", "fakeRawUrl", httpHeader, []byte("fakeBody"))
	assert.Empty(t, resp)
	assert.Error(t, err)

	MSList, err := registryClient.GetAllMicroServices()
	t.Log(MSList)
	assert.NotEmpty(t, MSList)
	assert.NoError(t, err)

	f1 := func(*client.MicroServiceInstanceChangedEvent) {}
	err = registryClient.WatchMicroService(MSList[0].ServiceId, f1)
	assert.NoError(t, err)

	var ms = new(proto.MicroService)
	var msdepreq = new(client.MircroServiceDependencyRequest)
	var msdepArr []*client.MicroServiceDependency
	var msdep1 = new(client.MicroServiceDependency)
	var msdep2 = new(client.MicroServiceDependency)
	var dep = new(client.DependencyMicroService)
	var m = make(map[string]string)

	m["abc"] = "abc"
	m["def"] = "def"

	dep.AppID = "AppId"

	msdep1.Consumer = dep
	msdep2.Consumer = dep

	msdepArr = append(msdepArr, msdep1)
	msdepArr = append(msdepArr, msdep2)

	ms.AppId = MSList[0].AppId
	ms.ServiceName = MSList[0].ServiceName
	ms.Version = MSList[0].Version
	ms.Environment = MSList[0].Environment
	ms.Properties = m

	msdepreq.Dependencies = msdepArr
	s1, err := registryClient.RegisterMicroServiceInstance(microServiceInstance)
	assert.Empty(t, s1)
	assert.Error(t, err)

	s1, err = registryClient.RegisterMicroServiceInstance(nil)
	assert.Empty(t, s1)
	assert.Error(t, err)

	msArr, err := registryClient.GetMicroServiceInstances("fakeConsumerID", "fakeProviderID")
	assert.Empty(t, msArr)
	assert.Error(t, err)

	msArr, err = registryClient.Health()
	assert.NotEmpty(t, msArr)
	assert.NoError(t, err)

	b, err := registryClient.UpdateMicroServiceProperties(MSList[0].ServiceId, ms)
	assert.Equal(t, true, b)
	assert.NoError(t, err)

	f1 = func(*client.MicroServiceInstanceChangedEvent) {}
	err = registryClient.WatchMicroService(MSList[0].ServiceId, f1)
	assert.NoError(t, err)

	f1 = func(*client.MicroServiceInstanceChangedEvent) {}
	err = registryClient.WatchMicroService("", f1)
	assert.Error(t, err)

	f1 = func(*client.MicroServiceInstanceChangedEvent) {}
	err = registryClient.WatchMicroService(MSList[0].ServiceId, nil)
	assert.NoError(t, err)

	str, err := registryClient.RegisterService(ms)
	assert.NotEmpty(t, str)
	assert.NoError(t, err)

	str, err = registryClient.RegisterService(nil)
	assert.Empty(t, str)
	assert.Error(t, err)
	t.Run("register service with name only", func(t *testing.T) {
		sid, err := registryClient.RegisterService(&proto.MicroService{
			ServiceName: "simpleService",
		})
		assert.NotEmpty(t, sid)
		assert.NoError(t, err)
		s, err := registryClient.GetMicroService(sid)
		assert.NoError(t, err)
		t.Log(s)
		assert.Equal(t, "0.0.1", s.Version)
		assert.Equal(t, "default", s.AppId)
		ok, err := registryClient.UnregisterMicroService(sid)
		assert.NoError(t, err)
		assert.True(t, ok)
		s, err = registryClient.GetMicroService(sid)
		assert.Nil(t, s)
	})
	t.Run("register service with invalid name", func(t *testing.T) {
		_, err := registryClient.RegisterService(&proto.MicroService{
			ServiceName: "simple&Service",
		})
		t.Log(err)
		assert.Error(t, err)
	})
	t.Run("get all apps", func(t *testing.T) {
		apps, err := registryClient.GetAllApplications()
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(apps))
		t.Log(apps)

	})
	ms1, err := registryClient.GetProviders("fakeconsumer")
	assert.Empty(t, ms1)
	assert.Error(t, err)

	err = registryClient.AddDependencies(msdepreq)
	assert.Error(t, err)

	err = registryClient.AddDependencies(nil)
	assert.Error(t, err)

	getms1, err := registryClient.GetMicroService(MSList[0].ServiceId)
	assert.NotEmpty(t, getms1)
	assert.NoError(t, err)

	getms2, err := registryClient.FindMicroServiceInstances("abcd", MSList[0].AppId, MSList[0].ServiceName, MSList[0].Version)
	assert.Empty(t, getms2)
	assert.Error(t, err)

	getmsstr, err := registryClient.GetMicroServiceID(MSList[0].AppId, MSList[0].ServiceName, MSList[0].Version, MSList[0].Environment)
	assert.NotEmpty(t, getmsstr)
	assert.NoError(t, err)

	getmsstr, err = registryClient.GetMicroServiceID(MSList[0].AppId, "Server112", MSList[0].Version, "")
	assert.Empty(t, getmsstr)
	//assert.Error(t, err)

	ms.Properties = nil
	b, err = registryClient.UpdateMicroServiceProperties(MSList[0].ServiceId, ms)
	assert.Equal(t, false, b)
	assert.Error(t, err)

	err = registryClient.AddSchemas("", "schema", "schema")
	assert.Error(t, err)

	b, err = registryClient.Heartbeat(MSList[0].ServiceId, "")
	assert.Equal(t, false, b)
	assert.Error(t, err)

	b, err = registryClient.UpdateMicroServiceInstanceStatus(MSList[0].ServiceId, "", MSList[0].Status)
	assert.Equal(t, false, b)
	assert.Error(t, err)

	b, err = registryClient.UnregisterMicroService("")
	assert.Equal(t, false, b)
	assert.Error(t, err)
	services, err := registryClient.GetAllResources("instances")
	assert.NotZero(t, len(services))
	assert.NoError(t, err)
	err = registryClient.Close()
	assert.NoError(t, err)

}
func TestRegistryClient_FindMicroServiceInstances(t *testing.T) {
	lager.Initialize("", "DEBUG", "",
		"size", true, 1, 10, 7)

	hostname, err := os.Hostname()
	if err != nil {
		openlogging.GetLogger().Error("Get hostname failed.")
		return
	}
	ms := &proto.MicroService{
		ServiceName: "scUTServer",
		AppId:       "default",
		Version:     "0.0.1",
		Schemas:     []string{"schema"},
	}
	var sid string
	registryClient := &client.RegistryClient{}

	err = registryClient.Initialize(
		client.Options{
			Addrs: []string{"127.0.0.1:30100"},
		})
	assert.NoError(t, err)
	sid, err = registryClient.RegisterService(ms)

	if err == client.ErrMicroServiceExists {
		sid, err = registryClient.GetMicroServiceID("default", "scUTServer", "0.0.1", "")
		assert.NoError(t, err)
		assert.NotNil(t, sid)
	}

	err = registryClient.AddSchemas(ms.ServiceId, "schema", "schema")
	assert.NoError(t, err)

	microServiceInstance := &proto.MicroServiceInstance{
		ServiceId: sid,
		Endpoints: []string{"rest://127.0.0.1:3000"},
		HostName:  hostname,
		Status:    client.MSInstanceUP,
	}
	t.Run("register instance and update props, should success", func(t *testing.T) {
		iid, err := registryClient.RegisterMicroServiceInstance(microServiceInstance)
		assert.NoError(t, err)
		assert.NotNil(t, iid)
		microServiceInstance.Properties = map[string]string{
			"project": "x"}
		ok, err := registryClient.UpdateMicroServiceInstanceProperties(microServiceInstance.ServiceId,
			iid, microServiceInstance)
		assert.True(t, ok)
		assert.NoError(t, err)
		instances, err := registryClient.FindMicroServiceInstances(microServiceInstance.ServiceId,
			"default",
			"scUTServer", "0.0.1")
		assert.NoError(t, err)
		assert.Equal(t, "x", instances[0].Properties["project"])
	})

	t.Log("find again, should get ErrNotModified")
	_, err = registryClient.FindMicroServiceInstances(sid, "default", "scUTServer", "0.0.1")
	assert.Equal(t, client.ErrNotModified, err)

	t.Log("find again without revision, should get nil error")
	_, err = registryClient.FindMicroServiceInstances(sid, "default", "scUTServer", "0.0.1", client.WithoutRevision())
	assert.NoError(t, err)

	t.Log("register new and find")
	microServiceInstance2 := &proto.MicroServiceInstance{
		ServiceId: sid,
		Endpoints: []string{"rest://127.0.0.1:3001"},
		HostName:  hostname + "1",
		Status:    client.MSInstanceUP,
	}
	_, err = registryClient.RegisterMicroServiceInstance(microServiceInstance2)
	time.Sleep(3 * time.Second)
	_, err = registryClient.FindMicroServiceInstances(sid, "default", "scUTServer", "0.0.1")
	assert.NoError(t, err)

	t.Log("after reset")
	registryClient.ResetRevision()
	_, err = registryClient.FindMicroServiceInstances(sid, "default", "scUTServer", "0.0.1")
	assert.NoError(t, err)
	_, err = registryClient.FindMicroServiceInstances(sid, "default", "scUTServer", "0.0.1")
	assert.Equal(t, client.ErrNotModified, err)

	_, err = registryClient.FindMicroServiceInstances(sid, "AppIdNotExists", "ServerNotExists", "0.0.1")
	assert.Equal(t, client.ErrMicroServiceNotExists, err)

	f := &proto.FindService{
		Service: &proto.MicroServiceKey{
			ServiceName: "scUTServer",
			AppId:       "default",
			Version:     "0.0.1",
		},
	}
	fs := []*proto.FindService{f}
	instances, err := registryClient.BatchFindInstances(sid, fs)
	t.Log(instances)
	assert.NoError(t, err)

	f1 := &proto.FindService{
		Service: &proto.MicroServiceKey{
			ServiceName: "empty",
			AppId:       "default",
			Version:     "0.0.1",
		},
	}
	fs = []*proto.FindService{f1}
	instances, err = registryClient.BatchFindInstances(sid, fs)
	t.Log(instances)
	assert.NoError(t, err)

	f2 := &proto.FindService{
		Service: &proto.MicroServiceKey{
			ServiceName: "empty",
			AppId:       "default",
			Version:     "latest",
		},
	}
	fs = []*proto.FindService{f}
	instances, err = registryClient.BatchFindInstances(sid, fs)
	t.Log(instances)
	assert.NoError(t, err)

	fs = []*proto.FindService{f2, f}
	instances, err = registryClient.BatchFindInstances(sid, fs)
	t.Log(instances)
	assert.NoError(t, err)

	fs = []*proto.FindService{}
	instances, err = registryClient.BatchFindInstances(sid, fs)
	assert.Equal(t, client.ErrEmptyCriteria, err)
}
func TestRegistryClient_GetDefaultHeaders(t *testing.T) {
	registryClient := &client.RegistryClient{}

	err := registryClient.Initialize(
		client.Options{
			Addrs:        []string{"127.0.0.1:30100"},
			ConfigTenant: "go-sc-tenant",
		})
	assert.Nil(t, err)

	header := registryClient.GetDefaultHeaders()
	tenant := header.Get(client.TenantHeader)
	assert.Equal(t, tenant, "go-sc-tenant")
}
