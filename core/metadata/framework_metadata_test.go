package metadata_test

import (
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/metadata"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
func TestFramework(t *testing.T) {
	metadata.Once = &sync.Once{}
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
