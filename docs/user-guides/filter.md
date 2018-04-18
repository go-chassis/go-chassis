# 负载均衡过滤器
## 概述

负载均衡过滤器实现在负载均衡模块中，它允许开发者定制自己的过滤器，使得经过负载均衡策略选择实例前，预先对实例进行筛选。在一次请求调用过程中可以使用多个过滤策略，对从本地cache或服务中心获取的实例组进行过滤，将经过精简的实例组交给负载均衡策略做后续处理。

## 配置

目前可配的filter只有根据Available Zone Filter。可根据微服务实例的region以及AZ信息进行过滤，优先寻找同Region与AZ的实例。

```
cse:
  loadbalance:
    serverListFilters: zoneaware
```

需要配置实例的Datacenter信息

```
region:
  name: us-east
  availableZone: us-east-1
```

## API

Go-chassis支持多种实现Filter接口的过滤器。FilterEndpoint支持通过实例访问地址过滤，FilterMD支持通过元数据过滤，FilterProtocol支持通过协议过滤，FilterAvailableZoneAffinity支持根据Zone过滤。

```go
type Filter func([]*registry.MicroServiceInstance) []*registry.MicroServiceInstance
```

## 示例

客户端实例过滤器Filter的使用支持用户通过API调用传入，并且可以一次传入多个Filter，对实例组进行层层条件筛选。客户端调用传入的方式是调用options的WithFilter方法。

```go
invoker.Invoke(ctx, "Server", "HelloServer", "SayHello",
    &helloworld.HelloRequest{Name: "Peter"},
    reply,
    core.WithProtocol("highway"),
    core.WithStrategy(loadbalance.StrategyRoundRobin),
    core.WithFilters(
      loadbalance.FilterEndpoint("highway://127.0.0.1:8080"),
      loadbalance.FilterProtocol("highway"),
    ),
)
```



