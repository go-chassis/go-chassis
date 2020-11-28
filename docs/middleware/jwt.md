# JWT
JWT中间件对Token进行认证，验证请求是否合法
## 使用
**编写业务逻辑：**
暴露一个API，允许申请token, 假设用户名和密码是admin就是合法用户，为他生成一个token
代码中的secret为密钥，建议商用系统进行妥善的安全设计
```go
type HelloAuth struct {
}

func (r *HelloAuth) Login(b *rf.Context) {
	u := &User{}
	if err := b.ReadEntity(u); err != nil {
		b.WriteError(http.StatusInternalServerError, err)
		return
	}
	if u.Name == "admin" && u.Pwd == "admin" {
		to, err := token.DefaultManager.Sign(map[string]interface{}{
			"user": u.Name,
			"pwd":  u.Pwd,
		}, []byte("my_secret"))
		if err != nil {
			b.WriteError(http.StatusInternalServerError, err)
		}
		b.Write([]byte(to))
	} else {
		b.WriteError(http.StatusInternalServerError, errors.New("wrong user or pwd"))
	}

}

```
定制token的认证逻辑，放开login api的访问
```go
	jwt.Use(&jwt.Auth{
		MustAuth: func(req *http.Request) bool {
			if strings.Contains(req.URL.Path, "/login") {
				return false
			}
			return true
		},
		Realm: "test-realm",
	    SecretFunc: func(claims interface{}, method token.SigningMethod) (interface{}, error) {
                   			return []byte("my_secret"), nil
                   		},
	})
```

更改配置文件, 将basicAuth handler添加到chain中，注意作为认证鉴权，一般说的都是服务端功能，所以要放到provider chain中
```yaml
servicecomb:
  handler:
    chain:
      Provider:
        default: jwt
```

验证与完整代码https://github.com/go-chassis/go-chassis/tree/master/examples/jwt
