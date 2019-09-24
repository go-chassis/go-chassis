# Contract management
## 概述
go chassis follow [Open API 2.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md)
you can manually edit documentation and put in go chassis schema folders, 
then go chassis will automatically upload them to service center.

in additional, if you write rest service, go chassis will automatically generate Open API spec, 
and upload them to service center.

## Configuration

契约文件必须为yaml格式文件，契约文件应放置于go-chassis的schema目录。

schema目录位于：

1，conf/{serviceName}/schema，其中conf表示go-chassis的conf文件夹

2，${SCHEMA\_ROOT}

2的优先级高于1。

## API

包路径

```go
import "github.com/go-chassis/go-chassis/core/config/schema"
```

契约字典，key值为契约文件名，value为契约文件内容

```go
var DefaultSchemaIDsMap map[string]string
```



## Example
the contract file structure is as below
    conf
    `-- myservice
        `-- schema
            |-- myschema1.yaml
            `-- myschema2.yaml


define API doc in URLPatterns function 
check https://github.com/go-chassis/go-chassis/blob/master/server/restful/router.go for more options

```go

func (r *RestFulHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/", ResourceFunc: r.Root,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodGet, Path: "/sayhello/{userid}", ResourceFunc: r.Sayhello,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/sayhi", ResourceFunc: r.Sayhi,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/sayjson", ResourceFunc: r.SayJSON,
			Returns: []*rf.Returns{{Code: 200}}},
	}
}
```
