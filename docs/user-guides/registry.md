# Registry
## 概述

微服务的注册发现默认通过服务中心完成。用户可以配置与服务中心的通信方式，服务中心地址，以及自身注册到服务中心的信息。微服务启动过程中，会自动向服务中心进行注册。在调用其他服务的时候，go-chassis会从服务中心查询其他服务的实例信息缓存到本地，比如访问地址、使用的协议以及元数据等。在微服务运行过程中，go-chassis自动更新实例在客户端的缓存，更新缓存时间可配置，也可以使用watch机制监听实例变化的事件。

## 配置

注册中心的配置分布在两个yaml文件中，分别为chassis.yaml和microservice.yaml文件。chassis.yaml中配置APPLICATION\_ID信息、使用的注册中心类型、注册中心的地址信息。

* type: 对接的服务中心类型默认加载servicecenter和file两种，registry也支持用户定制registry并注册。
* scope: 默认不允许跨应用间访问，只允许本应用间访问，当配置为full时则允许跨应用间访问，且能发现本租户全部微服务。
* autodiscovery: 是否开启自动发现，开启后将从服务中心注册发现注册的其他服务中心。
* register: 配置项默认为自动注册，即框架启动时完成实例的自动注册。当配置manual时，框架只会注册配置文件中的微服务，不会注册实例，使用者可以通过服务中心对外的API完成实例注册。
* api.version: 目前只支持v4版本。


**disabled**
> *(optional, bool)* 是否开启服务注册发现模块，默认为false

**type**
> *(optional, string)* 对接服务中心插件类型，默认为servicecenter

**scope**
> *(optional, bool)* 默认为空，不允许跨app发现，填入full以允许跨应用发现

**autodiscovery**
> *(optional, bool)* 自动发现 服务中心集群节点 默认为false

**address**
> *(optional, bool)*服务中心地址 允许配置多个以逗号隔开，默认为空

**register**
> *(optional, bool)* 是否自动自注册，默认为 auto，可选manual

**refeshInterval**
> *(optional, string)* 更新实例缓存的时间间隔，格式为数字加单位（s/m/h），如1s/1m/1h，默认为30s

**api.version**
> *(optional, string)* 访问服务中心的api版本，默认为v4

**watch**
> *(optional, bool)*  是否watch实例变化事件，默认为false




## API

Registry提供以下接口供用户注册微服务及实例。以下方法通过Registry提供的相关接口实现，内置有两种Registry的实现，默认为servicecenter。另外支持用户自行定义实现Registry接口的插件，用于服务注册发现。

##### 注册微服务实例

```go
RegisterMicroserviceInstances() error
```

##### 注册微服务

```go
RegisterMicroservice() error
```

##### 自定义Registry插件

```go
InstallPlugin(name string, f func(opts ...Option) Registry)
```

## 示例

服务中心最简化配置只需要registry的address，注册的微服务实例通过appId、服务名和版本决定。

```yaml
APPLICATION_ID: CSE #optional
cse:
  service:
    registry:
      disabled: false            #optional: 默认开启registry模块
      type: servicecenter        #optional: 默认类型为对接服务中心
      scope: full                #optional: scope为full注册时允许跨app，并且发现本租户全部微服务
      autodiscovery: true       
      address: http://10.0.0.1:30100,http://10.0.0.2:30100 
      register: auto             #optional：默认为自动 [auto manual]
      refeshInterval : 30s       
      watch: true                         
      api:
        version: v4
```




