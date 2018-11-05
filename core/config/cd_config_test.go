package config_test

import (
	"os"
	"testing"

	_ "github.com/go-chassis/go-chassis/initiator"

	"github.com/go-chassis/go-chassis/core/config"

	"github.com/stretchr/testify/assert"
)

func TestCDInit(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	config.Init()
}

func TestGetContractDiscoveryType(t *testing.T) {
	check := config.GetContractDiscoveryType()
	assert.Equal(t, "servicecenter", check)
}

func TestGetContractDiscoveryAddress(t *testing.T) {
	check := config.GetContractDiscoveryAddress()
	assert.Equal(t, "", check)
}

func TestGetContractDiscoveryTenant(t *testing.T) {
	check := config.GetContractDiscoveryTenant()
	assert.Equal(t, "default", check)
}

func TestGetContractDiscoveryAPIVersion(t *testing.T) {
	check := config.GetContractDiscoveryAPIVersion()
	assert.Equal(t, "", check)
}

func TestGetContractDiscoveryDisable(t *testing.T) {
	check := config.GetContractDiscoveryDisable()
	assert.Equal(t, false, check)
}
