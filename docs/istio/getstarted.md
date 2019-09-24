Get started
=====

go-chassis can be integrated with Istio for service discovery and routing. To enable Istio pilot support in go-chassis, the simple 2 steps are needed during development:

- **Import the istiov2 registry plugin from mesher**

```go
import _ "github.com/go-mesh/mesher/plugins/registry/istiov2"
```

   This will install the istiov2 service discovery plugin at runtime.

- **Configure service discovery in chassis.yaml**

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

Disable the registrator(since we don't have to register the service to Pilot explicitly) and change serviceDiscovery type to `pilotv2`(indicates the Pilot that provides xDS v2 API, the xDS v1 API is already deprecated), configure the address, typically `istio-pilot.istio-system:15010` in a Istio environment.

Then when deploying the micro services in Istio, make sure the Kubernetes Services' name and go-chassis service name is exactly the same, then go-chassis will discovery the service instances as expected from Pilot.


### The routing tags in Istio

In the original go-chassis configuration, user can specify tag based route rules, as described below:

```yaml
## router.yaml
router:
  infra: cse
routeRule:
  targetService:
    - precedence: 2
      route:
      - tags:
          version: v1
        weight: 40
      - tags:
          version: v2
          debug: true
        weight: 40
      - tags:
          version: v3
        weight: 20
```
Then in a typical Istio environment, which is likely to be Kubernetes cluster, user can specify the DestinationRules for targetService with the same tags:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: targetService
spec:
  host: targetService
  subsets:
  - name: v1
    labels:
      version: v1
  - name: v2
    labels:
      version: v2
      debug: "true"
  - name: v3
    labels:
      version: v3
```

Notice that the subsets' tags are the same with those in `router.yaml`, then go-chassis's tag based load balancing strategy works as it originally does.
