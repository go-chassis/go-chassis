# Micro Service Definition 

## Introduction
Use microservice.yaml to describe your service

Conceptions:
- instance: one process is a micro service instance, instances belong to one micro service
- service: service is a static information entity in storage, it has instances

you can consider a project as an micro service, after compile, build and run, it became a micro service instance


## Configurations

**name**
> *(required, string)* Micro service name

**hostname**
> *(optional, string)* hostname of host, it can be IP, $INTERNAL_IP placeholder or hostname, default is hostname return by os.hostname()
> When specify `hostname: $INTERNAL_IP` go-chassis will report ip address instead of hostname to service center, this is useful when hostname is meaningless in some scenes, such as a docker host.

**APPLICATION_ID**
> *(optional, string)* Application ID, default value is "default"

**version**
> *(optional, string)* version number default is 0.0.1

**properties**
> *(optional, map)* micro service metadata ，usually it is defined in project, and never changed

**instance_properties**
> *(optional, map)* instance metadata, during runtime, if can be different based on environment 

**paths**
> *(optional, array)* micro service API paths, will be registered with servicecenter

## Example

```yaml
service_description:
  name: Server
  hostname: 10.244.1.3
  properties:
    project: X1
  instance_properties:
    nodeIP: 192.168.0.111
  paths:
  - path: /rest/demoservice
    property:
      checksession: true
```
