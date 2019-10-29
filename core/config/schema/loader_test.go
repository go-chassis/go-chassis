package schema_test

import (
	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/config/schema"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	swagger "github.com/go-chassis/go-restful-swagger20"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"sort"
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
		├── microservice2
		│   └── schema
		│       ├── schema2_1.yaml
		│       └── schema2_2.yamll
		└── microservice3
	*/
	microserviceName1 := "microservice1"
	microserviceName2 := "microservice2"
	microserviceName3 := "microservice3"
	schemaID1_1 := "schema1_1"
	schemaID1_2 := "schema1_2"
	schemaID2_1 := "schema2_1"
	schemaID2_2 := "schema2_2"

	NoExistMicroserviceName := "NoExistMicroservice"

	//Fix the root directory otherwise the Schema dir will be created inside /tmp/go-buildXXX///
	os.Setenv("CHASSIS_HOME", os.Getenv("GOPATH"))

	schemaDirOfMs1 := fileutil.SchemaDir(microserviceName1)
	schemaDirOfMs2 := fileutil.SchemaDir(microserviceName2)
	Ms3Dir := filepath.Join(fileutil.GetConfDir(), microserviceName3)

	// 创建目录
	os.RemoveAll(schemaDirOfMs1)
	os.RemoveAll(schemaDirOfMs2)
	os.RemoveAll(Ms3Dir)
	err := os.MkdirAll(schemaDirOfMs1, os.ModePerm)
	assert.Nil(t, err)

	err = os.MkdirAll(schemaDirOfMs2, os.ModePerm)
	assert.Nil(t, err)

	err = os.MkdirAll(Ms3Dir, os.ModePerm)
	assert.Nil(t, err)

	// 创建schema文件
	schemaFiles := []string{
		filepath.Join(schemaDirOfMs1, schemaID1_1+".yaml"),
		filepath.Join(schemaDirOfMs1, schemaID1_2+".yml"),
		filepath.Join(schemaDirOfMs2, schemaID2_1+".yaml"),
		filepath.Join(schemaDirOfMs2, schemaID2_2+".yamll"),
	}

	for _, schemaFile := range schemaFiles {
		file, err := os.OpenFile(schemaFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		assert.Nil(t, err)
		file.Close()
	}

	t.Log("========加载schema")
	err = schema.LoadSchema(fileutil.GetConfDir())
	assert.Nil(t, err)

	t.Log("========查询schemaID")
	t.Log("====查询", microserviceName1)
	schemaIDs, err := schema.GetSchemaIDs(microserviceName1)
	sort.Strings(schemaIDs)
	t.Log(schemaIDs)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(schemaIDs))
	assert.Equal(t, schemaID1_1, schemaIDs[0], schemaID1_2, schemaIDs[1])

	t.Log("====查询是否为值拷贝")
	schemaIDs[0] = "test.huawei"
	schemaIDs, err = schema.GetSchemaIDs(microserviceName1)
	sort.Strings(schemaIDs)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(schemaIDs))
	assert.Equal(t, schemaID1_1, schemaIDs[0], schemaID1_2, schemaIDs[1])

	t.Log("====查询", microserviceName2)
	schemaIDs, err = schema.GetSchemaIDs(microserviceName2)
	sort.Strings(schemaIDs)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(schemaIDs))
	assert.Equal(t, schemaID2_1, schemaIDs[0])

	t.Log("====查询", microserviceName3)
	schemaIDs, err = schema.GetSchemaIDs(microserviceName3)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(schemaIDs))

	t.Log("====查询", NoExistMicroserviceName)
	schemaIDs, err = schema.GetSchemaIDs(NoExistMicroserviceName)
	assert.NotNil(t, err)

	t.Log("========查询microserviceNames")
	microserviceNames := schema.GetMicroserviceNamesBySchemas()
	sort.Strings(microserviceNames)
	assert.Equal(t, 3, len(microserviceNames))
	assert.Equal(t,
		microserviceName1, microserviceNames[0],
		microserviceName2, microserviceNames[1],
		microserviceName3, microserviceNames[2])

	err = os.RemoveAll(fileutil.GetConfDir())
	assert.Nil(t, err)
}

func TestSetSchemaIDs(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	log.Println(os.Getenv("CHASSIS_HOME"))
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
