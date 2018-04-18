# Metrics
## 概述

Metrics用于度量服务性能指标。开发者可通过配置文件来将框架自动手机的metrics导出并让prometheus收集。

如果有业务代码定制的metrics，也可以通过API来调用，来定制自己的的metrics

## 配置

**cse.metrics.enable**
> *(optional, bool)* 是否开启metrics功能，默认为*false*

**cse.metrics.apipath**
> *(optional, string)* metrics接口，默认为*/metrics*

**cse.metrics.enableGoRuntimeMetrics**
>*(optional, bool)* 是否开启go runtime监测，默认为*false*

## API

包路径

```go
import "github.com/ServiceComb/go-chassis/metrics"
```

获取go-chassis的metrics registry，用户定制的metrics，可以通过这个registry来添加，最终也会自动导出到API的返回中

```go
func GetSystemRegistry() metrics.Registry
```

获取go-chassis使用的prometheus registry，允许用户直接对Prometheus registry进行操作

```go
func GetSystemPrometheusRegistry() *prometheus.Registry
```

创建一个特定名称的metrics registry

```go
func GetOrCreateRegistry(name string) metrics.Registry
```

使用特定metrics registry向prometheus汇报metrics数据

```go
func ReportMetricsToPrometheus(r metrics.Registry)
```

汇报metrics数据的http handler

```go
func MetricsHandleFunc(req *restful.Request, rep *restful.Response)
```

## 示例

```yaml
cse:
  metrics:
    apiPath: /metrics      # we can also give api path having prefix "/" ,like /adas/metrics
    enable: true
    enableGoRuntimeMetrics: true
```

若rest监听在127.0.0.1:8080，则作上述配置后，可通过 [http://127.0.0.1:8080/metrics](http://127.0.0.1:8080/metrics) 获取metrics数据。

