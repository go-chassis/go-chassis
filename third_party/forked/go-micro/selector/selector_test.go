package selector_test

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestSelector(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "client"))
	t.Log(os.Getenv("CHASSIS_HOME"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()

	LBstr := make(map[string]string)

	LBstr["name"] = "RoundRobin"
	config.GlobalDefinition.Cse.Loadbalance.Strategy = LBstr
	loadbalance.Enable()
	assert.Equal(t, "RoundRobin", config.GlobalDefinition.Cse.Loadbalance.Strategy["name"])

	LBstr["name"] = ""
	config.GlobalDefinition.Cse.Loadbalance.Strategy = LBstr
	loadbalance.Enable()
	assert.Equal(t, "", config.GlobalDefinition.Cse.Loadbalance.Strategy["name"])

	LBstr["name"] = "ABC"
	config.GlobalDefinition.Cse.Loadbalance.Strategy = LBstr
	loadbalance.Enable()
	assert.Equal(t, "ABC", config.GlobalDefinition.Cse.Loadbalance.Strategy["name"])

}
