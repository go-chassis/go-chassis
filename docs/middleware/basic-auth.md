# Basic Auth
go chassis提供高级别通用中间件抽象层，其中的一个抽象是Basic Auth，免于用户学习[handler chain内部的复杂性](https://docs.go-chassis.com/dev-guides/how-to-implement-handler.html)，让用户只需关注与自身业务的开发

## 使用
编写业务代码
```go
	basicauth.Use(&basicauth.BasicAuth{
		Realm: "test-realm",
		Authorize: func(user, password string) error {
		    //check your user name and password
		    return nil
		},
	})
```
更改配置文件, 将basicAuth handler添加到chain中，注意作为认证鉴权，一般说的都是服务端功能，所以要放到provider chain中
```yaml
cse:
  handler:
    chain:
      Provider:
        default: basicAuth
```