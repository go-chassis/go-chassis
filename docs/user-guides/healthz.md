# 客户端健康检查

## 概述

客户端健康检查（Health Check）是指客户端对服务端实例缓存进行健康性的判断。

在网络分区或延时较大的环境下，客户端可能会出现上报心跳到服务中心失败的情况，导致结果是，客户端会将收到的实例下线事件，并移除本地实例缓存，最终影响业务调用。

为防止上述情况发生，go-chassis提供这样的健康检查机制：当客户端监听到某个（或latest）版本的服务端可用实例下降到0个时，
客户端会在同步移除实例缓存前进行一次健康检查，调用服务端暴露的健康检查接口（RESTful或Highway）并校验其返回值来确认是否要移除实例缓存。

## 配置

### 服务端注册

go-chassis默认不会主动注册服务端的健康检查接口，需要用户主动import到项目中。

```go
// 注册健康检查接口
import _ "github.com/ServiceComb/go-chassis/healthz/provider"
```

加入上述代码片段后，go-chassis会按照暴露的服务协议类型对应注册健康检查接口，接口描述如下

* RESTful:

  1. Method: GET
  1. Path: /healthz
  1. Response:
  ```js
  {
    "appId": "string",
    "serviceName": "string",
    "version": "string"
  }
  ```

* Highway: 

  1. Schema: _chassis_highway_healthz
  1. Operation: HighwayCheck
  1. Response: 
  ```proto
  // The response message containing the microservice key
  message Reply {
    string appId = 1;
    string serviceName = 2;
    string version = 3;
  }
  ```

### 客户端配置

客户端健康检查配置在chassis.yaml。

**healthCheck**
> *(optional, bool)* 允许对服务端的实例做健康检查，默认值为false。

###### 示例

```yaml
cse:
  service:
    Registry:
      healthCheck: true
      #serviceDiscovery:
      #  healthCheck: true # 同时支持单独开启服务发现能力时的客户端健康检查
```