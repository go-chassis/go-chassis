# Registry
## 概述

微服务的注册发现默认通过[service center](https://github.com/apache/servicecomb-service-center)完成。
用户可以配置与服务中心的通信方式，服务中心地址，以及自身注册到服务中心的信息。
微服务启动过程中，会自动向服务中心进行注册。
在微服务运行过程中，go-chassis会周期从服务中心查询其他服务的实例信息缓存到本地.

## 配置

注册中心相关配置分布在两个yaml文件中，分别为chassis.yaml和microservice.yaml文件。
chassis.yaml中配置使用的注册中心类型、注册中心的地址信息。


**disabled**
> *(optional, bool)* 是否开启服务注册发现模块，默认为false

**type**
> *(optional, string)* 对接的服务中心类型，默认为servicecenter，
> 也支持用户定制registry[插件](https://github.com/go-chassis/go-chassis-extension/tree/master/registry)。


**address**
> *(optional, bool)*服务中心地址 允许配置多个以逗号隔开，默认为空

**register**
> *(optional, bool)* 默认为 auto，可选manual。即默认为自动注册，框架启动时完成微服务实例的自动注册。
> 当配置manual时，框架只会注册配置文件中的微服务，不会注册微服务实例。
> 使用者可以通过服务中心的API完成实例注册。是否自动自注册，

**refreshInterval**
> *(optional, string)* 更新实例缓存的时间间隔，格式为数字加单位（s/m/h），如1s/1m/1h，默认为30s

**watch**
> *(optional, bool)*  是否watch实例变化事件，默认为false

## 示例

服务中心最简化配置只需要registry的address, 或者disabled设置为true。

```yaml
servicecomb:
  registry:
    disabled: false            #optional: 默认开启registry模块
    type: servicecenter        #optional: 默认类型为对接服务中心
    address: http://10.0.0.1:30100,http://10.0.0.2:30100 
    register: auto             #optional：默认为自动 [auto manual]
    refeshInterval: 30s       
    watch: true                         
  credentials:
    account:
      name: service_account  
      password: Complicated_password1
    cipher: default
```





