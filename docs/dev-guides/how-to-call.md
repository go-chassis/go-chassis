# How to call micro service
提供2种风格调用，rpc与http

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


## rpc
```go
// declare reply struct
reply := &helloworld.HelloReply{}
// use WithProtocol to specify rpc plugin
core.NewRPCInvoker().Invoke(context.Background(), "RPCServer", "helloworld.Greeter", "SayHello",
		&helloworld.HelloRequest{Name: "Peter"}, reply, core.WithProtocol("grpc"))
```