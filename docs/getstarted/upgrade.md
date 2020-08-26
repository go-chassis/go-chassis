# Upgrade from 1.8.3 to 2.0

## micro service definition
1.Migrate config from "service_description" to "servicecomb.service", for example:
1.8
```yaml
service_description:
  name: Server
  hostname: 10.244.1.3
```
2.0,
```yaml
servicecomb:
  service:
    name: Server
    hostname: 10.244.1.3
```

2. change "instance_properties" to "instanceProperties", for example:
1.8
```yaml
service_description:
  name: Server
  instance_properties:
    nodeIP: 192.168.0.111
```
2.0
```yaml
servicecomb:
  service:
    name: Server
    instanceProperties:
      nodeIP: 192.168.0.111
```

## change cse:// to http://

for example:

1.8
```go
arg, _ := rest.NewRequest("GET", "cse://Server/instances", nil)
```
2.0
```go
arg, _ := rest.NewRequest("GET", "http://Server/instances", nil)
```

## change all "cse:"" to "servicecomb:" in yaml
for example:

1.8
```yaml
cse:
 config:
```
2.0
```yaml
servicecomb:
 config:
```

## move "registry" under "servicecomb"

for example:

1.8
```yaml
cse:
 service:
   registry:
```
2.0
```yaml
servicecomb:
 registry:
```