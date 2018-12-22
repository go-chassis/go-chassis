Writing Rest service
==========================
Checkout full example in [here](https://github.com/go-chassis/go-chassis/tree/master/examples/rest)
### Provider
this section show you how to write a http server

Create 1 project or go package as recommended 

server/

├── conf

│ ├── chassis.yaml

│ └── microservice.yaml

└── main.go

1.Write a struct to hold http logic and url patterns
```go
type RestFulHello struct {}

func (r *RestFulHello) Sayhello(b *restful.Context) {
    b.Write([]byte("get user id: " + b.ReadPathParameter("userid")))
}
```
2.Write your url patterns
```go
func (s *RestFulHello) URLPatterns() []restful.Route {
    return []restful.RouteSpec{
        {Method: http.MethodGet, Path: ""/sayhello/{userid}"", ResourceFuncName: "Sayhello",
         			Returns: []*rf.Returns{{Code: 200}}},
    }
}
```
3.Modify chassis.yaml to describe the server you need
```yaml
cse:
  service:
    registry:
      address: http://127.0.0.1:30100 
  protocols: # what kind of server you want to launch
    rest: #launch a http server
      listenAddress: 127.0.0.1:5001
```

4.Register this struct

the first params means which server you want to register your struct to. Be aware API can separate by different server and ports
```go
chassis.RegisterSchema("rest", &RestFulHello{})
```

**Notice**
>>Must implement URLPatterns, and for other functions must use \*restful.Context as the only input, 
and certainly the method name must start with uppercase


5.Modify microservice.yaml
```yaml
service_description:
  name: RESTServer # name your provider
```
6.In main.go init and start the chassis 
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
### Consumer
this section show you how to write a http client

Create 1 project or go package as recommended 

client/

├── conf

│ ├── chassis.yaml

│ └── microservice.yaml

└── main.go

1. modify chassis.yaml
```yaml
cse:
  service:
    registry:
      address: http://127.0.0.1:30100
```
2. modify microservice.yaml
```yaml
service_description:
  name: RESTClient #name your consumer
```
3.in main.go call your service
```go
func main() {
    //Init framework
    if err := chassis.Init(); err != nil {
        lager.Logger.Error("Init failed.", err)
        return
    }
    req, _ := rest.NewRequest("GET", "http://RESTServer/sayhello/world")
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
**Notice**
>> if conf folder is not under work dir, plz export CHASSIS_HOME=/path/to/conf/parent_folder or CHASSIS_CONF_DIR==/path/to/conf_folder
