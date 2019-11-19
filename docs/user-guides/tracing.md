# Tracing
## Introduction
Go chassis use opentacing-go or skywalking to trace distributed system call
# 1 Opentacing-go
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

# 2 Skywalking
## Configurations
**In conf/monitoring.yaml**

**servicecomb.apm.tracing.tracer**
>  *(optional, string)* tracer'name, only skywalking now

**servicecomb.apm.tracing.settings**
>  *(optional, map)* options including :
>  URI server address of skywalking 
>  servertype: service type, match componentid in skywalking ex:  5001:servicecomb-mesher 5002:servicecomb-service-center 28:servicecomb-java-cahssis 
>  enable: if open

## Example
```yaml
servicecomb:
  apm:                                #application performance monitor
    tracing:
      tracer: skywalking
      settings:
        enable: true                  #enable tracing ability
        URI: 127.0.0.1:11800          #url of skywalking 
        servertype: 5001              #server type
```
## Step：
**handler name is defined in github.com/go-chassis/go-chassis/core/handler**
- [1] Adding handler name 'handler.SkyWalkingConsumer' in consumerChain.
- [2] Adding handler name 'handler.SkyWalkingProvider' in providerChain.
## Example
```golang
consumerChain := strings.Join([]string{
		chassisHandler.Router,
		chassisHandler.RatelimiterConsumer,
		chassisHandler.BizkeeperConsumer,
		chassisHandler.Loadbalance,
		chassisHandler.Transport,
		chassisHandler.SkyWalkingConsumer,
	}, ",")
	providerChain := strings.Join([]string{
		chassisHandler.RatelimiterProvider,
		chassisHandler.Transport,
		chassisHandler.SkyWalkingProvider,
	}, ",")
```
