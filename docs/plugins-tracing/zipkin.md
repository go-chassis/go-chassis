# Zipkin

Zipkin tracer is a plugin of go chassis, it reports tracing data to zipkin server

## Configurations
you must import tracing plugin pkg in main.go
```go
import _ "github.com/go-chassis/go-chassis-plugins/tracing/zipkin"
```

**tracing.settings.URI**
>  *(optional, string)* zipkin api url

**tracing.settings.batchSize**
>  *(optional, string)* after how many data collected, a tracer will report to server, default is 10000

**tracing.settings.batchInterval**
>  *(optional, string)* after how long the tracer running, a tracer will report to server, default is 10s

**tracing.settings.collector**
>  *(optional, string)* support http, namedPipe, default is http


## Example
```yaml
tracing:
  tracer: zipkin
  settings:
    URI: http://127.0.0.1:9411/api/v1/spans
    batchSize: 10000
    batchInterval: 10s
    collector: http
    
```
