package eventlistener

import (
	"encoding/json"
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/route"
	"strings"
)

// constants for dark launch key and prefix
const (
	//DarkLaunchKey & DarkLaunchPrefix is a variable of type string
	DarkLaunchKey    = "^cse\\.darklaunch\\.policy\\."
	DarkLaunchPrefix = "cse.darklaunch.policy."
)

//DarkLaunchEventListener is a struct
type DarkLaunchEventListener struct{}

//Event is method used for dark launch event listening
func (e *DarkLaunchEventListener) Event(event *core.Event) {
	rules := router.GetRouteRule()
	if rules == nil {
		rules = map[string][]*config.RouteRule{}
	}
	rule := &config.DarkLaunchRule{}
	serviceName := strings.Replace(event.Key, DarkLaunchPrefix, "", 1)

	switch event.EventType {
	case common.Update:
		if err := json.Unmarshal([]byte(event.Value.(string)), rule); err != nil {
			lager.Logger.Error("can not update route rule", err)
		}
		lager.Logger.Info("Route rule '" + serviceName + "' is updated to " + event.Value.(string))
		rules[serviceName] = config.TranslateRules(rule)
	case common.Create:
		if err := json.Unmarshal([]byte(event.Value.(string)), rule); err != nil {
			lager.Logger.Error("can not create route rule", err)
		}
		lager.Logger.Info("Route rule '" + serviceName + "' is created. Value=" + event.Value.(string))
		rules[serviceName] = config.TranslateRules(rule)
	case common.Delete:
		// delete route rule whose destination is the tail of event.key
		delete(rules, serviceName)
		lager.Logger.Info("Route rule '" + serviceName + "' is removed!")
	}
	router.SetRouteRule(rules)
}
