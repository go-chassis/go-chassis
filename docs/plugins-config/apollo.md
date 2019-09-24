# Ctrip Apollo
[Ctrip Apollo](https://github.com/ctripcorp/apollo) is a Configuration Server which can be used to store your configurations. Go-Chassis supports retrieving the configurations from Apollo and can be used for dynamic configurations of microservices. In this guide we will explain you how to configure Go-Chassis to use Apollo as a configuration server.

## Configurations
Use can use this [guide](https://github.com/ctripcorp/apollo/wiki) to start up the Ctrip Apollo and make the Project, NamesSpace and add Configurations to it. Once your Apollo Server is setup then you can do the following modification in Go-Chassis to make it work with Apollo.  
Update the chassis.yaml of your microservices with the following configuration.
```yaml
cse:
  config:
    client:
      serverUri: http://127.0.0.1:8080          # This should be the address of your Apollo Server
      type: apollo                              # The type should be Apollo
      refreshMode: 1                            # Refresh Mode should be 1 so that Chassis-pulls the Configuration periodically
      refreshInterval: 10                       # Chassis retrives the configurations from Apollo at this interval
      serviceName: apollo-chassis-demo          # This the name of the project in Apollo Server
      env: DEV                                  # This is the name of environment to which configurations belong in Apollo
      cluster: demo                             # This is the name of cluster to which your Project belongs in Apollo
      namespace: application                    # This is the NameSpace to which your configurations belong in the project.
```
Once these configurations are set the Chassis can retrieve the configurations from Apollo Server.  
To see the detailed use case of how to use Ctrip Apollo with Chassis please refer to this [example](https://github.com/asifdxtreme/chassis-apollo-example).