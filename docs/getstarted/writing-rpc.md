Writing gRPC service
==========================================
Checkout full example in [here](https://github.com/go-chassis/go-chassis-examples/tree/master/grpc)
### Define grpc contract
1个工程或者go package，推荐结构如下
```
schemas
`-- helloworld
    `-- helloworld.proto
```

1.定义helloworld.proto文件
```proto
syntax = "proto3";

package helloworld;

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}

```
2.通过pb生成go文件 helloworld.pb.go
```bash
protoc --go_out=plugins=grpc:. helloworld.proto
```

将生成的go文件拷贝到目录中

```
schemas
`-- helloworld
    |-- helloworld.pb.go
    `-- helloworld.proto
```

After generated, need to change one variable name
```go
var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.Greeter", // use this as the schemaID when consumer call provider
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "helloworld.proto",
}
```

change _Greeter_serviceDesc to *Greeter_serviceDesc*

### Provider Side
1个工程或者go package，推荐结构如下
```
server
|-- conf
|   |-- chassis.yaml
|   `-- microservice.yaml
`-- main.go
```

1.编写接口
```go
type Server struct{}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
```
2.注册接口

第一个参数表示你要向哪个协议注册，第三个为grpc serivce desc

```go
chassis.RegisterSchema("grpc", &Server{}, server.WithRPCServiceDesc(&pb.Greeter_serviceDesc))
```


3.修改配置文件chassis.yaml
```yaml
cse:
  service:
    registry:
      address: http://127.0.0.1:30100
  protocols:
    grpc:
      listenAddress: 127.0.0.1:5000
```
4.修改microservice.yaml
```yaml
service_description:
  name: RPCServer
```
5.In main.go import grpc server to enable grpc protocol and start go chassis
```go
import _ "github.com/go-chassis/go-chassis-protocol/server/grpc"
```

```go
func main() {
    //start all server you register in server/schemas.
    if err := chassis.Init(); err != nil {
        lager.Logger.Errorf("Init failed: %s", err)
        return
    }
    chassis.Run()
}
```
### Consumer Side
1个工程或者go package，推荐结构如下
```
client
|-- conf
|   |-- chassis.yaml
|   `-- microservice.yaml
`-- main.go
```

1.修改配置文件chassis.yaml
```yaml
cse:
  service:
    registry:
      address: http://127.0.0.1:30100
```

2.修改microservice.yaml
```yaml
service_description:
  name: Client
```

3.5.In main.go import grpc client to enable grpc protocol.
```go
import _ "github.com/go-chassis/go-chassis-protocol/client/grpc"
```

Use invoker to call remote function
```go
//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	//declare reply struct
	reply := &helloworld.HelloReply{}
	//Invoke with microservice name, schema ID and operation ID
	if err := core.NewRPCInvoker().Invoke(context.Background(), "RPCServer", "helloworld.Greeter", "SayHello",
		&helloworld.HelloRequest{Name: "Peter"}, reply, core.WithProtocol("grpc")); err != nil {
		lager.Logger.Error("error" + err.Error())
	}
	lager.Logger.Info(reply.Message)
}
```

**Notice**
>> if conf folder is not under work dir, plz export CHASSIS_HOME=/path/to/conf/parent_folder or CHASSIS_CONF_DIR==/path/to/conf_folder
