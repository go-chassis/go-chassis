package config_test

import (
	"github.com/go-chassis/cari/security"
	"github.com/go-chassis/go-chassis/v2/security/cipher"
	"testing"

	"os"
	"path/filepath"
	"time"

	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/pkg/util/fileutil"
	"github.com/stretchr/testify/assert"
)

func TestGetRegistratorRbacAccount(t *testing.T) {
	b := []byte(`
servicecomb:
  registry:
    disabled: false            #optional: 默认开启registry模块
    type: servicecenter        #optional: 默认类型为对接服务中心
    address: http://10.0.0.1:30100,http://10.0.0.2:30100
    register: auto             #optional：默认为自动 [auto manual]
    refeshInterval: 30s
    watch: true
    uploadSchema: false 
    heartbeat:
      mode: non-keep-alive
      interval: 30s
  credentials:
    account:
      name: service_account
      password: Complicated_password1
    cipher: default
`)
	d, _ := os.Getwd()
	filename1 := filepath.Join(d, "chassis.yaml")
	f1, err := os.OpenFile(filename1, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	assert.NoError(t, err)
	_, err = f1.Write(b)
	assert.NoError(t, err)

	os.Setenv(fileutil.ChassisConfDir, d)
	defer os.Unsetenv(fileutil.ChassisConfDir)
	err = config.Init()
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	cipher.InstallCipherPlugin("default", new)
	config.ReadGlobalConfigFromArchaius()
	c := config.GetRegistratorRbacAccount()
	assert.Equal(t, "service_account", c.Username)
	assert.Equal(t, "Complicated_password1", c.Password)
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
