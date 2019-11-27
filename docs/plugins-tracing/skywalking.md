# Skywalking

go-chassis-apm is a plugin of go chassis, it reports tracing data to skywalking server

## Configurations
**In conf/monitoring.yaml**

**servicecomb.apm.tracing.tracer**
>  *(optional, string)* tracer'name, only skywalking now

**servicecomb.apm.tracing.settings**
>  *(optional, map)* options including :
>  URI server address of skywalking 
>  servertype: service type, match componentid in skywalking ex:  5001:servicecomb-mesher 5002:servicecomb-service-center 28:servicecomb-java-cahssis 
>  enable: if open

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

```golang
consumerChain := strings.Join([]string{
    chassisHandler.Router,
    chassisHandler.RatelimiterConsumer,
    chassisHandler.BizkeeperConsumer,
    chassisHandler.Loadbalance,
    chassisHandler.SkyWalkingConsumer,
}, ",")

providerChain := strings.Join([]string{
    chassisHandler.RatelimiterProvider,
    chassisHandler.SkyWalkingProvider,
}, ",")
```
