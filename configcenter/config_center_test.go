package configcenter_test

import (
	_ "github.com/go-chassis/go-chassis/initiator"

	"github.com/go-chassis/go-chassis/configcenter"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfigCenterEndpoint(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{
		Cse: model.CseStruct{
			Config: model.Config{
				Client: model.ConfigClient{},
			},
		},
	}
	_, err := configcenter.GetConfigCenterEndpoint()
	assert.Error(t, err)
}
