# Rate limiting
## 概述

用户可以通过配置限流策略限制provider端或consumer端的请求频率，使每秒请求数限制在最大请求量的大小。其中provider端的配置可限制接收处理请求的频率，consumer端的配置可限制发往指定微服务的请求的频率。

## 配置

限流配置在rate\_limiting.yaml中，同时需要在chassis.yaml的handler chain中添加handler。其中qps.limit.\[service\] 是指限制从service 发来的请求的处理频率，若该项未配置则global.limit生效。Consumer端不支持global全局配置，其他配置项与Provider端一致。

**flowcontrol.qps.enabled**
> *(optional, bool)* 是否开启限流，默认true

**flowcontrol.qps.global.limit**
> *(optional, int)* 每秒允许的请求数，默认2147483647max int）

**flowcontrol.qps.limit.{service}**
> *(optional, string)* 针对某微服务每秒允许的请求数 ，默认2147483647max int）


#### Provider示例

provider端需要在chassis.yaml添加ratelimiter-provider。同时在rate\_limiting.yaml中配置具体的请求数。

```yaml
cse:
  handler:
    chain:
      Provider:
        default: ratelimiter-provider
```

```yaml
cse:
  flowcontrol
    Provider:
      qps:
        enabled: true  # enable rate limiting or not
        global:
          limit: 100   # default limit of provider
        limit:
          Server: 100  # rate limit for request from a provider
```

#### Consumer示例

在consumer端需要添加ratelimiter-consumer这个handler。同时在rate\_limiting.yaml中配置具体的请求数。

```yaml
cse:
  handler:
    chain:
      Consumer:
        default: ratelimiter-consumer
```

```yaml
cse:
  flowcontrol:
    Consumer:
      qps:
        enabled: true  # enable rate limiting or not
        limit:
          Server: 100  # rate limit for request to a provider
```

## API

qpslimiter提供获取流控实例的接口GetQpsTrafficLimiter和相关的处理接口。其中ProcessQpsTokenReq根据目标qpsRate在handler chain当中sleep相应时间实现限流，UpdateRateLimit提供更新qpsRate限制的接口，DeleteRateLimiter提供了删除流控实例的接口。

##### 对请求流控

```go
qpslimiter.GetQpsTrafficLimiter().ProcessQpsTokenReq(key string, qpsRate int)
```

##### 更新流控限制

```go
qpslimiter.GetQpsTrafficLimiter().UpdateRateLimit(key string, value interface{})
```

##### 删除流控实例

```go
qpslimiter.GetQpsTrafficLimiter().DeleteRateLimiter(key string)
```



