# Fault Tolerance
## 概述

go-chassis提供自动重试的容错能力，用户可配置retry及backOff策略自动启用重试功能。

## 配置

重试功能的相关配置与客户端负载均衡策略都在chassis.yaml的cse.loadbalance.配置下。当retryEnabled配置为true时，可通过配置retryOnSame和retryOnNext定制重试次数。另外可通过backOff定制重试策略，默认支持三种backOff策略。

- zero:  固定重试时间为0的重试策略，即失败后立即重试不等待。
- constant: 固定时间为backoff.minMs的重试策略，即失败后等待backoff.minMs再重试。
- jittered: 按指数增加重试时间的重试策略，初始重试时间为backoff.minMs，最大重试时间为backoff.MaxMs。

**retryEnabled**
> *(optional, bool)* 是否开启重试功能, 默认值为*false*

**retryOnSame**
> *(optional, int)* 请求失败后向同一个实例重试的次数，默认为*0*

**retryOnNext**
> *(optional, int)* 请求失败后向其他实例重试的次数，默认为*0*

**backoff.kind**
> *(optional, string)* 重试策略: [jittered或constant或zero] 默认为*zero*

**retryEnabled**
> *(optional, int)* 重试最小时间间隔 单位ms , 默认值为*0*

**retryEnabled**
> *(optional, int)* 重试最大时间间隔 单位ms, 默认值为*0*

## 示例

配置chassis.yaml负载均衡部分中的重试参数。

```yaml
cse:
  loadbalance:
    retryEnabled: true
    retryOnNext: 2
    retryOnSame: 3
    backoff:
      kind: jittered
      MinMs: 200
      MaxMs: 400
```



