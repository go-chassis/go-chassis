package lager

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	paaslager "github.com/go-chassis/paas-lager"
	"github.com/go-chassis/paas-lager/third_party/forked/cloudfoundry/lager"
	"github.com/go-mesh/openlogging"
)

// constant values for logrotate parameters
const (
	LogRotateDate     = 1
	LogRotateSize     = 10
	LogBackupCount    = 7
	RollingPolicySize = "size"
)

//Logger is the global variable for the object of lager.Logger
//Deprecated. plz use openlogging instead
var Logger lager.Logger

// logFilePath log file path
var logFilePath string

//Options is the struct for lager information(lager.yaml)
type Options struct {
	Writers        string `yaml:"writers"`
	LoggerLevel    string `yaml:"logger_level"`
	LoggerFile     string `yaml:"logger_file"`
	LogFormatText  bool   `yaml:"log_format_text"`
	RollingPolicy  string `yaml:"rollingPolicy"`
	LogRotateDate  int    `yaml:"log_rotate_date"`
	LogRotateSize  int    `yaml:"log_rotate_size"`
	LogBackupCount int    `yaml:"log_backup_count"`
}

// Init Build constructs a *Lager.Logger with the configured parameters.
func Init(option *Options) {
	Logger = newLog(option)
	initLogRotate(logFilePath, option)
	openlogging.SetLogger(Logger)
	openlogging.Debug("logger init success")
	return
}

// newLog new log
func newLog(option *Options) lager.Logger {
	checkPassLagerDefinition(option)

	if filepath.IsAbs(option.LoggerFile) {
		createLogFile("", option.LoggerFile)
		logFilePath = filepath.Join("", option.LoggerFile)
	} else {
		createLogFile(os.Getenv("CHASSIS_HOME"), option.LoggerFile)
		logFilePath = filepath.Join(os.Getenv("CHASSIS_HOME"), option.LoggerFile)
	}
	writers := strings.Split(strings.TrimSpace(option.Writers), ",")
	if len(strings.TrimSpace(option.Writers)) == 0 {
		writers = []string{"stdout"}
	}
	paaslager.Init(paaslager.Config{
		Writers:       writers,
		LoggerLevel:   option.LoggerLevel,
		LoggerFile:    logFilePath,
		LogFormatText: option.LogFormatText,
	})

	logger := paaslager.NewLogger(option.LoggerFile)
	return logger
}

// checkPassLagerDefinition check pass lager definition
func checkPassLagerDefinition(option *Options) {
	if option.LoggerLevel == "" {
		option.LoggerLevel = "DEBUG"
	}

	if option.LoggerFile == "" {
		option.LoggerFile = "log/chassis.log"
	}

	if option.RollingPolicy == "" {
		log.Println("RollingPolicy is empty, use default policy[size]")
		option.RollingPolicy = RollingPolicySize
	} else if option.RollingPolicy != "daily" && option.RollingPolicy != RollingPolicySize {
		log.Printf("RollingPolicy is error, RollingPolicy=%s, use default policy[size].", option.RollingPolicy)
		option.RollingPolicy = RollingPolicySize
	}

	if option.LogRotateDate <= 0 || option.LogRotateDate > 10 {
		option.LogRotateDate = LogRotateDate
	}

	if option.LogRotateSize <= 0 || option.LogRotateSize > 50 {
		option.LogRotateSize = LogRotateSize
	}

	if option.LogBackupCount < 0 || option.LogBackupCount > 100 {
		option.LogBackupCount = LogBackupCount
	}
}

// createLogFile create log file
func createLogFile(localPath, out string) {
	_, err := os.Stat(strings.Replace(filepath.Dir(filepath.Join(localPath, out)), "\\", "/", -1))
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(strings.Replace(filepath.Dir(filepath.Join(localPath, out)), "\\", "/", -1), os.ModePerm)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(strings.Replace(filepath.Join(localPath, out), "\\", "/", -1), os.O_CREATE, 0640)
	if err != nil {
		panic(err)
	}
	defer f.Close()
}
