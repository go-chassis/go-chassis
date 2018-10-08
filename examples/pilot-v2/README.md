## Istio pilot as Discovery Service for go-chassis

go-chassis can be integrated with Istio Pilot for service discovery. To enable Istio pilot support in go-chassis, the simple 2 steps are needed during development:

1. Import the istiov2 registry plugin of mesher

```go
import _ "github.com/go-mesh/mesher/plugins/registry/istiov2"
```

   This will install the istiov2 service discovery plugin at runtime.

2. Configure serviceDiscovery in chassis.yaml

```yaml
cse:
  service:
    registry:
      registrator:
        disabled: true
      serviceDiscovery:
        type: pilotv2
        address: grpc://istio-pilot.istio-system:15010
```

Disable the registrator(since we don't have to register the service to Pilot explicitly) and change serviceDiscovery type to `pilotv2`(indicates the Pilot that provides xDS v2 API), configure the address, typically `istio-pilot.istio-sytem:15010` in a Istio environment.



### Deploy the Kubernetes resources

All the changes mentioned above are already done in the demo project. The image building and Kubernetes resources deployment are already obtained in the Makefile(go explore it if you are interested), are few steps are needed:

```bash
$ cd server
$ make docker
......
Successfully tagged go-chassis/pilotv2server:latest
```

Build the client image with the same make command. Then we can make the image accessible to the nodes. We can either copy image tar files to nodes and load it with docker, or push the image to a public docker registry service(in this case, update the image reference in the Kubernetes resources yamls).



Then we can deploy the client and server to Kubernetes cluster with `make deploy`. The resources are listed and described below:

| Resources                                      | Description                                                  |
| ---------------------------------------------- | :----------------------------------------------------------- |
| client/build/pilotv2client.dep.yaml            | The client deployment                                        |
| client/build/istio-desitnationrule-reader.yaml | The cluter role and binding for DesitnationRule read permissions for the default service account |
| server/build/k8s/destinationrule.yaml          | The DestinationRules for the 3 versions of servers           |
| server/build/k8s/pilotv2server.v*.dep.yaml     | The deployment for the 3 versions of servers                 |
| server/build/k8s/service.yaml                  | The server's Kubernetes service                              |

And the deploy command:

```bash
$ cd server
make deploy
```

Run the same command under folder client, then make sure the client and servers are up and running, and trace the client's log:

```bash
$ kubectl get pods
NAME                                READY     STATUS        RESTARTS   AGE
pilotv2client-fcdbc4d9c-w52sk       1/1       Running       0          32s
pilotv2server-v1-6fc9bdc7bf-vhplh   1/1       Running       0          27s
pilotv2server-v2-76d47f567-27pcx    1/1       Running       0          27s
pilotv2server-v3-6f4c755d9c-mhv7n   1/1       Running       0          27s
$ kubectl logs -f pilotv2client-fcdbc4d9c-w52sk
......
2018-09-26 07:17:18.414 +00:00 INFO client/client_manager.go:86 Create client for rest:pilotv2server:172.30.61.5:5001
2018-09-26 07:17:18.416 +00:00 INFO client/main.go:37 REST Server sayhello[GET]: user world from 31
2018-09-26 07:17:23.416 +00:00 INFO client/client_manager.go:86 Create client for rest:pilotv2server:172.30.61.7:5001
2018-09-26 07:17:23.420 +00:00 INFO client/main.go:37 REST Server sayhello[GET]: user world from 84
2018-09-26 07:17:28.420 +00:00 INFO client/client_manager.go:86 Create client for rest:pilotv2server:172.30.61.6:5001
2018-09-26 07:17:28.423 +00:00 INFO client/main.go:37 REST Server sayhello[GET]: user world from 84
```

The server responses the client with a random number which is created when server starts. So client gets the response from different servers with a different random number. So now we have the client/server up and running, with a simple router rule(defined in client/conf/router.yaml)!