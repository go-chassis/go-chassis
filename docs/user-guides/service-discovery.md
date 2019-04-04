# Service Discovery
## 概述

Service discovery 是关于如何发现服务的配置。
和Registry的区别是，他仅负责发现服务，而不负责注册服务
Service Discovery与Registry只能选择其一进行配置
启用此功能可以与Istio的Pilot集成

## 配置

服务发现的配置在chassis.yaml。

**type**
> *(optional, string)* 对接服务中心插件类型，默认为servicecenter，另外可选择pilotv2以及kube

**NOTE: 当使用kube registry时，发布的service需要指定port name为以下格式 [protocol]-[suffix]**

**address**
> *(optional, bool)*服务中心地址 允许配置多个以逗号隔开，默认为空

**refreshInterval**
> *(optional, string)* 更新实例缓存的时间间隔，格式为数字加单位（s/m/h），如1s/1m/1h，默认为30s

## 示例

当registry type为pilotv2时需要指定pilot的地址address，当registry type为kube时需要指定与kube-apiserver交互所需的kubeconfig的配置文件位置，以下分别为registry的最小示例。

```yaml
cse:
  service:
    Registry:
      serviceDiscovery:
        type: pilotv2
        address: grpc://istio-pilot.istio-system:15010
        refeshInterval : 30s
```

```yaml
cse:
  service:
    Registry:
      serviceDiscovery:
        type: kube
        configPath: /etc/.kube/config
```

