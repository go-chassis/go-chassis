# ServiceComb

ServiceComb service center is the default plugin of go chassis, 
it support client side discovery, so need to set registry service. 
it implements both ServiceDiscovery and Registrator plugin.

Update the chassis.yaml of your microservices with the following configuration.

## Configurations

```yaml
servicecomb:
  registry:
      type: servicecenter
      address: http://10.0.0.1:30100,http://10.0.0.2:30100 
      refeshInterval : 30s       
      watch: true                         
      api:
        version: v4
```

## Config Servers

Go-Chassis supports retrieving configuration from different remote servers and can be used for dynamic configurations of microservices. The following will guide you in configuring remote configurations services ServiceComb-Kie.

### ServiceComb-Kie

[kie](https://github.com/apache/servicecomb-kie) is a service for configuration management in distributed system.

User can use this [guide](https://kie.readthedocs.io/en/latest/get-started.html) to start up the start a ServiceComb-Kie server.

Update the chassis.yaml of your microservices with the following configuration.

```yaml
servicecomb:
  config:
    client:
      serverUri: http://127.0.0.1:30110         # This should be the address of your Kie Server
      type: kie                                 # The type should be kie
      refreshMode: 1                            # If Refresh Mode is set to 1, chassis will pull the configuration periodically. If Refresh Mode is set to 0, chassis will use a long connection to watch the configuration changes and update immediately when the configuration changes.
      refreshInterval: 10                       # If Refresh Mode is set to 1, chassis retrieves the configurations from Kie at this interval. If Refresh Mode is set to 0, chassis uses this time as the timeout for long pulling connections.
```

Update the microservice.yaml of your microservices with the following configuration.

```yaml
servicecomb:
  service:
    name: servicename        # your microservices name
    version: 0.1.0           # your microservices version
    environment: prod        # the environment where your application runs
```
