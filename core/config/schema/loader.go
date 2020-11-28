package schema

import (
	"errors"
	"fmt"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/go-chassis/v2/pkg/util/fileutil"
	swagger "github.com/go-chassis/go-restful-swagger20"
	"github.com/go-chassis/openlog"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// MicroserviceMeta is the struct for micro service meta
type MicroserviceMeta struct {
	MicroserviceName string
	SchemaIDs        []string
}

// NewMicroserviceMeta gives the object of MicroserviceMeta
func NewMicroserviceMeta(microserviceName string) *MicroserviceMeta {
	return &MicroserviceMeta{
		MicroserviceName: microserviceName,
		SchemaIDs:        make([]string, 0),
	}
}

// defaultMicroserviceMetaMgr default micro-service meta-manager
var defaultMicroserviceMetaMgr map[string]*MicroserviceMeta

// schemaIDsMap default schema schema IDs map
var schemaIDsMap map[string]string

// defaultMicroServiceNames default micro-service names
var defaultMicroServiceNames = make([]string, 0)

//GetSchemaPath calculate the schema root path and return
func GetSchemaPath(name string) string {
	schemaEnv := os.Getenv(common.EnvSchemaRoot)
	var p string
	if schemaEnv != "" {
		p = filepath.Join(schemaEnv, name, fileutil.SchemaDirectory)
	} else {
		p = fileutil.SchemaDir(name)
	}
	return p
}

// LoadSchema to load the schema files and micro-service information under the conf directory
//path is the conf path
func LoadSchema(path string) error {
	/*
		conf/
		├── chassis.yaml
		├── microservice1
		│   └── schema
		│       ├── schema1.yaml
	*/
	schemaNames, err := getSchemaNames(path)
	if err != nil {
		return err
	}

	for _, msName := range schemaNames {
		var (
			microsvcMeta *MicroserviceMeta
			schemaError  error
		)
		p := GetSchemaPath(msName)
		microsvcMeta, schemaError = loadSchemaFileContent(GetSchemaPath(msName))

		if schemaError != nil {
			return schemaError
		}

		defaultMicroserviceMetaMgr[msName] = microsvcMeta
		openlog.Info(fmt.Sprintf("found schema files in %s %s", p, microsvcMeta))
	}
	return nil
}

// getSchemaNames 目录名为服务名
func getSchemaNames(confDir string) ([]string, error) {
	schemaNames := make([]string, 0)
	// 遍历confDir下的microservice文件夹
	err := filepath.Walk(confDir,
		func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return err
			}
			// 仅读取负一级目录
			if !info.IsDir() || filepath.Dir(path) != confDir {
				return nil
			}
			schemaNames = append(schemaNames, info.Name())
			return nil
		})
	return schemaNames, err
}

// SetMicroServiceNames set micro service names
func SetMicroServiceNames(confDir string) error {
	fileFormatName := `microservice(\.yaml|\.yml)$`

	err := filepath.Walk(confDir,
		func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return err
			}
			// 仅读取负一级目录
			if !info.IsDir() || filepath.Dir(path) != confDir {
				return nil
			}

			filesExist, err := getFiles(filepath.Join(confDir, info.Name()))
			if err != nil {
				return err
			}

			for _, name := range filesExist {
				ret, _ := regexp.MatchString(fileFormatName, name)
				if ret {
					defaultMicroServiceNames = append(defaultMicroServiceNames, info.Name())
				}
			}

			return nil
		})
	return err
}

// loadSchemaFileContent load scheme file content
func loadSchemaFileContent(schemaPath string) (*MicroserviceMeta, error) {
	microserviceMeta := NewMicroserviceMeta(filepath.Base(schemaPath))
	schemaFiles, err := getFiles(schemaPath)
	if err != nil {
		return microserviceMeta, err
	}

	for _, fullPath := range schemaFiles {
		schemaFile := filepath.Base(fullPath)
		dat, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return nil, errors.New("cannot find the schema file")
		}

		schemaID := strings.TrimSuffix(schemaFile, filepath.Ext(schemaFile))
		microserviceMeta.SchemaIDs = append(microserviceMeta.SchemaIDs, schemaID)
		schemaIDsMap[schemaID] = string(dat)
	}

	return microserviceMeta, nil
}

//GetContent get schema content by id
func GetContent(schemaID string) string {
	return schemaIDsMap[schemaID]
}

// getFiles get files
func getFiles(fPath string) ([]string, error) {
	files := make([]string, 0)
	_, err := os.Stat(fPath)
	if os.IsNotExist(err) {
		return files, nil
	}
	// schema文件名规则
	pat := `^.+(\.yaml|\.yml)$`
	// 遍历schemaPath下的schema文件
	err = filepath.Walk(fPath,
		func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return err
			}
			// 仅读取负一级文件
			if info.IsDir() || filepath.Dir(path) != fPath {
				return nil
			}
			ret, _ := regexp.MatchString(pat, info.Name())
			if !ret {
				return nil
			}
			files = append(files, path)
			return nil
		})
	return files, err
}

// GetMicroserviceNames get micro-service names
func GetMicroserviceNames() []string {
	return defaultMicroServiceNames
}

// GetSchemaIDs get schema IDs
func GetSchemaIDs(microserviceName string) ([]string, error) {
	microsvcMeta, ok := defaultMicroserviceMetaMgr[microserviceName]
	if !ok {
		return nil, fmt.Errorf("microservice %s not found", microserviceName)
	}
	schemaIDs := make([]string, 0)
	schemaIDs = append(schemaIDs, microsvcMeta.SchemaIDs...)
	return schemaIDs, nil
}

// init is for to initialize the defaultMicroserviceMetaMgr, and schemaIDsMap
func init() {
	defaultMicroserviceMetaMgr = make(map[string]*MicroserviceMeta)
	schemaIDsMap = make(map[string]string)
}

// SetSchemaInfo is for fill defaultMicroserviceMetaMgr and schemaIDsMap
func SetSchemaInfo(sws *swagger.SwaggerService) error {
	schemaInfoList, err := sws.GetSchemaInfoList()
	if err != nil {
		openlog.Error("get schema Info err: " + err.Error())
		return err
	}
	microsvcMeta := NewMicroserviceMeta(fileutil.SchemaDirectory)
	microsvcMeta.SchemaIDs = append(microsvcMeta.SchemaIDs, runtime.ServiceName)
	defaultMicroserviceMetaMgr[runtime.ServiceName] = microsvcMeta
	for _, schemaInfo := range schemaInfoList {
		schemaIDsMap[runtime.ServiceName] = schemaInfo
	}
	return nil
}

// SetSchemaInfoByMap is for fill defaultMicroserviceMetaMgr and schemaIDsMap
func SetSchemaInfoByMap(schemaMap map[string]string) error {
	if len(schemaMap) == 0 {
		return nil
	}
	microsvcMeta := NewMicroserviceMeta(runtime.ServiceName)
	for id, schemaInfo := range schemaMap {
		microsvcMeta.SchemaIDs = append(microsvcMeta.SchemaIDs, id)
		schemaIDsMap[id] = schemaInfo
	}

	// already read from conf/ServiceName dir
	if _, ok := defaultMicroserviceMetaMgr[runtime.ServiceName]; !ok {
		defaultMicroserviceMetaMgr[runtime.ServiceName] = microsvcMeta
	}

	return nil
}
