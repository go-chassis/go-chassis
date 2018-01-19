package metadata_test

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/metadata"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestFramework(t *testing.T) {
	metadata.Once = &sync.Once{}
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	assert := assert.New(t)
	t.Log("Service started by SDK")
	f := metadata.NewFramework()
	assert.Equal(f.Name, metadata.SdkName)
	assert.Equal(f.Version, metadata.SdkVersion)
	assert.Equal(f.Register, metadata.SdkRegistrationComponent)
}

func TestFrameworkSetNameVersionRegister(t *testing.T) {
	assert := assert.New(t)
	t.Log("setting framework name, version and registration component by exported Method")
	f := metadata.NewFramework()
	f.SetName("MyFramework")
	f.SetVersion("0.5")
	f.SetRegister("MyRegistrationComponent")
	assert.Equal(f.Name, "MyFramework")
	assert.Equal(f.Version, "0.5")
	assert.Equal(f.Register, "MyRegistrationComponent")
}
