# 熔断与降级
## **概述**

降级策略是当服务请求异常时，微服务所采用的异常处理策略。

降级策略有三个相关的技术概念：“隔离”、“熔断”、“容错”：

* “隔离”是一种异常检测机制，常用的检测方法是请求超时、流量过大等。一般的设置参数包括超时时间、同时并发请求个数等。
* “熔断”是一种异常反应机制，“熔断”依赖于“隔离”。熔断通常基于错误率来实现。一般的设置参数包括统计请求的个数、错误率等。
* “容错”是一种异常处理机制，“容错”依赖于“熔断”。熔断以后，会调用“容错”的方法。一般的设置参数包括调用容错方法的次数等。

把这些概念联系起来：当"隔离"措施检测到N次请求中共有M次错误的时候，"熔断"不再发送后续请求，调用"容错"处理函数。这个技术上的定义，是和Netflix Hystrix一致的，通过这个定义，非常容易理解它提供的配置项，参考：[https://github.com/Netflix/Hystrix/wiki/Configuration](https://github.com/Netflix/Hystrix/wiki/Configuration)。当前ServiceComb提供两种容错方式，分别为返回null值和抛出异常。

## **配置**

配置格式为：

cse.{namespace}.Consumer.{serviceName}.{property}: {configuration}

字段意义：

{namespace}取值为：isolation\|circuitBreaker\|fallback\|fallbackpolicy，分别表示隔离、熔断、降级、降级策略。

{serviceName}表示服务名，即某个服务提供者。

{property}表示具体配置项。

{configuration}表示具体配置内容。

为了方便描述，下表中的配置项均省略了Consumer和{serviceName}。

| 配置项 | 默认值 | 取值范围 | 是否必选 | 含义 | 注意 |
| :--- | :--- | :--- | :--- | :--- | :--- |
| cse.isolation.timeout.enabled | FALSE | - | 否 | 是否启用超时检测 |  |
| cse.isolation.timeoutInMilliseconds | 30000 | - | 否 | 超时时间阈值 |  |
| cse.isolation.maxConcurrentRequests | 1000 | - | 否 | 最大并发数阈值 |  |
| cse.circuitBreaker.enabled | TRUE | - | 否 | 是否启用熔断措施 |  |
| cse.circuitBreaker.forceOpen | FALSE | - | 否 | 不管失败次数，都进行熔断 |  |
| cse.circuitBreaker.forceClosed | FALSE | - | 否 | 任何时候都不熔断 | 当与forceOpen同时配置时，forceOpen优先。 |
| cse.circuitBreaker.sleepWindowInMilliseconds | 15000 | - | 否 | 熔断后，多长时间恢复 | 恢复后，会重新计算失败情况。注意：如果恢复后的调用立即失败，那么会立即重新进入熔断。 |
| cse.circuitBreaker.requestVolumeThreshold | 20 | - | 否 | 10s内统计错误发生次数阈值，超过阈值则触发熔断 | 由于10秒还会被划分为10个1秒的统计周期，经过1s中后才会开始计算错误率，因此从调用开始至少经过1s，才会发生熔断。 |
| cse.circuitBreaker.errorThresholdPercentage | 50 | - | 否 | 错误率阈值，达到阈值则触发熔断 |  |
| cse.fallback.enabled | TRUE | - | 否 | 是否启用出错后的故障处理措施 |  |
| cse.fallbackpolicy.policy | returnnull | returnnull \| throwexception | 否 | 出错后的处理策略 |  |

## **示例**

```yaml
---
cse:
  isolation:
    Consumer:
      timeout:
        enabled: false
      timeoutInMilliseconds: 1
      maxConcurrentRequests: 100
      Server:
        timeout:
          enabled: true
        timeoutInMilliseconds: 1000
        maxConcurrentRequests: 1000
  circuitBreaker:
    Consumer:
      enabled: false
      forceOpen: false
      forceClosed: true
      sleepWindowInMilliseconds: 10000
      requestVolumeThreshold: 20
      errorThresholdPercentage: 10
      Server:
        enabled: true
        forceOpen: false
        forceClosed: false
        sleepWindowInMilliseconds: 10000
        requestVolumeThreshold: 20
        errorThresholdPercentage: 5
  fallback:
    Consumer:
      enabled: true
      maxConcurrentRequests: 20
  fallbackpolicy:
    Consumer:
      policy: throwexception
```



