package router

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
)

// Init initialize router config
func Init() error {
	// init dests and templates
	routerConfigFromFile := config.RouterDefinition
	//TODO from config
	BuildRouter("cse")
	if routerConfigFromFile != nil {
		if routerConfigFromFile.Destinations != nil {
			DefaultRouter.SetRouteRule(routerConfigFromFile.Destinations)
		}
		if routerConfigFromFile.SourceTemplates != nil {
			Templates = routerConfigFromFile.SourceTemplates
		}
	}

	lager.Logger.Info("Router init success")
	return nil
}

// ValidateRule validate the route rules of each service
func ValidateRule(rules map[string][]*model.RouteRule) bool {
	for name, rule := range rules {

		for _, route := range rule {
			allWeight := 0
			for _, routeTag := range route.Routes {
				allWeight += routeTag.Weight
			}

			if allWeight > 100 {
				lager.Logger.Warnf("route rule for [%s] is not valid: ruleTag weight is over 100%", name)
				return false
			}
		}

	}
	return true
}
