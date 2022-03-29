package config_test

import (
	"os"
	"testing"

	_ "github.com/go-chassis/go-chassis/v2/initiator"

	"github.com/go-chassis/go-chassis/v2/core/config"

	"io"
	"path/filepath"
	"time"

	"github.com/go-chassis/go-chassis/v2/pkg/util/fileutil"
	"github.com/stretchr/testify/assert"
)

func TestCDInit(t *testing.T) {
	b := []byte(`
servicecomb:
  registry:
      #disabled: false           optional:禁用注册发现选项，默认开始注册发现
      type: servicecenter           #optional:可选zookeeper/servicecenter，zookeeper供中软使用，不配置的情况下默认为servicecenter
      scope: full                   #optional:scope不为full时，只允许在本app间访问，不允许跨app访问；为full就是注册时允许跨app，并且发现本租户全部微服务
      address: http://127.0.0.1:30100
      #register: manual          optional：register不配置时默认为自动注册，可选参数有自动注册auto和手动注册manual
      refeshInterval : 30s
      watch: true
`)
	d, _ := os.Getwd()
	filename1 := filepath.Join(d, "chassis.yaml")
	f1, err := os.OpenFile(filename1, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	assert.NoError(t, err)
	_, err = f1.Write(b)
	assert.NoError(t, err)
	defer f1.Close()

	b = []byte(`
---
servicecomb:
  service:
	  name: Client

`)
	d, _ = os.Getwd()
	filename2 := filepath.Join(d, "microservice.yaml")
	os.Remove(filename2)
	f2, err := os.Create(filename2)
	assert.NoError(t, err)
	defer f2.Close()
	_, err = io.WriteString(f2, string(b))
	assert.NoError(t, err)

	os.Setenv(fileutil.ChassisConfDir, d)
	defer os.Unsetenv(fileutil.ChassisConfDir)
	err = config.Init()
	assert.NoError(t, err)

	check := config.GetContractDiscoveryType()
	assert.Equal(t, "servicecenter", check)

	check = config.GetContractDiscoveryAddress()
	assert.Equal(t, "http://127.0.0.1:30100", check)

	check = config.GetContractDiscoveryAPIVersion()
	assert.Equal(t, "", check)

	dis := config.GetContractDiscoveryDisable()
	assert.Equal(t, false, dis)
	assert.NoError(t, err)
	t.Run("TestCDInit2", func(t *testing.T) {
		b := []byte(`
servicecomb:
  registry:
    type: servicecenter           #optional:可选zookeeper/servicecenter，zookeeper供中软使用，不配置的情况下默认为servicecenter
    scope: full                   #optional:scope不为full时，只允许在本app间访问，不允许跨app访问；为full就是注册时允许跨app，并且发现本租户全部微服务
    address: http://10.0.0.1:30100
    refreshInterval: 30s
    api:
      version: v1
    disabled: true
`)
		d, _ := os.Getwd()
		filename1 := filepath.Join(d, "chassis.yaml")
		f1, err := os.OpenFile(filename1,
			os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
		assert.NoError(t, err)
		_, err = f1.Write(b)
		assert.NoError(t, err)
		defer f1.Close()

		os.Setenv(fileutil.ChassisConfDir, d)
		defer os.Unsetenv(fileutil.ChassisConfDir)
		time.Sleep(1 * time.Second)
		config.ReadGlobalConfigFromArchaius()
		check := config.GetContractDiscoveryType()
		assert.Equal(t, "servicecenter", check)

		check = config.GetContractDiscoveryAddress()
		assert.Equal(t, "http://10.0.0.1:30100", check)

		check = config.GetContractDiscoveryAPIVersion()
		assert.Equal(t, "v1", check)

		dis := config.GetContractDiscoveryDisable()
		assert.Equal(t, true, dis)

	})
}
