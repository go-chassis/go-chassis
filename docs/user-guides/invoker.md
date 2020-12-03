# Invoker
## Introduction
Invoker is the entry point to call remote service.
support 2 kinds of invoker, rpc and http
## http
支持传入原生net/http提供的request，将request传入invoker即可将请求交由go chassis内核处理，
最终返回net/http原生response。
```go
invoker:= core.NewRestInvoker()
req, _ := rest.NewRequest("GET", "http://orderService/hello", nil)
invoker.ContextDo(context.TODO(), req)
```
如果你不想引入service center,consul这类注册发现组件，甚至连kubernetes的发现插件也不想用，
而是想使用kubernetes的原生容器网络，或者结合任意服务网格，可传入WithoutSD, 另外需要添加端口信息。
```go
req, _ := rest.NewRequest("GET", "http://orderService:8080/hello", nil)
invoker.ContextDo(context.TODO(), req, core.WithoutSD())
```

#### RPC Invoker

使用NewRPCInvoker创建invoker实例，可接受chain等自定义选项

指定远端的服务名，struct name，以及func name，以及请求参数和返回接口即可进行调用

最终结果会赋值到reply参数中
```go
// declare reply struct
reply := &helloworld.HelloReply{}
// use WithProtocol to specify rpc plugin
core.NewRPCInvoker().Invoke(context.Background(), "RPCServer", "helloworld.Greeter", "SayHello",
		&helloworld.HelloRequest{Name: "Peter"}, reply, core.WithProtocol("grpc"))
```
无论Rest还是RPC调用都能够接受多种选项对一次调用进行控制，参考options.go查看更多选项

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
```yaml
servicecomb:
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




