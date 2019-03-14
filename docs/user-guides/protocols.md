# Protocol Servers

## Introduction
you can extend your own protocol in go chassis, currently support rest(http) and gRPC

## Configurations

**protocols.{protocol_server_name}**
> *(required, string)* the name of the protocol server, it must be protocol name or consist of protocol name and a suffix.
 the suffix and protocol is connect with hyphen "-" like <protocol>-{suffix}

**protocols.{protocol_server_name}.advertiseAddress**
> *(optional, string)* server advertise address, if you use registry like service center, 
this address will be registered in registry, so that other service can discover your address

**protocols.{protocol_server_name}.listenAddress**
> *(required, string)* server listen address, recommend to use 0.0.0.0:{port}, 
then go chassis will automatically generate advertise address, it is convenience to run in container
 because the internal IP is not sure until container runs



## Example
this config will launch 2 http server and 1 grpc server
```
cse:
  protocols:
    rest:
      listenAddress: 0.0.0.0:5000
    rest-admin:
      listenAddress: 0.0.0.0:5001
    grpc:
      listenAddress: 0.0.0.0:6000
```

for ipv6, need quotation marks. because [] is object list in yaml format
```
cse:
  protocols:
    rest:
      listenAddress: "[2407:c080:17ff:ffff::7274:83a]:5000"
```

if you do not want to specify a port, you can leave the port empty (use quotes) or use 0, the system will give a random port for you
```
cse:
  protocols:
    rest:
      listenAddress: 0.0.0.0:0
    grpc:
      listenAddress: "127.0.0.1:"
```