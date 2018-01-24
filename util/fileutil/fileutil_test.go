package fileutil_test

import (
	"github.com/ServiceComb/go-chassis/util/fileutil"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGetWorkDirHmNotSet(t *testing.T) {
	os.Setenv("CHASSIS_HOME", "test")
	assert.Equal(t, "conf", filepath.Base(fileutil.GetConfDir()))

}
func TestGetWorkDir(t *testing.T) {
	os.Setenv("CHASSIS_HOME", "test")
	assert.Equal(t, "conf", filepath.Base(fileutil.GetConfDir()))
}

func TestHystricDefinaiton(t *testing.T) {
	os.Setenv("CHASSIS_HOME", "test")
	def := fileutil.HystrixDefinition()
	assert.Equal(t, "test/conf/circuit_breaker.yaml", def)
}
func TestMicroserviceDefinition(t *testing.T) {
	def := fileutil.MicroserviceDefinition("micro")
	assert.Equal(t, "test/conf/micro/microservice.yaml", def)
}
func TestGlobalDefinition(t *testing.T) {
	def := fileutil.GlobalDefinition()
	assert.Equal(t, "test/conf/chassis.yaml", def)
}
func TestPassLagerDefinition(t *testing.T) {
	def := fileutil.PassLagerDefinition()
	assert.Equal(t, "test/conf/lager.yaml", def)
}
func TestSchemaDir(t *testing.T) {
	def := fileutil.SchemaDir("micro")
	assert.Equal(t, "test/conf/micro/schema", def)
}
func TestGetDefinition(t *testing.T) {
	def := fileutil.GetDefinition()
	assert.Equal(t, "test/conf/microservice.yaml", def)
}
func TestGetWorkDirConfSet(t *testing.T) {
	os.Setenv("CHASSIS_CONF_DIR", "conf")
	assert.Equal(t, "conf", filepath.Base(fileutil.GetConfDir()))
}
