# Tracing
## 概述

调用跟踪模块主要实现与“服务监控”的对接，按照服务监控的要求，产生span和SLA数据，按照配置文件的定义将产生的调用链数据输出到zipkin或文件中，用于分析调用链的调用过程和状态。

## 配置

使用调用链追踪功能，必须先在handler chain中添加对应handler：tracing-provider或tracing-consumer,默认存在。

**tracing.collectorType**

> *(requied, string)*  定义调用数据向什么服务发送，支持 *zipkin*，*namedPipe*

**tracing.collectorTarget**

>  *(requied, string)* 服务的URI，比如文件路径，http地址。
>  collectorType为http时，collectorTarget为zipkin地址否则，namedPipe为文件路径

## 示例

追踪数据发送至zipkin:

```yaml
cse:
  handler:
    chain:
      Provider:
        default: tracing-provider,bizkeeper-provider
tracing:
  collectorType: zipkin
  collectorTarget: http://localhost:9411/api/v1/spans
```

追踪数据写入linux named pipe:

```yaml
tracing:
  collectorType: namedPipe
  collectorTarget: /home/chassis.trace
```



