# Restful Development

## 概述
go chassis 提供了Restful风格的http服务，
基于 [go-restful](https://github.com/emicklei/go-restful)，
版本号参考[go.mod](https://github.com/go-chassis/go-chassis/blob/master/go.mod)
使用方法请参阅下面的内容。

> yaml的设置请参考 [Get started](../getstarted/writing-rest.html)

## 路由注册
### 通过路由函数URLPatterns定义API route
基于go-restful，拥有它全部的API定义能力，可参考go-restful的[readme](https://github.com/emicklei/go-restful)来确认
```go
type DummyResource struct {
}

func (r *DummyResource) Sayhello(b *restful.Context) {
	id := b.ReadPathParameter("userid")
	b.Write([]byte(id))
}

//URLPatterns helps to respond for corresponding API calls
func (r *DummyResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/sayhello/{userid}", ResourceFunc: r.Sayhello,
			Returns: []*restful.Returns{{Code: 200}}},
	}
}
```
### Open API 文档管理
go chassis可以自动生成Open API2.0的文档，通过编写restful.Route结构的各个属性，可以丰富Open API文档内容
文档自动生成后,可以通过2种方式查看

1. 会被注册到service center中，可以通过服务id与schema id进行查询，schema id为微服务名
/v4/{project}/registry/microservices/{serviceId}/schemas/{schemaId}

2. 本地http服务启动后，可以在/apidocs.json路径查看文档，也可以在本地进程运行目录下的conf/schema文件夹下找到

## 路由分组

相同资源的多个路由可能会具备同样的路由前缀，这时候你可以使用路由分组来避免重复的路由前缀

```go
type DummyResource struct {
}

//GroupPath returns group path will auto 
func (r *DummyResource) GroupPath() string {
	return "/demo"
}

func (r *DummyResource) Sayhello(b *restful.Context) {
	id := b.ReadPathParameter("userid")
	b.Write([]byte(id))
}

func (r *DummyResource) Panic(b *restful.Context) {
	panic("panic msg")
}

//URLPatterns helps to respond for corresponding API calls
func (r *DummyResource) URLPatterns() []restful.Route {
	return []restful.Route{ 
		// will register path:/demo/sayhello
		{Method: http.MethodGet, Path: "/sayhello", ResourceFunc: r.Sayhello,
		    Returns: []*restful.Returns{{Code: 200}}},
        // will register path:/demo/sayhelloagain
        {Method: http.MethodGet, Path: "/sayhelloagain", ResourceFunc: r.Sayhello,
            Returns: []*restful.Returns{{Code: 200}}},
	}
}
```

