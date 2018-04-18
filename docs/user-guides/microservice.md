# 微服务定义

## 概述
每个微服务都具有自己的定制的microservice.yaml文件，用来配置微服务私有信息。
其中properties以key-value对的形式为微服务添加非通用的属性，
比如allowCrossApp: false可用于配置是否允许跨应用发现和访问。


## 配置

**name**
> *(required, string)* 微服务名称

**version**
> *(optional, string)* 版本号 默认0.0.1

**properties**
> *(optional, map)* 微服务元数据

## 例子

```yaml
service_description:
  name: Server
  properties:
    yourkey: test
    otherkey: test2
```