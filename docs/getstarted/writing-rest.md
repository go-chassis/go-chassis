Writing Rest service
==========================
服务端
1个工程或者go package，推荐结构如下

server/

├── conf

│ ├── chassis.yaml

│ └── microservice.yaml

└── main.go

1.编写接口
```go
type RestFulHello struct {}

func (r *RestFulHello) Sayhello(b *restful.Context) {
    b.Write([]byte("get user id: " + b.ReadPathParameter("userid")))
}
```
2.注册路由
```go
func (s *RestFulHello) URLPatterns() []restful.Route {
    return []restful.Route{
        {http.MethodGet, "/sayhello/{userid}", "Sayhello"},
    }
}
```
3.注册接口

第一个参数表示你要向哪个协议注册，第三个为schema ID,会在调用中使用

chassis.RegisterSchema("rest", &RestFulHello{}, server.WithSchemaId("RestHelloService"))
说明:

想注册的rest协议的接口，必须实现URLPatterns方法定义路由
路由中暴露为API的方法都要入参均为*restful.Context
3.修改配置文件chassis.yaml
```yaml
cse:
  service:
    registry:
      address: http://127.0.0.1:30100
  protocols:
    rest:
      listenAddress: 127.0.0.1:5001
```
4.修改microservice.yaml
```yaml
service_description:
  name: RESTServer
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

1.修改配置文件chassis.yaml
```yaml
cse:
  service:
    registry:
      address: http://127.0.0.1:30100
```
2.修改microservice.yaml

service_description:
  name: RESTClient
3.main中调用服务端，请求包括服务名，schema，operation及参数
```go
//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
    //Init framework
    if err := chassis.Init(); err != nil {
        lager.Logger.Error("Init failed.", err)
        return
    }
    req, _ := rest.NewRequest("GET", "cse://RESTServer/sayhello/world")
    defer req.Close()
    resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
    if err != nil {
        lager.Logger.Error("error", err)
        return
    }
    defer resp.Close()
    lager.Logger.Info(string(resp.ReadBody()))
}
```