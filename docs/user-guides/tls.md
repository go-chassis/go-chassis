# TLS
## 概述

用户可以通过配置TLS启动HTTPS通信，保障数据的安全传输。包括客户端与服务端TLS，通过配置来自动启用Consumer与Provider的TLS配置。

## 配置

tls可配置于chassis.yaml文件或单独的tls.yaml文件。在如下格式中，tag指明服务名以及服务类型，key指定对应configuration的配置项。

```yaml
ssl:
  [tag].[key]: [configuration]
```

### TAG

tag为空时ssl配置为公共配置。registry.consumer及configcenter.consumer是作为消费者访问服务中心和配置中心时的ssl配置。protocol.serviceType允许协议和类型的任意组合。name.protocol.serviceType在协议和类型的基础上可定制服务名。

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

**{protocol}.{serviceType}**
>协议为任意协议目前包括 *grpc*，*rest*，用户扩展协议后，即可使用新的协议配置。
>类型为*Consumer*,*Provider* |

**{name}.{protocol}.{serviceType}**
>定制某微服务的独有的TLS配置 name为微服务名

### KEY

ssl支持以下配置项，其中若私钥KEY文件加密，则需要指定加解密插件及密码套件等信息进行解密。

                       
**keyFile**
> *(optional, string)* RSA Private Key file path

**verifyPeer**
>*(optional, bool)* 
是否验证对端,默认*false*

**cipherSuits**
> *(optional, string)* *TLS\_ECDHE\_RSA\_WITH\_AES\_128\_GCM\_SHA256*, *TLS\_ECDHE\_RSA\_WITH\_AES\_256\_GCM\_SHA384*
> 密码套件                           |

**protocol**
> *(optional, string)* TLS协议的最小版本,默认为*TLSv1.2*

**caFile**
> *(optional, string)* if verifyPeer is true, you need to supply ca files in here
as a consumer, you need server cert files, as a provider, it needs client cert files
check [example](https://github.com/go-chassis/go-chassis-examples/tree/master/mutualtls)

**certFile**
> *(optional, string)* Certificate file path

**certPwdFile**
> *(optional, string)* a file path, this file's content is Passphrase of keyFile

**cipherPlugin**
> *(optional, string)* you can custom 
[Cipher](https://docs.go-chassis.com/dev-guides/how-to-write-cipher.html) 
to decrypt "certPwdFile" content, by default no decryption        

## API

通过为Provider和Consumer配置ssl，go-chassis会自动为其加载相关配置。用户也可以通过chassis暴露的接口直接使用相关API。以下API主要用于获取ssl配置以及tls.Config。

##### 获取默认SSL配置

```go
GetDefaultSSLConfig() *common.SSLConfig
```

##### 获取指定SSL配置

```go
GetSSLConfigByService(svcName, protocol, svcType string) (*common.SSLConfig, error)
```

##### 获取指定TLSConfig

```go
GetTLSConfigByService(svcName, protocol, svcType string) (*tls.Config, *common.SSLConfig, error)
```

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
  rest.Provider.caFile: client.crt
  rest.Provider.certPwdFile: pwd 
```

### Consumer config
set verifyPeer to true to tell go chassis to verify TLSService 
add server.crt to caFile, it will be used as root CA during verification
```yaml
ssl:
  TLSService.rest.Consumer.cipherSuits: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  TLSService.rest.Consumer.protocol: TLSv1.2
  TLSService.rest.Consumer.caFile: server.crt
  TLSService.rest.Consumer.certFile: client.crt
  TLSService.rest.Consumer.keyFile: client.key
  TLSService.rest.Provider.verifyPeer: true
```