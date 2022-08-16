package main

import (
	"errors"
	"github.com/go-chassis/go-chassis/v2/middleware/jwt"
	"net/http"
	"strings"

	"github.com/go-chassis/go-chassis/v2"
	_ "github.com/go-chassis/go-chassis/v2/middleware/jwt"
	"github.com/go-chassis/go-chassis/v2/security/token"
	rf "github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/go-chassis/openlog"
)

// if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/server/
func main() {
	chassis.RegisterSchema("rest", &HelloAuth{})

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
	//start all server you register in server/schemas.
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}

type User struct {
	Name string `json:"name"`
	Pwd  string `json:"password"`
}
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

func (r *HelloAuth) Access(b *rf.Context) {
	b.Write([]byte("success"))
}

// URLPatterns helps to respond for corresponding API calls
func (r *HelloAuth) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodPost, Path: "/login", ResourceFunc: r.Login,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodGet, Path: "/resource", ResourceFunc: r.Access,
			Returns: []*rf.Returns{{Code: 200}}},
	}
}
