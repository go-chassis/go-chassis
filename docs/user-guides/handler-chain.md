# Handler chain
## 概述
处理链中包含一系列handler, 在一次调用中可以通过扩展调用链的方式来插入定制的逻辑处理，如何开发新的Handler，可参考Developer Guide，本章节只讨论如何进行配置，以及框架已经实现的handler有哪些

## 配置
```yaml
cse:
  protocols:
    {name}:
  handler:
    Consumer:
      {name}:{handler_names}
    Provider:
      {name}:{handler_names}
```
Consumer表示，当你要调用别的服务时，会通过的处理链  
Provider表示，当你被别人调用人，会通过的处理链  
支持在Consumer和Provider中定义多个不同的chain name  
当handler配置为空，那么框架会自动为Consumer与Provider加载默认的handlers，chain的名称为default  
当handler配置不为空，那么在Provider链中，如果protocol有相同名称的chain，对应的protocol服务将加载相同名称的chain，如果protocol没有配置相同名称的chain，那么该协议将默认加载名称为default的chain

### Consumer的默认chain为

名称	功能


ratelimiter-consumer	客户端限流

router	路由策略

loadbalance	负载均衡

tracing-consumer	客户端调用链追踪

transport	各协议客户端处理请求，如果你使用自定义处理链配置，那么结尾处必须加入这个handler

### Provider的默认chain为

名称	功能

ratelimiter-provider	服务端限流

tracing-provider	服务端调用链追踪

## API
当处理链配置为空，用户也可自定义自己的默认处理链
```go
//SetDefaultConsumerChains your custom chain map for Consumer,if there is no config, this default chain will take affect
func SetDefaultConsumerChains(c map[string]string)
//SetDefaultProviderChains set your custom chain map for Provider,if there is no config, this default chain will take affect
func SetDefaultProviderChains(c map[string]string)


```

you can check build-in handler list in [handler.go](https://github.com/go-chassis/go-chassis/blob/master/core/handler/handler.go)
the const part shows handler list
```go
const (
	//consumer chain
	Transport           = "transport"
	Loadbalance         = "loadbalance"
	BizkeeperConsumer   = "bizkeeper-consumer"
	TracingConsumer     = "tracing-consumer"
	RatelimiterConsumer = "ratelimiter-consumer"
	Router              = "router"
	FaultInject         = "fault-inject"

	//provider chain
	RatelimiterProvider = "ratelimiter-provider"
	TracingProvider     = "tracing-provider"
	BizkeeperProvider   = "bizkeeper-provider"
)
```
## 实例

### 自定义了custom-handler，放入到默认链中
```yaml
handler:
  chain:
    Consumer:
      default: custom-handler, bizkeeper-consumer, router, loadbalance, ratelimiter-consumer,transport
```

### 为不同协议定制不同的链
```yaml
protocols:
  rest:
    listenAddress: 127.0.0.1:5001
  rest-admin:
    listenAddress: 127.0.0.1:5002
  highway:
    listenAddress: 127.0.0.1:5003
  grpc:
    listenAddress: 127.0.0.1:5004
handler:
  chain:
    Provider:
      rest: custom-handler, ratelimiter-provider
      highway: tracing-provider
```
如上配置，rest协议将会加载rest chain，highway协议将会加载highway chain。由于rest-admin与grpc协议没有配置相同名称的chain，所以他们将默认加载default chain