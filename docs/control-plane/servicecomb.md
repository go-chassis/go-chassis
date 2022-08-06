# ServiceComb

ServiceComb service center is the default plugin of go chassis for service discovery.
ServiceComb kie is the default plugin of go chassis for configuration management.


## Configurations
Update the chassis.yaml of your microservices with the following configuration.
```yaml
servicecomb:
  registry:
      #type: servicecenter
      address: http://127.0.0.1:30110
      refeshInterval : 30s       
      watch: true
  config:
    client:
      serverUri: http://127.0.0.1:30110         # This should be the address of your Kie Server
      #type: kie
  credentials:
    account:
      name: service_account
      password: Complicated_password1
    cipher: default
```

## Config Servers
Go-Chassis leverage [go-archaius](https://github.com/go-chassis/go-archaius) to retrieve configuration from remote server

### ServiceComb-Kie

[kie](https://github.com/apache/servicecomb-kie) is a service for configuration management in distributed system.

use this [guide](https://kie.readthedocs.io/en/latest/get-started.html) to set up a ServiceComb-Kie server.

