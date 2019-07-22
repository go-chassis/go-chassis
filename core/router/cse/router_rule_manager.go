package cse

import (
	"errors"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-mesh/openlogging"
	"strings"
)

// constant for route rule keys
const (
	DarkLaunchKey      = "^cse\\.darklaunch\\.policy\\."
	DarkLaunchPrefix   = "cse.darklaunch.policy."
	DarkLaunchTypeRule = "RULE"
	DarkLaunchTypeRate = "RATE"
)

//GetRouterRuleFromArchaius get router config from archaius, including memory,local file and config center
func GetRouterRuleFromArchaius() (map[string][]*model.RouteRule, error) {
	destinations := make(map[string][]*model.RouteRule, 0)
	//set config from file first
	for k, v := range config.RouterDefinition.Destinations {
		destinations[k] = v
	}
	//then get config from archaius and simply overwrite rule from file
	configMap := archaius.GetConfigs()
	//filter out key:value pairs which are not route rules
	for k := range configMap {
		if !strings.HasPrefix(k, DarkLaunchPrefix) {
			delete(configMap, k)
		}
	}

	for k, v := range configMap {
		// todo bug fix
		value, ok := v.(string)
		if !ok {
			return nil, errors.New("route rule is not a json string format please check the configuration in config center")
		}

		service := strings.Replace(k, DarkLaunchPrefix, "", 1)
		r, err := ConvertJSON2RouteRule(value)
		if err != nil {
			openlogging.Error("convert failed: " + err.Error())
		}
		destinations[service] = r
	}
	return destinations, nil
}
