# Metrics
## 概述

go chassis runtime will produce metrics, 
developer is also able to custom metrics.
by default go chassis use prometheus lib to produce metrics data.
user can custom metrics plugin to replace prometheus.

## 配置

**cse.metrics.enable**
> *(optional, bool)* if it is true, 
a new http API defined in "cse.metrics.apipath" will serve for client
default is *false*

**cse.metrics.apipath**
> *(optional, string)* metrics接口，默认为*/metrics*

**cse.metrics.enableGoRuntimeMetrics**
>*(optional, bool)* 是否开启go runtime监测，默认为*true*

**cse.metrics.enableCircuitMetrics**
>*(optional, bool)* report circuit breaker metrics to go-metrics, default is *true*

**cse.metrics.flushInterval**
> *(optional, string)* interval flush metrics from go-metrics to prometheus exporter, 
for example 10s, 1m

**cse.metrics.circuitMetricsConsumerNum**
> *(optional, int)* should be careful about this option, default is 3, 
there is 3 go routines consume metrics, if there is so many consumers, during high concurrency, 
it will affect service performance

## Custom Metrics
The API is in
```go
github.com/go-chassis/go-chassis/pkg/metrics/metrics.go
``` 


## 示例

```yaml
cse:
  metrics:
    apiPath: /metrics      # we can also give api path having prefix "/" ,like /adas/metrics
    enable: true
    enableGoRuntimeMetrics: true
    enableCircuitMetrics: true
```

若rest监听在127.0.0.1:8080，则作上述配置后，
可通过 [http://127.0.0.1:8080/metrics](http://127.0.0.1:8080/metrics) 获取metrics数据。

