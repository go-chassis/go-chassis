# TLS
## 概述
```eval_rst
.. image:: tls.png
```
as in this figure, provider distributes its cert to consumer1 and 2, 
it means both consumer trust this provider, because they add provider cert to their trust bundle.
consumer1 also gives its cert to provider, 
provider add it to its trust bundle so that provider allow consumer1's request.
but provider does not has consumer2's cert. so consumer2 is not allowed to call provider.
by setting this trust bundle, you can simply protect your services.
## 配置

use tls.yaml to set tls config
the format is as below, will explain what is tag and key
role has only 2 type, Consumer and Provider
```yaml
ssl:
  [tag].[role].[key]: [configuration]
```

### tag and role
tag indicates what is the tls config target.   

registry.Consumer, configcenter.Consumer etc is build-in config to define the tls settings for 
control plane services(like service center, config center), you can not use them.

you can custom tls settings by following rules. 

tag usually comprises of service name, role(Consumer or Provider) and protocol. 

**registry.Consumer**
> 服务注册中心TLS配置

**serviceDiscovery.Consumer**
> 服务发现TLS配置

**contractDiscovery.Consumer**
> 契约发现TLS配置

**registrator.Consumer**
> 服务注册中心TLS配置

**configcenter.Consumer**
>配置中心TLS配置                                     |



**{protocol}.{Consumer|Provider}**
>define which protocol should use TLS communication

**{service_name}.{protocol}.{Consumer}**
>only works for consumer, it means using TLS communication to call serviceXX

###  key


                       
**Provider.keyFile**
> *(required, string)* RSA Private Key file path for server

**Provider.certFile**
> *(required, string)* Certificate file path for server

**Consumer.keyFile**
> *(optional, string)* RSA Private Key file path for client

**Consumer.certFile**
> *(optional, string)* Certificate file path for client

**{Consumer|Provider}.verifyPeer**
>*(optional, bool)* 
verify the other service or not, default is false.

**{Consumer|Provider}.cipherSuits**
> *(optional, string)* *TLS\_ECDHE\_RSA\_WITH\_AES\_128\_GCM\_SHA256*, *TLS\_ECDHE\_RSA\_WITH\_AES\_256\_GCM\_SHA384*
> 密码套件                           |

**{Consumer|Provider}.protocol**
> *(optional, string)* TLS protocol version, default is *TLSv1.2*

**{Consumer|Provider}.caFile**
> *(optional, string)* Define trust CA bundle in here. if verifyPeer is true, 
you must supply ca file list in here. 
During communication as a consumer, you need to add server cert files.
as a provider, it need to add client cert files
check [example](https://github.com/go-chassis/go-chassis-examples/tree/master/mutualtls)



**{Consumer|Provider}.certPwdFile**
> *(optional, string)* a file path, this file's content is Passphrase of keyFile, 
if you set Passphrase for you keyFile, you must set this config

**{Consumer|Provider}.cipherPlugin**
> *(optional, string)* you can custom 
[Cipher](https://docs.go-chassis.com/dev-guides/how-to-write-cipher.html) 
to decrypt "certPwdFile" content, by default no decryption        

## Example1: Simple TLS communication

### Generate files for a service
1. you can generate private key file with Passphrase 
```bash
#generate priviate key with passphrase
openssl genrsa -des3 -out server.key 1024
# save your passphrase
echo {your Passphrase} > pwd
```
or without passphrase
```bash
#generate private key without passphrase 
openssl genrsa -out server.key 2048
```

2. you can sign cert with csr and key 
```bash
openssl req -new -key server.key -out server.csr
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

```
or only with key
```bash
openssl req -new -x509 -key server.key -out server.crt -days 3650
```
### Provider配置

以下为rest类型provider提供HTTPS访问的ssl配置，其中tag为protocol.serviceType的形式。

```yaml
ssl:
  rest.Provider.cipherPlugin: default
  rest.Provider.cipherSuits: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  rest.Provider.protocol: TLSv1.2
  rest.Provider.keyFile: server.key
  rest.Provider.certFile: server.crt
  rest.Provider.certPwdFile: pwd # include Passphrase
```

### Consumer配置

以下为访问rest类型服务的消费者的ssl配置。tag为name.protocol.serviceType的形式，
其中TLSService为要访问的服务名，rest为协议。


```yaml
ssl:
  TLSService.rest.Consumer.cipherPlugin: default
  TLSService.rest.Consumer.cipherSuits: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  TLSService.rest.Consumer.protocol: TLSv1.2
```

## Example2: Mutual TLS communication
check complete [example](https://github.com/go-chassis/go-chassis-examples/tree/master/mutualtls)
### Generate client cert file
```bash
openssl genrsa -out client.key 2048
openssl req -new -x509 -key client.key -out client.crt -days 3650

```

### Provider config
that define: as a provider, which client can call this it
set verifyPeer to true to verify all clients. 
add client.crt in caFile, it will be used as client CA during verification
```yaml
ssl:
  rest.Provider.cipherPlugin: default
  rest.Provider.cipherSuits: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  rest.Provider.protocol: TLSv1.2
  rest.Provider.keyFile: server.key
  rest.Provider.certFile: server.crt
  rest.Provider.verifyPeer: true
  rest.Provider.caFile: client.crt,xxx.crt
  rest.Provider.certPwdFile: pwd 
```

### Consumer config
that define: as a consumer, how to call a service which enabled TLS config
the provider's name is TLSService.
set verifyPeer to true to tell go chassis to verify TLSService during communication
add provider's server.crt to caFile, it will be used as root CA during verification
```yaml
ssl:
  TLSService.rest.Consumer.cipherSuits: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  TLSService.rest.Consumer.protocol: TLSv1.2
  TLSService.rest.Consumer.caFile: server.crt,xxx.crt
  TLSService.rest.Consumer.certFile: client.crt
  TLSService.rest.Consumer.keyFile: client.key
  TLSService.rest.Provider.verifyPeer: true
```