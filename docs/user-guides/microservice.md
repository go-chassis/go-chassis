# 微服务定义

## 概述
一个微服务需要通过microservice.yaml文件对自己进行定义。
微服务与微服务实例概念：
每个进程对应一个微服务实例，多个实例可能属于同一个微服务。
比如：以程序为单位，编译打包后运行便产生了一个微服务实例，他们的代码是同源的，也就是说他们属于同一个微服务


## 配置

**name**
> *(required, string)* 微服务名称

**APPLICATION_ID**
> *(optional, string)* 所属应用 默认default

**version**
> *(optional, string)* 版本号 默认0.0.1

**properties**
> *(optional, map)* 微服务元数据，通常在开发期就已经定死

**instance_properties**
> *(optional, map)* 微服务实例元数据，运行期每个实例的内容可能会根据运行环境而有差异

## 例子

```yaml
service_description:
  name: Server
  properties:
    project: test
  instance_properties:
    nodeIP: 192.168.0.111
```