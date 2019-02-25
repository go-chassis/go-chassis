# Tracing
## Introduction
Go chassis use opentacing-go to trace distributed system call
## Configuration

the config is in monitoring.yaml

**tracing.tracer**

> *(optional, string)*  what kind of opentracing plugin go chassis should use, default is *zipkin*.
But you must import zipkin plugin to enable tracing. otherwise it will report WARN log to say tracing
is not working.

**tracing.settings**

>  *(optional, map)* options like URI, batchSize, BatchInterval can be custom in here
>  go chassis tracing pkg is highly extensible, to deal with varies of different tracer settings, 
it use map to define options, so that developers can freely custom options for tracer


## Example

you must import tracing plugin pkg in main.go, below use zipkin for tracing
```go
import _ "github.com/go-chassis/go-chassis-plugins/tracing/zipkin"
```

this config means send data to zipkin, tracing-provider must to be added in handler chain

```yaml
cse:
  handler:
    chain:
      Provider:
        default: tracing-provider,bizkeeper-provider
tracing:
  tracer: zipkin
  settings:
    URI: http://127.0.0.1:9411/api/v1/spans
    batchSize: 1
```


When you have more than 2-levels service calling like A->B->C

in B client you must deliver ctx to C, so that go chassis can keep tracing,

```go
//Trace is a method
func (r *TracingHello) Trace(b *rf.Context) {
	req, err := rest.NewRequest("GET", "http://RESTServerB/sayhello/world")
	if err != nil {
		b.WriteError(500, err)
		return
	}
	defer req.Close()
    // must set b.Ctx as input for next calling
	resp, err := core.NewRestInvoker().ContextDo(b.Ctx, req)
	if err != nil {
		b.WriteError(500, err)
		return
	}
	b.Write(resp.ReadBody())
}
```

check [examples](https://github.com/go-chassis/go-chassis-examples/tree/master/monitoring)