# Invoker
## Introduction

---

Invoker is the entry point for a developer to call remote service

## API

#### Rest Invoker

使用NewRestInvoker创建一个invoker实例，可接受chain等自定义选项

ContextDo可以接受一个http request作为参数，开发者可通过request的API对request进行操作，并作为参数传入该方法

```go
func NewRestInvoker(opt ...Option) *RestInvoker
func (ri *RestInvoker) ContextDo(ctx context.Context, req *http.Request, options ...InvocationOption) (*rest.Response, error)
```

#### RPC Invoker

使用NewRPCInvoker创建invoker实例，可接受chain等自定义选项

指定远端的服务名，struct name，以及func name，以及请求参数和返回接口即可进行调用

最终结果会赋值到reply参数中

```go
func NewRPCInvoker(opt ...Option) *RPCInvoker 
func (ri *RPCInvoker) Invoke(ctx context.Context, microServiceName, schemaID, operationID string, arg interface{}, reply interface{}, options ...InvocationOption) error
```

无论Rest还是RPC调用方法都能够接受多种选项对一次调用进行控制，参考options.go查看更多选项

## Examples

#### RPC


```go
invoker.Invoke(ctx, "Server", "HelloServer", "SayHello",
    &helloworld.HelloRequest{Name: "Peter"},
    reply,

)
```

#### Rest


在初始化invoker时还指定了这次请求要经过的处理链名称custom

```go
req, _ := rest.NewRequest("GET", "http://RESTServer/sayhello/world")
defer req.Close()
resp, err := core.NewRestInvoker(core.ChainName("custom")).ContextDo(context.TODO(), req)
```

#### Multiple Port
if you define different port for the same protocol, like below
```go
cse:
  protocols:
    rest:
      listenAddress: 0.0.0.0:5000
    rest-admin:
      listenAddress: 0.0.0.0:5001
```
then you can use suffix "admin" as port to access rest-admin server
```go
req, _ := rest.NewRequest("GET", "http://RESTServer:admin/sayhello/world")
```
use only service name to access rest server
```go
req, _ := rest.NewRequest("GET", "http://RESTServer/sayhello/world")
```




