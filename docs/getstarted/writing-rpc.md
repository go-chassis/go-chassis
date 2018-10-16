Writing gRPC service
==========================================
### Define grpc contract
1个工程或者go package，推荐结构如下

schemas

├── helloworld

│ ├──helloworld.proto

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

protoc --go_out=. helloworld.proto
将生成的go文件拷贝到目录中

schemas

├── helloworld

│ ├──helloworld.proto

│ └──helloworld.pb.go

After generated, need to change one variable name
```go
var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.Greeter",
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

server/

├── conf

│ ├── chassis.yaml

│ └── microservice.yaml

└── main.go

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
chassis.RegisterSchema("grpc", &Server{}, server.WithGRPCServiceDesc(&pb.Greeter_serviceDesc))
```


3.修改配置文件chassis.yaml
```yaml
cse:
  service:
    registry:
      address: http://127.0.0.1:30100
  protocols:
    highway:
      listenAddress: 127.0.0.1:5000
```
4.修改microservice.yaml
```yaml
service_description:
  name: Server
```
5.main.go中启动服务
```go
func main() {
    //start all server you register in server/schemas.
    if err := chassis.Init(); err != nil {
        lager.Logger.Error("Init failed.", err)
        return
    }
    chassis.Run()
}
```
### Consumer Side
1个工程或者go package，推荐结构如下

client/

├── conf

│ ├── chassis.yaml

│ └── microservice.yaml

└── main.go

1.拿到pb文件生成go代码

protoc --go_out=. hello.proto
2.修改配置文件chassis.yaml

```yaml
cse:
  service:
    registry:
      address: http://127.0.0.1:30100
```
3.修改microservice.yaml
```yaml
service_description:
  name: Client
```
4.main中调用服务端，指定微服务名，schema，operation与参数和返回
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