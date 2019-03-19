package cse

import (
	"testing"

	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-mesh/openlogging"
	"github.com/stretchr/testify/assert"
)

func TestDarkLaunchRule2RouteRule(t *testing.T) {
	openlogging.GetLogger().Info("check translate,type DarkLaunchTypeRule")

	routeRules := DarkLaunchRule2RouteRule(getRule(DarkLaunchTypeRule, "version=1.0",
		"foo=bar", []string{"1.0"}))
	// check  routeRules
	assert.NotNil(t, routeRules)
	assert.NotZero(t, len(routeRules))

	// check Routes
	for _, routeRule := range routeRules {
		assert.NotNil(t, routeRule)
		assert.NotZero(t, len(routeRule.Routes))
		for _, r := range routeRule.Routes {
			assert.NotNil(t, r)
			assert.Equal(t, r.Tags, map[string]string{"version": "1.0"})
		}
		// check match headers
		v, ok := routeRule.Match.Headers["Foo"]
		assert.True(t, ok)
		v1, ok := v["exact"]
		assert.True(t, ok)
		assert.Equal(t, v1, "bar")
		v1, ok = v["caseInsensitive"]
		assert.True(t, ok)
		assert.Equal(t, v1, "false")
	}

	routeRules = DarkLaunchRule2RouteRule(getRule(DarkLaunchTypeRule, "version=",
		"Foo", []string{"1.0"}))
	// check  routeRules
	assert.NotNil(t, routeRules)
	assert.NotZero(t, len(routeRules))
	// check Routes
	for _, routeRule := range routeRules {
		assert.NotNil(t, routeRule)
		assert.NotZero(t, len(routeRule.Routes))
		for _, r := range routeRule.Routes {
			assert.NotNil(t, r)
			assert.Equal(t, r.Tags, map[string]string{"version": ""})
		}
		// check match headers
		v, ok := routeRule.Match.Headers["foo"]
		assert.False(t, ok)
		assert.Nil(t, v)
		v, ok = routeRule.Match.Headers["Foo"]
		assert.False(t, ok)
		assert.Nil(t, v)
	}

	routeRules = DarkLaunchRule2RouteRule(getRule(DarkLaunchTypeRule, "version",
		"Foo=", []string{"1.0"}))
	// check  routeRules
	assert.NotNil(t, routeRules)
	assert.NotZero(t, len(routeRules))
	for _, routeRule := range routeRules {
		assert.NotNil(t, routeRule)
		assert.NotZero(t, len(routeRule.Routes))
		for _, r := range routeRule.Routes {
			assert.NotNil(t, r)
			assert.Equal(t, r.Tags, map[string]string{"version": "version"})
		}
		// check match headers
		v, ok := routeRule.Match.Headers["Foo"]
		assert.True(t, ok)
		v1, ok := v["exact"]
		assert.True(t, ok)
		assert.Equal(t, v1, "")
		v1, ok = v["caseInsensitive"]
		assert.True(t, ok)
		assert.Equal(t, v1, "false")
	}

	//	 DarkLaunchTypeRate

	openlogging.GetLogger().Info("check translate,type DarkLaunchTypeRate")
	routeRules = DarkLaunchRule2RouteRule(getRule(DarkLaunchTypeRate,
		"version=1.0", "50", []string{"1.0"}))
	for _, routeRule := range routeRules {
		assert.NotNil(t, routeRule)
		assert.NotZero(t, len(routeRule.Routes))

		for _, r := range routeRule.Routes {
			assert.NotNil(t, r)
			assert.Equal(t, r.Tags, map[string]string{"version": "1.0"})
			assert.Equal(t, r.Weight, 50)
		}
	}

	routeRules = DarkLaunchRule2RouteRule(getRule(DarkLaunchTypeRate,
		"version=", "", []string{"1.0"}))
	for _, routeRule := range routeRules {
		assert.NotNil(t, routeRule)
		assert.NotZero(t, len(routeRule.Routes))

		for _, r := range routeRule.Routes {
			assert.NotNil(t, r)
			assert.Equal(t, r.Tags, map[string]string{"version": ""})
			assert.Equal(t, r.Weight, 0)
		}
	}
}
func getRule(darkType, groupCondition, policyCondition string, Versions []string) *model.DarkLaunchRule {
	return &model.DarkLaunchRule{
		Type: darkType,
		Items: []*model.RuleItem{
			{
				GroupName:       "test",
				GroupCondition:  groupCondition,
				PolicyCondition: policyCondition,
				Versions:        Versions,
			},
		},
	}
}
