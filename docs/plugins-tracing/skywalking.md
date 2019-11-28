# Skywalking

go-chassis-apm is a plugin of go chassis, it reports tracing data to skywalking server

## Configurations
you must import go-chassis-apm plugin pkg in main.go
```go
import _ "github.com/go-chassis/go-chassis-apm/tracing/skywalking"
```

**In conf/monitoring.yaml**

**servicecomb.apm.tracing.tracer**
>  *(optional, string)* tracer'name, only skywalking now

**servicecomb.apm.tracing.settings.URI**
>  *(optional, string)* URI server address of skywalking

**servicecomb.apm.tracing.settings.servertype**
>  *(optional, string)* service type, match componentid in skywalking 
>  ex:  5001:servicecomb-mesher 5002:servicecomb-service-center 28:servicecomb-java-cahssis 

**servicecomb.apm.tracing.settings.enable**
>  *(optional, bool)* enable skywalking tracing ability

**Add handler name which are defined in github.com/go-chassis/go-chassis/core/handler**
>  Adding handler name *handler.SkyWalkingConsumer* in consumerChain.

>  Adding handler name *handler.SkyWalkingProvider* in providerChain.

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

```
  handler:
    chain:
      Consumer:
        outgoing:  #consumer handlers
        #ex: router, ratelimiter-consumer, bizkeeper-consumer, loadbalance, port-selector, transport, skywalking-consumer
      Provider:
        incoming:   #provider handlers
        #ex:  ratelimiter-provider, transport, skywalking-provider

```
