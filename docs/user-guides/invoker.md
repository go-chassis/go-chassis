# Invoker
## 概述

---

框架提供Rest调用与RPC调用2种方式

## API

#### Rest调用

使用NewRestInvoker创建一个invoker实例，可接受chain, filters等自定义选项

ContextDo可以接受一个http request作为参数，开发者可通过request的API对request进行操作，并作为参数传入该方法

```go
func NewRestInvoker(opt ...Option) *RestInvoker
func (ri *RestInvoker) ContextDo(ctx context.Context, req *rest.Request, options ...InvocationOption) (*rest.Response, error)
```

#### RPC调用

使用NewRPCInvoker创建invoker实例，可接受chain, filters等自定义选项

指定远端的服务名，struct name，以及func name，以及请求参数和返回接口即可进行调用

最终结果会赋值到reply参数中

```go
func NewRPCInvoker(opt ...Option) *RPCInvoker 
func (ri *RPCInvoker) Invoke(ctx context.Context, microServiceName, schemaID, operationID string, arg interface{}, reply interface{}, options ...InvocationOption) error
```

无论Rest还是RPC调用方法都能够接受多种选项对一次调用进行控制，参考options.go查看更多选项

## 示例

#### RPC

添加了2个具体调用选项，使用highway rpc，并使用roundrobin路由策略

```go
invoker.Invoke(ctx, "Server", "HelloServer", "SayHello",
    &helloworld.HelloRequest{Name: "Peter"},
    reply,
    core.WithProtocol("highway"),
    core.WithStrategy(loadbalance.StrategyRoundRobin),
    core.WithVersion("0.0.1"),

)
```

#### Rest

与普通的http调用不同的是url参数不使用ip:port而是服务名并且[http://变为cse://](http://变为cse://)

在初始化invoker时还指定了这次请求要经过的处理链名称custom

```go
req, _ := rest.NewRequest("GET", "cse://RESTServer/sayhello/world")
defer req.Close()
resp, err := core.NewRestInvoker(core.ChainName("custom")).ContextDo(context.TODO(), req)
```



