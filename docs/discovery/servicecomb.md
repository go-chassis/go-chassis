# ServiceComb

ServiceComb service center is the default plugin of go chassis, it support client side discovery, so need to set registry service. 
it implements both ServiceDiscovery and Registrator plugin

## Configurations

```yaml
cse:
  service:
    registry:
      type: servicecenter
      address: http://10.0.0.1:30100,http://10.0.0.2:30100 
      refeshInterval : 30s       
      watch: true                         
      api:
        version: v4
```
