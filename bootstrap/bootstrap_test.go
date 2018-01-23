package bootstrap_test

import (
	"github.com/ServiceComb/go-chassis/bootstrap"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var success map[string]bool

type bootstrapPlugin struct {
	Name string
}

func initialize() {
	os.Setenv("CHASSIS_HOME", "/tmp/")
	chassisConf := filepath.Join("/tmp/", "conf")
	os.MkdirAll(chassisConf, 0600)
	os.Create(filepath.Join(chassisConf, "chassis.yaml"))
	os.Create(filepath.Join(chassisConf, "microservice.yaml"))
}

func (b *bootstrapPlugin) Init() error {
	success[b.Name] = true
	return nil
}

func TestBootstrap(t *testing.T) {
	initialize()
	config.Init()
	archaius.Init()

	t.Log("Test bootstrap.go")
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	success = make(map[string]bool)

	plugin1 := &bootstrapPlugin{Name: "plugin1"}
	plugin2 := &bootstrapPlugin{Name: "plugin2"}
	plugin3 := bootstrap.BootstrapFunc(func() error {
		success["plugin3"] = true
		return nil
	})

	t.Log("Install Plugins")
	bootstrap.InstallPlugin(plugin1.Name, plugin1)
	bootstrap.InstallPlugin(plugin2.Name, plugin2)
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.ServerURI = ""
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
