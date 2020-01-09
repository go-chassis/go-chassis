# skywalking
is an essential tool to help optimize and monitor the performance of your microservice, it reports tracing data to skywalking server

## Configurations

1.init skywalking manager

**Calling function Init() in tracing/skywalking.go manually to init skywalking manager.**

2.edit conf/monitoring.yaml

**servicecomb.apm.tracing.settings.URI**
>  *(optional, string)* URI server address of skywalking

**servicecomb.apm.tracing.settings.servertype**
>  *(optional, string)* service type, match componentid in skywalking 
>  ex:  5001:servicecomb-mesher 5002:servicecomb-service-center 28:servicecomb-java-cahssis 

**servicecomb.apm.tracing.settings.enable**
>  *(optional, bool)* enable skywalking tracing ability

3.Add add this in provider chain and consumer chain

**skywalking-consumer and skywalking-provider**

4.import tracing package
```go
// should import after import go-chassis
	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/middleware/tracing"
```

## Example
```yaml
servicecomb:
  apm:                                #application performance monitor
    tracing:
      settings:
        enable: true                  #enable tracing ability
        URI: 127.0.0.1:11800          #url of skywalking 
        serverType: 5001              #server type
```

```yaml
handler:
  chain:
    Provider:
      default: skywalking-provider
    Consumer:
      default: skywalking-consumer
```



