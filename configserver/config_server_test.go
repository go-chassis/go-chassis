package configserver_test

import (
	"testing"

	"github.com/go-chassis/go-chassis/v2/configserver"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	_ "github.com/go-chassis/go-chassis/v2/core/registry/servicecenter"
	_ "github.com/go-chassis/go-chassis/v2/initiator"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigServerEndpoint(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{
		ServiceComb: model.ServiceComb{
			Config: model.Config{
				Client: model.ConfigClient{},
			},
		},
	}
	_, err := configserver.GetConfigServerEndpoint()
	assert.Error(t, err)
}
