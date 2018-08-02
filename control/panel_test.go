package control_test

import (
	"github.com/go-chassis/go-chassis/control"
	_ "github.com/go-chassis/go-chassis/control/archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstallPlugin(t *testing.T) {
	control.InstallPlugin("test", nil)

}
func TestInit(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{
		Panel: model.ControlPanel{
			Infra: "",
		},
	}
	err := control.Init()
	assert.NoError(t, err)
	config.GlobalDefinition.Panel.Infra = "xxx"
	err = control.Init()
	assert.Error(t, err)
}
