Istio
=====

go-chassis can also leverage Istio just like any simple application. 
To use istio, the simple 2 steps are needed:

### How to

1.edit chassis.yaml.

**registry.disabled** set it to true. and only reserve transport handler

```yaml
servicecomb:
  registry:
    disabled: true
  handler:
    chain:
      Consumer:
        default: transport
```

2.call remote service with Option "WithoutSD" and add port number
```go
req, _ := rest.NewRequest("GET", "http://orderService:8080/hello", nil)
invoker.ContextDo(context.TODO(), req, core.WithoutSD())
```
