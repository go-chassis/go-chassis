package eventlistener_test

import (
	"testing"
	"time"

	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-archaius/sources/memory-source"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/router"
	"github.com/ServiceComb/go-chassis/core/router/cse"
	"github.com/ServiceComb/go-chassis/eventlistener"

	"github.com/stretchr/testify/assert"
)

const (
	svcDarkLaunch       = "svcDarkLaunch"
	svcDarkLaunchConfig = `{"policyType":"RATE","ruleItems":[{"groupName":"s0"},{"groupName":"s1"}]}`
)

func TestDarkLaunchEventListenerEvent(t *testing.T) {
	lager.Initialize("", "DEBUG", "", "size", true, 1, 10, 7)
	c, err := archaius.NewConfig(make([]string, 0), make([]string, 0))
	if err != nil {
		t.Error(err)
	}
	archaius.DefaultConf = c
	c.ConfigFactory.AddSource(memoryconfigsource.NewMemoryConfigurationSource())

	err = router.Init()
	assert.NoError(t, err)

	e := &core.Event{
		EventSource: cse.RouteDarkLaunchGovernSourceName,
		EventType:   core.Create,
		Key:         svcDarkLaunch,
		Value:       svcDarkLaunchConfig,
	}

	t.Log("Before event, there should be no router config")
	assert.Nil(t, router.DefaultRouter.FetchRouteRuleByServiceName(svcDarkLaunch))

	t.Log("After event, there should exists router config")
	archaius.AddKeyValue(eventlistener.DarkLaunchPrefix+svcDarkLaunch, svcDarkLaunchConfig)
	l := &eventlistener.DarkLaunchEventListener{}
	l.Event(e)
	time.Sleep(100 * time.Millisecond)
	r := router.DefaultRouter.FetchRouteRuleByServiceName(svcDarkLaunch)
	assert.NotNil(t, r)
	if r == nil {
		t.FailNow()
	}
	assert.Equal(t, 2, len(r[0].Routes))
}
