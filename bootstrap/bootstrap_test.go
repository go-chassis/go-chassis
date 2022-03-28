package bootstrap_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-chassis/go-chassis/v2/bootstrap"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	_ "github.com/go-chassis/go-chassis/v2/core/registry/servicecenter"
	_ "github.com/go-chassis/go-chassis/v2/initiator"
	"github.com/stretchr/testify/assert"
)

var success map[string]bool

type bootstrapPlugin struct {
	Name string
}

func initialize(t *testing.T) {
	os.Setenv("CHASSIS_HOME", "/tmp/")
	t.Cleanup(func() {
		os.Unsetenv("CHASSIS_HOME")
	})
	chassisConf := filepath.Join("/tmp/", "conf")
	os.MkdirAll(chassisConf, 0700)
	os.Create(filepath.Join(chassisConf, "chassis.yaml"))
	os.Create(filepath.Join(chassisConf, "microservice.yaml"))
}

func (b *bootstrapPlugin) Init() error {
	success[b.Name] = true
	return nil
}

func TestBootstrap(t *testing.T) {
	initialize(t)
	config.Init()
	time.Sleep(1 * time.Second)
	config.GlobalDefinition = &model.GlobalCfg{}
	config.MicroserviceDefinition = &model.ServiceSpec{}
	config.GlobalDefinition.ServiceComb.Registry.APIVersion.Version = "v2"

	t.Log("Test bootstrap.go")

	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	success = make(map[string]bool)

	plugin1 := &bootstrapPlugin{Name: "plugin1"}
	plugin2 := &bootstrapPlugin{Name: "plugin2"}
	plugin3 := bootstrap.Func(func() error {
		success["plugin3"] = true
		return nil
	})

	t.Log("Install Plugins")
	bootstrap.InstallPlugin(plugin1.Name, plugin1)
	bootstrap.InstallPlugin(plugin2.Name, plugin2)
	config.GlobalDefinition.ServiceComb.Config.Client.ServerURI = ""
	bootstrap.Bootstrap()

	t.Log("verifying Plugins")
	assert.Equal(t, 2, len(success))
	assert.True(t, success[plugin1.Name])
	assert.True(t, success[plugin2.Name])

	bootstrap.InstallPlugin("plugin3", plugin3)
	bootstrap.Bootstrap()

	assert.Equal(t, 3, len(success))
	assert.True(t, success["plugin3"])
}
