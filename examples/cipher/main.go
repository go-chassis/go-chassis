package main

import (
	"github.com/go-chassis/cari/security"
	"github.com/go-chassis/go-chassis/v2"
	_ "github.com/go-chassis/go-chassis/v2/middleware/ratelimiter"
	"github.com/go-chassis/go-chassis/v2/security/cipher"
	"github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/go-chassis/openlog"
	"net/http"
)

type DemoResource struct {
}

func (r *DemoResource) Limit(b *restful.Context) {
	b.ReadResponseWriter().WriteHeader(http.StatusOK)
	d, _ := cipher.Decrypt("ok")
	b.ReadResponseWriter().Write([]byte(d))
}

// URLPatterns returns routes
func (r *DemoResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/decrypt", ResourceFunc: r.Limit},
	}
}

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/{project_root}/

func main() {
	cipher.InstallCipherPlugin("default", new)
	chassis.RegisterSchema("rest", &DemoResource{})
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}

//DefaultCipher is a struct
type DefaultCipher struct {
}

func new() security.Cipher {
	return &DefaultCipher{}
}

//Encrypt is method used for encryption
func (c *DefaultCipher) Encrypt(src string) (string, error) {
	return src, nil
}

//Decrypt is method used for decryption
func (c *DefaultCipher) Decrypt(src string) (string, error) {
	return "d: " + src, nil
}
