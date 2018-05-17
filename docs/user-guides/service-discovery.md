# 服务发现
## 概述

Service discovery 是关于如何发现服务的配置。
和Registry的区别是，他仅负责发现服务，而不负责注册服务
Service Discovery与Registry只能选择其一进行配置
启用此功能可以与Istio的Pilot集成

## 配置

服务发现的配置在chassis.yaml。

**type**
> *(optional, string)* 对接服务中心插件类型，默认为servicecenter

**address**
> *(optional, bool)*服务中心地址 允许配置多个以逗号隔开，默认为空

**refreshInterval**
> *(optional, string)* 更新实例缓存的时间间隔，格式为数字加单位（s/m/h），如1s/1m/1h，默认为30s

## 示例

```yaml
cse:
  service:
    Registry:
      serviceDiscovery:
        type: pilot      
        address: http://istio-pilot.istio-system:8080
        refeshInterval : 30s                    
```



