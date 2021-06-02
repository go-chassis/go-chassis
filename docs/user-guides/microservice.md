# Micro Service Definition 

## Introduction
Use microservice.yaml to describe your service

Conceptions:
- instance: one system process is a micro service instance, instances belong to one micro service
- service: service is a static information entity in storage, it has instances

you can consider a project as a micro service, after build, ship and run, it becomes a micro service instance


## Configurations

**name**
> *(required, string)* microservice name. 
> it is the registered name in registry service.
> In best practice 
> this name is never changed in your entire software lifecycle(decided by a developer).


**hostname**
> *(optional, string)* 
> hostname of host, it can be IP, $INTERNAL_IP placeholder or hostname, default is hostname return by os.hostname()
> When specify `hostname: $INTERNAL_IP` go-chassis will report ip address instead of hostname to service center, 
> this is useful when hostname is meaningless in some scenes, such as a docker host.

**app**
> *(optional, string)* application id, default value is "default".
> In best practice you can build and run your service in different system.
> so better to decide it in runtime(decided by an operator)

**version**
> *(optional, string)* 
> version number, default is 0.0.1

**properties**
> *(optional, map)* 
> microservice metadata, In best practice it is defined in project, and never changed

**instanceProperties**
> *(optional, map)* instance metadata, it can be different in runtime(decided by an operator)

**paths**
> *(optional, array)* microservice API paths, will be registered to service center

**schemas**
>*(optional, array)* schema id, which will be registered to service center

## Example

```yaml
servicecomb:
  service:
    name: Server
    hostname: 10.244.1.3
    properties:
      project: X1
    instanceProperties:
      nodeIP: 192.168.0.111
    paths:
      - path: /rest/demoservice
    schemas:
      - "schema"
      - "schema1"
      - "schema2"
```
