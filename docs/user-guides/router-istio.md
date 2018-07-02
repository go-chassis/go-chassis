# How to use istio to manage route

Instead of using CSE and route config to manage route, go-chassis supports istio as a control plane to set route rule and follows the envoy API reference to manage route. This page gives the examples to show how requests are routed between micro services.

## Go-chassis Configurations

In **Consumer** router.yaml, you can set router.infra to define which router plugin go-chassis fetches from.  The default router.infra  is cse, which means the routerule comes from route config in CSE config-center. If router.infra is set to be pilot, the router.address is necessary, such as the in-cluster istio-pilot grpc address.

```yaml
router:
  infra: pilot # pilot or cse
  address: http://istio-pilot.istio-system:15010
```

In **Both** consumer and provider registry configurations, the recommended one shows below.

```yaml
cse:
  service:
    registry:
      registrator:
        disabled: true
      serviceDiscovery:
        type: pilot
        address: http://istio-pilot.istio-system:8080
```

## Kubernetes Configurations

The provider applications of v1, v2 and v3 version could be deployed in kubernetes cluster as **Deployment** with differenent labels. The labels of version is necessary now,  and you need to set env to generate nodeID in istio system, such as **POD_NAMESPACE, POD_NAME** and **INSTANCE_IP**.

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    version: v1
    app: pilot
    name: istioserver
  name: istioserver-v1
  namespace: default
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: pilot
      version: v1
      name: istioserver
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: pilot
        version: v1
        name: istioserver
    spec:
      containers:
      - image: gosdk-istio-server:latest
        imagePullPolicy: Always
        name: istioserver-v1
        ports:
        - containerPort: 8084
          protocol: TCP
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        env:
        - name: CSE_SERVICE_CENTER
          value: http://istio-pilot.istio-system:8080
        - name: POD_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
        - name: INSTANCE_IP
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.podIP
        volumeMounts:
        - mountPath: /etc/certs/
          name: istio-certs
          readOnly: true
      dnsPolicy: ClusterFirst
      initContainers:
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - name: istio-certs
        secret:
          defaultMode: 420
          optional: true
          secretName: istio.default
```

## Istio v1alpha3 router configurations

 [Traffic-management](https://istio.io/docs/tasks/traffic-management/request-routing/) gives references and examples of istio new router rule schema. First, subsets is defined according to labels. Then you can set route rule of differenent weight for virtual services.

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: istioserver
spec:
  host: istioserver
  subsets:
  - name: v1
    labels:
      version: v1
  - name: v2
    labels:
      version: v2
  - name: v3
    labels:
      version: v3
```

> **NOTICE: The subsets only support labels of version to distinguish differenent virtual services, this constrains will canceled later.**

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: istioserver
spec:
  hosts:
    - istioserver
  http:
  - route:
    - destination:
        host: istioserver
        subset: v1
      weight: 25
    - destination:
        host: istioserver
        subset: v2
      weight: 25
    - destination:
        host: istioserver
        subset: v3
      weight: 50

```