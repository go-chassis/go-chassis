# RouterConfigure

### Introduction

A Router plugin should has ability to fetch router rule, 
it decides where the route rule comes from.
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
			//your implementation
		})
```

Second, specify your plugin name in router.yaml
```yaml
servicecomb:
  service:
    router:
      plugin: istio
      address: "xxx"
``` 

go chassis will use your router implementation as router rule configuration source, 
to know how to manage request traffics, 
refer to [Router](https://go-chassis.readthedocs.io/en/latest/user-guides/router.html) 