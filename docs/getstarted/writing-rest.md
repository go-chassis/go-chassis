Writing Rest service
==========================
Checkout full example in [here](https://github.com/go-chassis/go-chassis/tree/master/examples/rest)
### Provider
this section show you how to write a http server

Create 1 project or go package as recommended 
```
server

+-- main.go

+-- conf

    +-- chassis.yaml

    +-- microservice.yaml
```

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
        {Method: http.MethodGet, Path: ""/sayhello/{userid}"", ResourceFunc: s.Sayhello,
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
```
client

+-- main.go

+-- conf
    
    +-- chassis.yaml
    
    +-- microservice.yaml
```
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


### Define the contents of the schema file in URLPatterns()

Here you can configure Method、Path、Parameters，and so on.
For example:

```go
func (d *Data) URLPatterns() []rf.Route {
	return []rf.Route{
		{
			Method:http.MethodGet,
			Path:"/price/{id}",
			ResourceFunc:d.GetPrice, #schema=operationId
			Consumes: []string{goRestful.MIME_JSON,goRestful.MIME_XML},
			Produces: []string{goRestful.MIME_JSON},
			Returns: []*rf.Returns{{Code: http.StatusOK,Message:"true",Model: Data{}}},
			Parameters:[]*rf.Parameters{#schema=parameter
				&rf.Parameters{"x-auth-token","string",goRestful.HeaderParameterKind,"this is a token"},
				&rf.Parameters{"x-auth-token2","string",goRestful.HeaderParameterKind,"this is a token"},
			},
		},
	}
}
````
you can find yor open API doc in http://ip:port/apidocs.json, after you start your program
,you can copy the response body into http://editor.swagger.io/ to read document.
```yaml
swagger: "2.0"
info:
  title: ""
  version: ""
basePath: /
paths:
  /price/{id}:
    get:
      operationId: GetPrice
      parameters:
      - name: x-auth-token
        in: header
        description: this is a token
        type: string
      - name: x-auth-token2
        in: header
        description: this is a token
        type: string
      consumes:
      - application/json
      - application/xml
      produces:
      - application/json
      responses:
        "200":
          description: "true"
          schema:
            $ref: '#/definitions/Data'
definitions:
  Data:
    type: object
    properties:
      err:
        $ref: '#/definitions/ErrorCode'
      priceID:
        type: string
      type:
        type: string
      value:
        type: string
  ErrorCode:
    type: object
    properties:
      code:
        type: integer
        format: ""
```

Paramater type：

```go
	// PathParameterKind = indicator of Request parameter type "path"
	PathParameterKind = iota

	// QueryParameterKind = indicator of Request parameter type "query"
	QueryParameterKind

	// BodyParameterKind = indicator of Request parameter type "body"
	BodyParameterKind

	// HeaderParameterKind = indicator of Request parameter type "header"
	HeaderParameterKind

	// FormParameterKind = indicator of Request parameter type "form"
	FormParameterKind
```

### Automatically generate schema file
The program will generate the schema file locally by default, 
If you want to define your own instead of automatically generating a schema，
You can modify the configuration '**noRefreshSchema: true**' in chassis.yaml

```yaml
cse:
  service:
    registry:
      address: http://127.0.0.1:30100 
  protocols:
    rest:
      listenAddress: "127.0.0.1:5003"
  noRefreshSchema: true
```


**Notice**
>> if conf folder is not under work dir, plz export CHASSIS_HOME=/path/to/conf/parent_folder or CHASSIS_CONF_DIR==/path/to/conf_folder
