package cse

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/lager"

	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/router"
	"github.com/stretchr/testify/assert"
)

//ExternalConfigurationSource is a struct
type ExternalConfigurationSource struct {
	Configurations map[string]interface{}
	callback       core.DynamicConfigCallback
	sync.RWMutex
	CallbackCheck chan bool
	ChanStatus    bool
}

var externalConfigSource *ExternalConfigurationSource

//NewExternalConfigurationSource initializes all necessary components for external configuration
func NewExternalConfigurationSource() *ExternalConfigurationSource {
	if externalConfigSource == nil {
		externalConfigSource = new(ExternalConfigurationSource)
		externalConfigSource.Configurations = make(map[string]interface{})
		externalConfigSource.CallbackCheck = make(chan bool)
	}

	return externalConfigSource
}

//GetConfigurations gets all external configurations
func (confSrc *ExternalConfigurationSource) GetConfigurations() (map[string]interface{}, error) {
	configMap := make(map[string]interface{})

	confSrc.Lock()
	defer confSrc.Unlock()
	for key, value := range confSrc.Configurations {
		configMap[key] = value
	}

	return configMap, nil
}

//GetConfigurationByKey gets required external configuration for a particular key
func (confSrc *ExternalConfigurationSource) GetConfigurationByKey(key string) (interface{}, error) {
	confSrc.Lock()
	defer confSrc.Unlock()
	value, ok := confSrc.Configurations[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}

	return value, nil
}

//GetPriority returns priority of the external configuration
func (*ExternalConfigurationSource) GetPriority() int {
	return 100
}

//GetSourceName returns name of external configuration
func (*ExternalConfigurationSource) GetSourceName() string {
	return "ExternalConfigurationSource"
}

//DynamicConfigHandler dynamically handles a external configuration
func (confSrc *ExternalConfigurationSource) DynamicConfigHandler(callback core.DynamicConfigCallback) error {
	confSrc.callback = callback
	confSrc.CallbackCheck <- true
	return nil
}

//AddKeyValue creates new configuration for corresponding key and value pair
func (confSrc *ExternalConfigurationSource) AddKeyValue(key string, value interface{}) error {
	if !confSrc.ChanStatus {
		<-confSrc.CallbackCheck
		confSrc.ChanStatus = true
	}

	event := new(core.Event)
	event.EventSource = confSrc.GetSourceName()
	event.Key = key
	event.Value = value

	confSrc.Lock()
	if _, ok := confSrc.Configurations[key]; !ok {
		event.EventType = core.Create
	} else {
		event.EventType = core.Update
	}

	confSrc.Configurations[key] = value
	confSrc.Unlock()

	if confSrc.callback != nil {
		confSrc.callback.OnEvent(event)
	}

	return nil
}

//DeleteKey deletes a key from source
func (confSrc *ExternalConfigurationSource) DeleteKey(key string) error {
	if !confSrc.ChanStatus {
		<-confSrc.CallbackCheck
		confSrc.ChanStatus = true
	}

	event := new(core.Event)
	event.EventSource = confSrc.GetSourceName()
	event.EventType = core.Delete
	event.Key = key

	confSrc.Lock()
	delete(confSrc.Configurations, key)
	confSrc.Unlock()

	if confSrc.callback != nil {
		confSrc.callback.OnEvent(event)
	}

	return nil
}

//Cleanup cleans a particular external configuration up
func (confSrc *ExternalConfigurationSource) Cleanup() error {
	confSrc.Configurations = nil

	return nil
}

//GetConfigurationByKeyAndDimensionInfo gets a required external configuration for particular key and dimension info pair
func (*ExternalConfigurationSource) GetConfigurationByKeyAndDimensionInfo(key, di string) (interface{}, error) {
	return nil, nil
}

//AddDimensionInfo adds dimension info for a external configuration
func (*ExternalConfigurationSource) AddDimensionInfo(dimensionInfo string) (map[string]string, error) {
	return nil, nil
}

//GetConfigurationsByDI gets required external configuration for a particular dimension info
func (*ExternalConfigurationSource) GetConfigurationsByDI(dimensionInfo string) (map[string]interface{}, error) {
	return nil, nil
}

// the service key in router config file
const (
	svcNone               = "svcNone"
	svcRoute              = "svcRoute"
	svcDarkLaunch         = "svcDarkLaunch"
	svcRouteAndDarkLaunch = "svcRouteAndDarkLaunch"
)

//to diff from the router config between different service and operation
const (
	darkLaunchRuleNumSvcDarkLaunch         = 2
	darkLaunchRuleNumSvcRouteAndDarkLaunch = 3
	darkLaunchRuleNumAfterAdd              = 4
	darkLaunchRuleNumAfterUpdate           = 5
)

func genSvcRouteRule() []*model.RouteRule {
	r := []*model.RouteRule{
		{
			Precedence: 0,
			Routes: []*model.RouteTag{
				{
					Tags: map[string]string{
						common.BuildinTagVersion: "0.0.1",
						common.BuildinTagApp:     svcRoute,
					},
					Weight: 20,
				},
			},
		},
	}
	return r
}

func genSvcDarkLaunchRule() string {
	return `{"policyType":"RATE","ruleItems":[{"groupName":"s0"},{"groupName":"s1"}]}`
}

func genSvcRouteAndDarkLaunchRule() ([]*model.RouteRule, string) {
	r := []*model.RouteRule{
		{
			Precedence: 1,
			Routes: []*model.RouteTag{
				{
					Tags: map[string]string{
						common.BuildinTagVersion: "0.0.1",
						common.BuildinTagApp:     svcRouteAndDarkLaunch,
					},
					Weight: 20,
				},
			},
		},
	}
	return r, `{"policyType":"RATE","ruleItems":[{"groupName":"s0"},{"groupName":"s1"},{"groupName":"s2"}]}`
}

func genDarkLaunchRuleAfterAdd() string {
	return `{"policyType":"RATE","ruleItems":[{"groupName":"s0"},{"groupName":"s1"},{"groupName":"s2"},{"groupName":"s3"}]}`
}

func genDarkLaunchRuleAfterUpdate() string {
	return `{"policyType":"RATE","ruleItems":[{"groupName":"s0"},{"groupName":"s1"},{"groupName":"s2"},{"groupName":"s3"},{"groupName":"s4"}]}`
}

func addDarkLaunchRule(s string) {
	key := DarkLaunchPrefix + s
	NewExternalConfigurationSource().AddKeyValue(key, genDarkLaunchRuleAfterAdd())
	e := &core.Event{
		EventSource: RouteDarkLaunchGovernSourceName,
		EventType:   core.Create,
		Key:         s,
		Value:       genDarkLaunchRuleAfterAdd(),
	}
	NewRouteDarkLaunchGovernSource().Callback(e)
}
func updateDarkLaunchRule(s string) {
	key := DarkLaunchPrefix + s
	NewExternalConfigurationSource().AddKeyValue(key, genDarkLaunchRuleAfterUpdate())
	e := &core.Event{
		EventSource: RouteDarkLaunchGovernSourceName,
		EventType:   core.Update,
		Key:         s,
		Value:       genDarkLaunchRuleAfterUpdate(),
	}
	NewRouteDarkLaunchGovernSource().Callback(e)
}
func deleteDarkLaunchRule(s string) {
	key := DarkLaunchPrefix + s
	NewExternalConfigurationSource().DeleteKey(key)
	e := &core.Event{
		EventSource: RouteDarkLaunchGovernSourceName,
		EventType:   core.Delete,
		Key:         s,
	}
	NewRouteDarkLaunchGovernSource().Callback(e)
}

func preInit(t *testing.T) {
	lager.Initialize("", "DEBUG", "", "size", true, 1, 10, 7)
	c, err := archaius.NewConfig(make([]string, 0), make([]string, 0))
	if err != nil {
		t.Error(err)
	}
	archaius.DefaultConf = c
	c.ConfigFactory.AddSource(NewExternalConfigurationSource())
}

func TestInitRouterManager(t *testing.T) {
	preInit(t)
	SetRouteRuleByKey(svcRoute, genSvcRouteRule())
	r, s := genSvcRouteAndDarkLaunchRule()
	SetRouteRuleByKey(svcRouteAndDarkLaunch, r)
	NewExternalConfigurationSource().AddKeyValue(DarkLaunchPrefix+svcRouteAndDarkLaunch, s)
	NewExternalConfigurationSource().AddKeyValue(DarkLaunchPrefix+svcDarkLaunch, genSvcDarkLaunchRule())
	err := router.Init()
	if err != nil {
		t.Error(err)
	}
	testRouteManager(t, svcNone)
	testRouteManager(t, svcRoute)
	testRouteManager(t, svcDarkLaunch)
	testRouteManager(t, svcRouteAndDarkLaunch)
}

// can test svcNone/svcRouter
func testRouteManager(t *testing.T, svc string) {
	t.Logf("====Route manager test for [%s]", svc)
	r := GetRouteRuleByKey(svc)
	if svc == svcNone {
		t.Log("After init, route should be nil")
		assert.Nil(t, r)
	} else {
		assert.NotNil(t, r)
		if r == nil {
			t.FailNow()
		}
		switch svc {
		case svcRoute:
			t.Log("After init, route should from route config")
			assert.Equal(t, r[0].Routes[0].Tags[common.BuildinTagApp], svc)
		case svcDarkLaunch:
			t.Log("After init, route should from darklaunch config")
			assert.Equal(t, darkLaunchRuleNumSvcDarkLaunch, len(r[0].Routes))
		case svcRouteAndDarkLaunch:
			t.Log("After init, route should from darklaunch config")
			assert.Equal(t, darkLaunchRuleNumSvcRouteAndDarkLaunch, len(r[0].Routes))
		}
	}

	t.Log("After add dark launch governance, route should be updated")
	addDarkLaunchRule(svc)
	time.Sleep(100 * time.Millisecond)
	r = GetRouteRuleByKey(svc)
	assert.NotNil(t, r)
	if r == nil {
		t.FailNow()
	}
	assert.Equal(t, darkLaunchRuleNumAfterAdd, len(r[0].Routes))

	t.Log("After update dark launch governance, route should be updated")
	updateDarkLaunchRule(svc)
	time.Sleep(100 * time.Millisecond)
	r = GetRouteRuleByKey(svc)
	assert.NotNil(t, r)
	assert.Equal(t, darkLaunchRuleNumAfterUpdate, len(r[0].Routes))

	deleteDarkLaunchRule(svc)
	time.Sleep(100 * time.Millisecond)
	r = GetRouteRuleByKey(svc)
	if svc == svcNone || svc == svcDarkLaunch {
		t.Log("After delete dark launch governance, route should be nil")
		assert.Nil(t, r)
	} else {
		assert.NotNil(t, r)
		if r == nil {
			t.FailNow()
		}
		t.Log("After delete dark launch governance, route should from route config")
		assert.Equal(t, r[0].Routes[0].Tags[common.BuildinTagApp], svc)
	}
}

func TestAddRouteRuleSource(t *testing.T) {
	t.Log("Add nil source should get err")
	err := AddRouteRuleSource(nil)
	assert.Error(t, err)

	t.Log("Before init, add source should get err")
	routeRuleMgr = nil
	err = AddRouteRuleSource(NewRouteDarkLaunchGovernSource())
	assert.Error(t, err)
}
