# 负载均衡策略
## 概述

用户可以通过配置选择不同的负载均衡策略，当前支持轮询、随机、基于响应时间的权值、会话保持等多种负载均衡策略。

负载均衡功能作用于客户端，且依赖注册中心。

## 配置

负载均衡的配置项为cse.loadbalance.[MicroServiceName].[PropertyName]，其中若省略MicroServiceName，则为全局配置；若指定MicroServiceName，则为针对特定微服务的配置。优先级：针对特定微服务的配置 > 全局配置。

为便于描述，以下配置项说明仅针对PropertyName字段

**strategy.name**
>*(optional, bool)* RoundRobin | 策略，可选值：*RoundRobin*,*Random*,*SessionStickiness*,*WeightedResponse*。
>SessionStickiness目前只支持Rest调用。

**注意：**

1. **使用SessionStickiness**策略，需要业务代码存储cookie，并在http请求中带入Cookie。使用go-chassis进行调用时，http头中将返回如下信息：Set-Cookie: SERVICECOMBLB=0406060d-0009-4e06-4803-080008060f0d，若用户使用SessionStickiness策略，需要将将该头部信息保存，并在发送后续请求时带上如下http头：Cookie: SERVICECOMBLB=0406060d-0009-4e06-4803-080008060f0d**
2. **使用 WeightedResponse策略，启用后30s 策略会计算好数据并生效，80%左右的请求会被发送到延迟最低的实例里**

## API

除了通过配置文件传入负载均衡策略，还支持用户客户端调用传入WithStrategy的方式。

```go
invoker.Invoke(ctx, "Server", "HelloServer", "SayHello",
    &helloworld.HelloRequest{Name: "Peter"},
    reply,
    core.WithContentType("application/json"),
    core.WithProtocol("highway"),
    core.WithStrategy(loadbalance.StrategyRoundRobin),
)
```

## 示例

配置chassis.yaml的负载均衡部分，以及添加处理链。

```yaml
cse:
  loadbalance:                 # 全局负载均衡配置
    strategy:
      name: RoundRobin
    microserviceA:              # 微服务级别的负载均衡配置
      strategy:
        name: SessionStickiness
```



