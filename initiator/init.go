//Package initiator init necessary module
// before every other package init functions
package initiator

import (
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

// PaasLagerDefinition is having the information about loging
var PaasLagerDefinition *PassLagerCfg

func init() {
	InitLogger()
}

// InitLogger initiate config file and openlogging before other modules
func InitLogger() {
	err := ParseLoggerConfig(fileutil.PaasLagerDefinition())
	//initialize log in any case
	if err != nil {
		lager.Initialize("", "", "",
			"", false, 1, 10, 7)
		if os.IsNotExist(err) {
			lager.Logger.Infof("[%s] not exist", fileutil.PaasLagerDefinition())
		} else {
			log.Panicln(err)
		}
	} else {
		lager.Initialize(PaasLagerDefinition.Writers, PaasLagerDefinition.LoggerLevel,
			PaasLagerDefinition.LoggerFile, PaasLagerDefinition.RollingPolicy,
			PaasLagerDefinition.LogFormatText, PaasLagerDefinition.LogRotateDate,
			PaasLagerDefinition.LogRotateSize, PaasLagerDefinition.LogBackupCount)
	}
}

// ParseLoggerConfig unmarshals the logger configuration file(lager.yaml)
func ParseLoggerConfig(file string) error {
	PaasLagerDefinition = &PassLagerCfg{}
	err := unmarshalYamlFile(file, PaasLagerDefinition)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return err
}

func unmarshalYamlFile(file string, target interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, target)
}

//PassLagerCfg is the struct for lager information(passlager.yaml)
type PassLagerCfg struct {
	Writers        string `yaml:"writers"`
	LoggerLevel    string `yaml:"logger_level"`
	LoggerFile     string `yaml:"logger_file"`
	LogFormatText  bool   `yaml:"log_format_text"`
	RollingPolicy  string `yaml:"rollingPolicy"`
	LogRotateDate  int    `yaml:"log_rotate_date"`
	LogRotateSize  int    `yaml:"log_rotate_size"`
	LogBackupCount int    `yaml:"log_backup_count"`
}
