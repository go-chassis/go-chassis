# Handler chain
## 概述
处理链中包含一系列handler, 在一次调用中可以通过扩展调用链的方式来插入定制的逻辑处理，如何开发新的Handler，可参考Developer Guide，本章节只讨论如何进行配置，以及框架已经实现的handler有哪些

## 配置
```yaml
cse:
  handler:
    Consumer:
      {name}:{handler_names}
    Provider:
      {name}:{handler_names}
```
Consumer表示，当你要调用别的服务时，会通过的处理链
Provider表示，当你被别人调用人，会通过的处理链
支持在Consumer和Provider中定义多个不同的chain name
如果handler配置为空那么框架会自动为Consumer与Provider加载默认的handlers，chain的名称为default

### Consumer的默认chain为

名称	功能


ratelimiter-consumer	客户端限流

bizkeeper-consumer	熔断降级

router	路由策略

loadbalance	负载均衡

tracing-consumer	客户端调用链追踪

transport	各协议客户端处理请求，如果你使用自定义处理链配置，那么结尾处必须加入这个handler

### Provider的默认chain为

名称	功能

ratelimiter-provider	服务端限流

tracing-provider	服务端调用链追踪

bizkeeper-provider	服务端熔断

## API
当处理链配置为空，用户也可自定义自己的默认处理链
```go
//SetDefaultConsumerChains your custom chain map for Consumer,if there is no config, this default chain will take affect
func SetDefaultConsumerChains(c map[string]string)
//SetDefaultProviderChains set your custom chain map for Provider,if there is no config, this default chain will take affect
func SetDefaultProviderChains(c map[string]string)
```
## 实例
```yaml
handler:
  chain:
    Consumer:
      default: bizkeeper-consumer, router, loadbalance, ratelimiter-consumer,transport
      custom: some-handler
```