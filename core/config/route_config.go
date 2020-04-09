package config

import "github.com/go-chassis/go-archaius"

//DefaultRouterType set the default router type
const DefaultRouterType = "cse"

// GetRouterType returns the type of router
func GetRouterType() string {
	return archaius.GetString("servicecomb.service.router.infra", DefaultRouterType)
}

// GetRouterEndpoints returns the router address
func GetRouterEndpoints() string {
	return archaius.GetString("servicecomb.service.router.address", "")
}
