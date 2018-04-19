# Handler
## 概述

Go chassis以插件的形式支持在一次请求调用中，插入自己的处理逻辑。

## 实现

实现handler需要以下三个步骤：实现Handler接口，并根据其名称注册，最后在chassis.yaml中添加handler chain相关配置。其中\[service\_type\] 可配置为Provider或Consumer，\[chain\_name\]默认为default。

##### 注册处理逻辑

```go
RegisterHandler(name string, f func() Handler) error
```

##### 实现处理接口

```go
type Handler interface {
    Handle(*Chain, *invocation.Invocation, invocation.ResponseCallBack)
    Name() string
}
```

##### 添加配置

```yaml
cse:
  handler:
    chain:
      [service_type]:
        [chain_name]: [your_handler_name]
```

## 示例

示例中注册的是名为fake-handler的处理链，其实现的Handle方法仅记录inv的endpoint信息。

```go
package handler
import (
    "github.com/ServiceComb/go-chassis/core/handler"
    "github.com/ServiceComb/go-chassis/core/invocation"
    "log"
)
const Name = "fake-handler"
type FakeHandler struct{}

func init()                         { handler.RegisterHandler(Name, New) }
func New() handler.Handler          { return &FakeHandler{} }
func (h *FakeHandler) Name() string { return Name }

func (h *FakeHandler) Handle(chain *handler.Chain, inv *invocation.Invocation,
    cb invocation.ResponseCallBack) {
    log.Printf("fake handler running for %v", inv.Endpoint)
    chain.Next(inv, func(r *invocation.InvocationResponse) error {
        return cb(r)
    })
}
```

chassis.yaml配置示例如下

```yaml
cse:
  handler:
    chain:
      Provider:
        default: fake-handler
```



