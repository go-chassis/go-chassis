package configserver_test

import (
	"testing"

	"github.com/go-chassis/go-chassis/configserver"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	_ "github.com/go-chassis/go-chassis/initiator"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigServerEndpoint(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{
		Cse: model.CseStruct{
			Config: model.Config{
				Client: model.ConfigClient{},
			},
		},
	}
	_, err := configserver.GetConfigServerEndpoint()
	assert.Error(t, err)
}
