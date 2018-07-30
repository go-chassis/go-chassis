# Archaius
## 概述

go-archaius是go-chassis的动态配置框架，目前支持CSE 配置中心，本地文件，ENV，CMD等配置管理。如果用户希望接入自己的配置服务中，可以参考本章节实现。

## 使用说明

go-archaius支持同时配置多种源，包括命令行，环境变量，外部源及文件等。用户可通过ConfigurationFactory接口的AddSource方法添加自己的配置源，并通过RegisterListener方法注册EventListener。

```go
AddSource(core.ConfigSource) error
```

```go
RegisterListener(listenerObj core.EventListener, key ...string) error
```

其中配置源需要实现ConfigSource接口。其中GetPriority和GetSourceName必须实现且有有效返回，分别用于获取配置源优先级和配置源名称。GetConfigurations和GetConfigurationByKey方法用于获取全部配置和指定配置项，需要用户实现。其他方法可以返回空。

- GetPriority方法用于确定配置源的优先级，go-archaius内置的5个配置源优先级由高到底分别是配置中心，命令行，环境变量，文件，外部配置源，对应着0到4五个整数值。用户自己接入的配置源可自行配置优先级级别，数值越小则优先级越高。GetSourceName方法用于返回配置源名称。
- 若没有区域区分的集中式配置中心，DemensionInfo相关接口可返回nil不实现。
- Cleanup用于清空本地缓存的配置。
- DynamicConfigHandler接口可根据需要实现，用于实现动态配置动态更新的回调方法。

```go
type ConfigSource interface {
    GetSourceName() string
    GetConfigurations() (map[string]interface{}, error)
    GetConfigurationsByDI(dimensionInfo string) (map[string]interface{}, error)
    GetConfigurationByKey(string) (interface{}, error)
    GetConfigurationByKeyAndDimensionInfo(key, dimensionInfo string) (interface{}, error)
    AddDimensionInfo(dimensionInfo string) (map[string]string, error)
    DynamicConfigHandler(DynamicConfigCallback) error
    GetPriority() int
    Cleanup() error
}
```

注册EventListener用于在配置源更新时由Dispatcher分发事件，由注册的listener处理。

```go
type EventListener interface {
    Event(event *Event)
}
```

## 示例

##### 实现configSource

```go
type fakeSource struct {
	Configuration  map[string]interface{}
	changeCallback core.DynamicConfigCallback
	sync.Mutex
}

func (*fakeSource) GetSourceName() string { return "TestingSource" }
func (*fakeSource) GetPriority() int      { return 0 }

func (f *fakeSource) GetConfigurations() (map[string]interface{}, error) {
	config := make(map[string]interface{})
	f.Lock()
	defer f.Unlock()
	for key, value := range f.Configuration {
		config[key] = value
	}
	return config, nil
}

func (f *fakeSource) GetConfigurationByKey(key string) (interface{}, error) {
	f.Lock()
	defer f.Unlock()
	configValue, ok := f.Configuration[key]
	if !ok {
		return nil, errors.New("invalid key")
	}
	return configValue, nil
}

func (f *fakeSource) DynamicConfigHandler(callback core.DynamicConfigCallback) error {
	f.Lock()
	defer f.Unlock()
	f.changeCallback = callback
	return nil
}

func (f *fakeSource) Cleanup() error {
	f.Lock()
	defer f.Unlock()
	f.Configuration = make(map[string]interface{})
	f.changeCallback = nil
	return nil
}

func (*fakeSource) AddDimensionInfo(d string) (map[string]string, error) { return nil, nil }
func (*fakeSource) GetConfigurationByKeyAndDimensionInfo(k, d string) (interface{}, error) { return nil, nil }
func (*fakeSource) GetConfigurationsByDI(d string) (map[string]interface{}, error) { return nil, nil }

```

##### 添加configSource

```go
func NewConfigSource() core.ConfigSource {
    return &fakeSource{
      Configuration: make(map[string]interface{}),
    }
}
```

```go
factory, _ := goarchaius.NewConfigFactory(lager.Logger)
err := factory.AddSource(fakeSource.NewConfigSource())
```

##### 注册evnetListener

```go
type EventHandler struct{
    Factory goarchaius.ConfigurationFactory
}
func (h EventHandler) Event(e *core.Event) { 
  value := h.Factory.GetConfigurationByKey(e.Key)
  log.Printf("config value after change %s | %s", e.Key, value)
}
```

```go
factory.RegisterListener(&EventHandler{Factory: factory}, "a*")
```

