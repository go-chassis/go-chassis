# 动态配置
## 概述

go-chassis提供动态配置管理能力，支持CSE配置中心，本地文件，环境变量及命令行等多种配置管理，并由archaius包提供统一的接口获取配置。

## API

archaius包提供获取全部配置和四种获取指定配置值的Get方法，同时提供默认值，即若指定配置值未配置则使用默认值。另外提供Exist方法判断指定配置值是否存在。用户可使用UnmarshalConfig方法提供反序列化配置到结构体的方法。go-archaius内置5个配置源，获取配置时优先从优先级高的配置源获取指定配置项，若无则按优先级高低依次从后续配置源获取，直到遍历所有配置源。内置配置源生效的优先级由高到底分别是配置中心，命令行，环境变量，文件，外部配置源。

##### 获取指定配置值

```go
GetBool(key string, defaultValue bool) bool
GetFloat64(key string, defaultValue float64) float64
GetInt(key string, defaultValue int) int
GetString(key string, defaultValue string) string
```

```go
Get(key string) interface{}
Exist(key string) bool
```

##### 获取全部配置

```go
GetConfigs() map[string]interface{}
```

##### 反序列化配置到结构体

```go
UnmarshalConfig(obj interface{}) error
```

在go-archaius默认纳入动态管理的配置文件外，提供了AddFile方法允许用户添加其他文件到动态配置管理中。AddKeyValue可额外为外部配置源添加配置对。除默认加载的配置源外，允许用户实现自己的配置源，并通过RegisterListener注册到动态配置管理框架中。

##### 添加文件源

```go
AddFile(file string) error
```

##### 外部配置源添加配置对

```go
AddKeyValue(key string, value interface{}) error
```

##### 注册/注销动态监听

```go
RegisterListener(listenerObj core.EventListener, key ...string) error
UnRegisterListener(listenerObj core.EventListener, key ...string) error
```

在对接config center配置中心时请求中需指定demensionsInfo信息来确定获取配置的实例。该接口允许为配置项分区域DI配置和查询。

##### 添加DI及获取指定DI的配置值

```go
GetConfigsByDI(dimensionInfo string) map[string]interface{}
GetStringByDI(dimensionInfo, key string, defaultValue string) string
```

```go
AddDI(dimensionInfo string) (map[string]string, error)
```

## 示例

示例中文件配置如下，可通过archaius包的Get方法读取指定文件配置项。

```yaml
cse:
  fallback:
    Consumer:
      enabled: true
      maxConcurrentRequests: 20
```

```go
archaius.GetInt("cse.fallback.Consumer.Consumer.maxConcurrentRequests", 10)
archaius.GetBool("cse.fallback.Consumer.Consumer.enabled", false)
```



