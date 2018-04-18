# 服务发现
## 概述

Service discovery 是关于服务如何发现的配置，当在平台运行时，可以与Istio的Pilot集成，框架不负责注册工作，而只使用服务发现功能

## 配置

服务发现的配置在chassis.yaml。

| 配置项 | 默认值 | 配置说明 |
| --- | --- | --- |
| type | servicecenter | 对接服务发现系统类型 支持pilot,file,servicecenter |
| address | [http://127.0.0.1:30100](http://127.0.0.1:30100) | 服务中心地址 允许配置多个以逗号隔开 |
| refeshInterval | 30s | 更新实例缓存的时间间隔，格式为数字加单位（s/m/h），如1s/1m/1h |
| watch | false | 是否watch实例变化事件 |
| api.version | v4 | 访问服务中心的api版本 |

## 示例

```yaml
cse:
  service:
    serviceDiscovery:
      type: pilot        #optional: 默认类型为对接服务中心   
      address: http://10.0.0.1:30100,http://10.0.0.2:30100 
      refeshInterval : 30s                    
      api:
        version: v4
```



