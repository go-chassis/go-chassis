# TLS
## 概述

用户可以通过配置SSL/TLS启动HTTPS通信，保障数据的安全传输。包括客户端与服务端TLS，通过配置来自动启用Consumer与Provider的TLS配置。

## 配置

tls可配置于chassis.yaml文件或单独的tls.yaml文件。在如下格式中，tag指明服务名以及服务类型，key指定对应configuration的配置项。

```yaml
ssl:
  [tag].[key]: [configuration]
```

#### TAG

tag为空时ssl配置为公共配置。registry.consumer及configcenter.consumer是作为消费者访问服务中心和配置中心时的ssl配置。protocol.serviceType允许协议和类型的任意组合。name.protocol.serviceType在协议和类型的基础上可定制服务名。

| 标签名                       | 配置说明                                     |
| ------------------------- | ---------------------------------------- |
| N/A                       | 公共配置                                     |
| registry.consumer         | 服务注册中心                                   |
| configcenter.consumer     | 配置中心                                     |
| protocol.serviceType      | 协议: \[highway或rest\]  类型: \[Consumer或Provider\] |
| name.protocol.serviceType | 定制标签                                     |

#### KEY

ssl支持以下配置项，其中若私钥KEY文件加密，则需要指定加解密插件及密码套件等信息进行解密。

| 配置项          | 默认值                                      | 配置说明                           |
| ------------ | ---------------------------------------- | ------------------------------ |
| cipherPlugin | default                                  | 指定加解密插件 内部插件支持 \[default aes\] |
| verifyPeer   | false                                    | 是否验证对端                         |
| cipherSuits  | TLS\_ECDHE\_RSA\_WITH\_AES\_128\_GCM\_SHA256, TLS\_ECDHE\_RSA\_WITH\_AES\_256\_GCM\_SHA384 | 密码套件                           |
| protocol     | TLSv1.2                                  | TLS协议的最小版本                     |
| caFile       |                                          | ca文件路径                         |
| certFile     |                                          | 私钥cert文件路径                     |
| keyFile      |                                          | 私钥key文件路径                      |
| certPwdFile  |                                          | 私钥key加密的密码文件                   |

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

## 示例

### Provider配置

以下为rest类型provider提供HTTPS访问的ssl配置，其中tag为protocol.serviceType的形式。

```yaml
ssl:
  rest.Provider.cipherPlugin: default
  rest.Provider.verifyPeer: true
  rest.Provider.cipherSuits: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  rest.Provider.protocol: TLSv1.2
  rest.Provider.keyFile: /etc/ssl/server_key.pem
  rest.Provider.certFile: /etc/ssl/server.cer
  rest.Provider.certPwdFile: /etc/ssl/cert_pwd_plain
  rest.Provider.caFile: /etc/ssl/trust.cer
```

### Consumer配置

以下为访问rest类型服务的消费者的ssl配置。tag为name.protocol.serviceType的形式，其中Server为要访问的服务名，rest为协议。verifyPeer若配置为true将启动双向认证，否则客户端将忽略对服务端的校验。

```yaml
ssl:
  Server.rest.Consumer.cipherPlugin: default
  Server.rest.Consumer.verifyPeer: true
  Server.rest.Consumer.cipherSuits: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  Server.rest.Consumer.protocol: TLSv1.2
  Server.rest.Consumer.keyFile: /etc/ssl/server_key.pem
  Server.rest.Consumer.certFile: /etc/ssl/server.cer
  Server.rest.Consumer.certPwdFile: /etc/ssl/cert_pwd_plain
  Server.rest.Consumer.caFile: /etc/ssl/trust.cer
```



