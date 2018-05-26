# 通信协议

## 概述
go-chassis支持不同协议扩展，不同的协议都可以接入到统一的微服务管理，控制，监控中。
目前支持HTTP1/2与RPC协议highway

## 配置

**protocols.{protocol_name}**
> *(required, string)* 协议名，目前内置rest与highway

**protocols.{protocol_name}.advertiseAddress**
> *(optional, string)* 协议广播地址，也就是向注册中心注册时的地址，在发现后进行通信时使用的网络地址

**protocols.{protocol_name}.listenAddress**
> *(required, string)* 协议监听地址，建议配置为0.0.0.0:{port}，
go chassis会自动为你计算advertiseAddress，无需手动填写，适合运行在容器中的场景，因为ip地址无法确定。

## 例子
```
cse:
  protocols:
    rest:
      listenAddress: 0.0.0.0:5000
    highway:
      listenAddress: 0.0.0.0:6000
```
