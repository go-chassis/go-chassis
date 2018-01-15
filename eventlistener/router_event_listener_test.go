package eventlistener_test

import (
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/route"
	"github.com/ServiceComb/go-chassis/eventlistener"
	"github.com/stretchr/testify/assert"
	"testing"
)

var darkLaunch = []byte(`
{
  "policyType": "RATE",
  "ruleItems": [
    {
      "groupName": "s1",
      "groupCondition": "version=0.3",
      "policyCondition": "30"
    },
    {
      "groupName": "s2",
      "groupCondition": "version=0.4",
      "policyCondition": "70"
    }
  ]
}`)

var darkLaunch1 = []byte(`
{
  "policyType": "RULE",
  "ruleItems": [
    {
      "groupName": "s1",
      "groupCondition": "version=0.3",
      "policyCondition": "test!=30"
    },
    {
      "groupName": "s2",
      "groupCondition": "version=0.4",
      "policyCondition": "t>3"
    }
  ]
}`)

func TestDarkLaunchEventListener_Event(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	archaius.Init()
	eventlistener.Init()
	eventListen := &eventlistener.DarkLaunchEventListener{}
	e := &core.Event{EventType: "CREATE", Key: eventlistener.DarkLaunchPrefix + "service", Value: string(darkLaunch)}
	eventListen.Event(e)
	assert.Equal(t, router.GetRouteRule()["service"][0].Routes[0].Weight, 30)

	e = &core.Event{EventType: "UPDATE", Key: eventlistener.DarkLaunchPrefix + "service1", Value: string(darkLaunch1)}
	eventListen.Event(e)
	t.Log(router.GetRouteRule())
	assert.Equal(t, router.GetRouteRule()["service1"][1].Match.HTTPHeaders["t"]["greater"], "3")
	assert.Equal(t, router.GetRouteRule()["service1"][0].Match.HTTPHeaders["test"]["noEqu"], "30")

	e2 := &core.Event{EventType: "DELETE", Key: eventlistener.DarkLaunchPrefix + "service", Value: ""}
	eventListen.Event(e2)
	assert.Equal(t, len(router.GetRouteRule()), 1)
}
