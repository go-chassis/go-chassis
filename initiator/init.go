//Package initiator init necessary module
// before every other package init functions
package initiator

import (
	"fmt"
	"github.com/go-chassis/openlog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
)

// LoggerOptions has the configuration about logging
var LoggerOptions *lager.Options

func init() {
	InitLogger()
}

// InitLogger initiate config file and openlog before other modules
func InitLogger() {
	err := ParseLoggerConfig(fileutil.LogConfigPath())
	//initialize log in any case
	if err != nil {
		lager.Init(&lager.Options{
			LoggerLevel:   lager.LevelInfo,
			RollingPolicy: "size",
			Writers:       lager.Stdout,
		})
		if os.IsNotExist(err) {
			openlog.Info(fmt.Sprintf("[%s] not exist", fileutil.LogConfigPath()))
		} else {
			log.Panicln(err)
		}
	} else {
		lager.Init(&lager.Options{
			Writers:        LoggerOptions.Writers,
			LoggerLevel:    LoggerOptions.LoggerLevel,
			RollingPolicy:  LoggerOptions.RollingPolicy,
			LoggerFile:     LoggerOptions.LoggerFile,
			LogFormatText:  LoggerOptions.LogFormatText,
			LogRotateDate:  LoggerOptions.LogRotateDate,
			LogRotateSize:  LoggerOptions.LogRotateSize,
			LogBackupCount: LoggerOptions.LogBackupCount,
		})

	}
}

// ParseLoggerConfig unmarshals the logger configuration file(lager.yaml)
func ParseLoggerConfig(file string) error {
	LoggerOptions = &lager.Options{}
	err := unmarshalYamlFile(file, LoggerOptions)
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
