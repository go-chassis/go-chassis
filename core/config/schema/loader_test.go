package schema_test

import (
	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/config/schema"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	swagger "github.com/go-chassis/go-restful-swagger20"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSchema(t *testing.T) {
	/*
		目录结构
		conf/
		├── chassis.yaml
		├── microservice1
		│   └── schema
		│       ├── schema1_1.yaml
		│       └── schema1_2yaml
	*/
	microserviceName1 := "microservice1"
	schemaID1_1 := "schema1"
	schemaID1_2 := "schema2"

	//Fix the root directory otherwise the Schema dir will be created inside /tmp/go-buildXXX///
	os.Setenv("CHASSIS_HOME", os.Getenv("GOPATH"))

	schemaDirOfMs1 := fileutil.SchemaDir(microserviceName1)

	// 创建目录
	os.RemoveAll(schemaDirOfMs1)
	err := os.MkdirAll(schemaDirOfMs1, os.ModePerm)
	assert.Nil(t, err)

	// 创建schema文件
	schemaFiles := []string{
		filepath.Join(schemaDirOfMs1, schemaID1_1+".yaml"),
		filepath.Join(schemaDirOfMs1, schemaID1_2+".yml"),
	}

	for _, schemaFile := range schemaFiles {
		file, err := os.OpenFile(schemaFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		assert.NoError(t, err)
		_, err = file.Write([]byte("test"))
		assert.NoError(t, err)
		file.Close()
	}

	err = schema.LoadSchema(fileutil.GetConfDir())
	assert.Nil(t, err)

	schemaIDs, err := schema.GetSchemaIDs(microserviceName1)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(schemaIDs))
	assert.Equal(t, schemaID1_1, schemaIDs[0], schemaID1_2, schemaIDs[1])

	assert.Equal(t, "test", schema.GetContent(schemaID1_1))

	err = os.RemoveAll(fileutil.GetConfDir())
	assert.Nil(t, err)
}

func TestSetSchemaIDs(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	config := swagger.Config{
		WebServices: restful.DefaultContainer.RegisteredWebServices(),
		OpenService: true,
		SwaggerPath: "/apidocs/",
		OutFilePath: filepath.Join(os.Getenv("CHASSIS_HOME"), "api.yaml")}
	config.Info.Description = "This is a sample server Book server"
	config.Info.Title = "swagger Book"
	sws := swagger.RegisterSwaggerService(config, restful.DefaultContainer)
	err := schema.SetSchemaInfo(sws)
	assert.NoError(t, err)
	s, e := schema.GetSchemaIDs("aaa")
	assert.Error(t, e)
	assert.Equal(t, 0, len(s))
}
