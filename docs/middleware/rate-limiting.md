# Rate limiting v1
## 概述

用户可以通过配置限流策略限制provider端的请求频率，使每秒请求数限制在最大请求量的大小。配置可限制接收处理请求的频率

## 配置

限流配置在rate_limiting.yaml中，同时需要在chassis.yaml的handler chain中添加handler。

**cse.flowcontrol.Provider.qps.enabled**
> *(optional, bool)* 是否开启限流，默认true

**cse.flowcontrol.Provider.qps.global.limit**
> *(optional, int)* 每秒允许的请求数，默认2147483647max int）

引入middleware
```go
import _ github.com/go-chassis/go-chassis/v2/middleware/ratelimiter
```
#### Provider示例

provider端需要在chassis.yaml添加ratelimiter-provider。同时在rate\_limiting.yaml中配置具体的请求数。

```yaml
servicecomb:
  handler:
    chain:
      Provider:
        default: ratelimiter-provider
```

```yaml
cse:
  flowcontrol:
    Provider:
      qps:
        enabled: true  # enable rate limiting or not
        global:
          limit: 100 
```
