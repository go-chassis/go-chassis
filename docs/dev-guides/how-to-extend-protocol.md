# Protocol
## 概述

---

框架支持[grpc协议](https://github.com/go-chassis/go-chassis-protocol)，
用户可扩展自己的RPC协议，并使用RPCInvoker调用

## 如何实现

---

#### 客户端

* 实现协议的客户端接口

```go
type Client interface
```

* 实现以下接口来返回客户端插件

```go
func(...clientOption.Option) Client
```

* 安装客户端插件

```go
func InstallPlugin(protocol string, f ClientNewFunc)
```

* 处理链默认自带名为transport的handler，他将根据协议名加载对应的协议客户端,指定协议的方式如下

```go
invoker.Invoke(ctx, "Server", "HelloServer", "SayHello",
    &helloworld.HelloRequest{Name: "Peter"},
    reply,
    core.WithProtocol("grpc"),
)
```

#### 服务端

* 实现协议的Server端

```
type Server interface
```

* 修改配置文件以启动协议监听

```
cse:
  protocols:
    grpc:
      listenAddress: 127.0.0.1:5000
      advertiseAddress: 127.0.0.1:5000
```



