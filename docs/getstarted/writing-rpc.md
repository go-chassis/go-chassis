Writing RPC service
==========================================
定义请求与返回结构体
1个工程或者go package，推荐结构如下

schemas

├── helloworld

│ ├──helloworld.proto

1.定义helloworld.proto文件
```proto
syntax = "proto3";

package helloworld;


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

protoc --go_out=. hello.proto
将生成的go文件拷贝到目录中

schemas

├── helloworld

│ ├──helloworld.proto

│ └──helloworld.pb.go

服务端
1个工程或者go package，推荐结构如下

server/

├── conf

│ ├── chassis.yaml

│ └── microservice.yaml

└── main.go

1.编写接口
```go
type HelloServer struct {
}
func (s *HelloServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
    return &helloworld.HelloReply{Message: "Go Hello  " + in.Name}, nil
}
```
2.注册接口

第一个参数表示你要向哪个协议注册，第三个为schema ID,会在调用中使用

chassis.RegisterSchema("highway", &HelloServer{}, server.WithSchemaId("HelloService"))
说明:

想暴露为API的方法都要符合以下条件

第一个参数为context.Context

第二个参数必须是结构体指针

返回的第一个必须是结构体指针

返回的第二个为error

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
客户端
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
        lager.Logger.Error("Init failed.", err)
        return
    }
    //declare reply struct
    reply := &helloworld.HelloReply{}
    //Invoke with microservice name, schema ID and operation ID
    if err := core.NewRPCInvoker().Invoke(context.Background(), "Server", "HelloService", "SayHello", &helloworld.HelloRequest{Name: "Peter"}, reply); err != nil {
        lager.Logger.Error("error", err)
    }
    lager.Logger.Info(reply.Message)
}
```