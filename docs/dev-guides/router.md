# Router 
this guide shows how to customize your own router management.

### Introduction

A Router plugin should has ability to fetch router rule, 
it decides where the route rule comes from.

go chassis has standardized model "config.RouteRule", you must adapt to it
```go
type Router interface {
	Init(Options) error
	SetRouteRule(map[string][]*config.RouteRule)
	FetchRouteRuleByServiceName(service string) []*config.RouteRule
	ListRouteRule() map[string][]*config.RouteRule
}
```
### Usage
First, install your plugin 
```go
router.InstallRouterPlugin("istio", func() (router.Router, error) {
			//return your router implementation
		})
```

Second, specify your plugin name in router.yaml
```yaml
servicecomb:
  router:
    plugin: istio
    address: "xxx"
``` 

go chassis will use your router implementation as router rule configuration source, 
to know how to manage request traffics, 
refer to [Router](https://go-chassis.readthedocs.io/en/latest/user-guides/router.html) 