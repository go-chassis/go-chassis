# Restful Development

## 概述
go chassis 提供了Restful风格的http服务，基于 [go-restful](https://github.com/emicklei/go-restful)，使用方法请参阅下面的节点。

> yaml的设置请参考 [Get started](../getstarted/writing-rest.html)

## 路由注册
### 通过路由函数

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

### 通过函数名称（不推荐）

这是兼容旧版本的一种方式，由于实现使用了reflect，性能会有所损耗

```go
type DummyResource struct {
}

func (r *DummyResource) GroupPath() string {
	return "/demo"
}

func (r *DummyResource) Sayhello(b *restful.Context) {
	id := b.ReadPathParameter("userid")
	b.Write([]byte(id))
}

//URLPatterns helps to respond for corresponding API calls
func (r *DummyResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/sayhello/{userid}", ResourceFuncName: "Sayhello",
			Returns: []*restful.Returns{{Code: 200}}},
	}
}
```
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