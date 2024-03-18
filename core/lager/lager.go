package lager

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chassis/openlog"
	"github.com/go-chassis/seclog"
	"github.com/go-chassis/seclog/third_party/forked/cloudfoundry/lager"
)

// constant values for log rotate parameters
const (
	LogRotateDate  = 1
	LogRotateSize  = 10
	LogBackupCount = 7
)

// log level
const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
)

// output type
const (
	Stdout = "stdout"
	File   = "file"
)

// logFilePath log file path
var logFilePath string

// Options is the struct for lager information(lager.yaml)
type Options struct {
	Writers       string `yaml:"logWriters"`
	LoggerLevel   string `yaml:"logLevel"`
	LoggerFile    string `yaml:"logFile"`
	LogFormatText bool   `yaml:"logFormatText"`
	LogColorMode  string `yaml:"logColorMode"`

	LogRotateDisable  bool `yaml:"logRotateDisable"`
	LogRotateCompress bool `yaml:"logRotateCompress"`
	LogRotateAge      int  `yaml:"logRotateAge"`
	LogRotateSize     int  `yaml:"logRotateSize"`
	LogBackupCount    int  `yaml:"logBackupCount"`

	AccessLogFile string `yaml:"accessLogFile"`
}

// Init Build constructs a *Lager.logger with the configured parameters.
func Init(option *Options) {
	var err error
	logger, err := NewLog(option)
	if err != nil {
		panic(err)
	}
	openlog.SetLogger(logger)
	openlog.Debug("logger init success")
}

// NewLog returns a logger
func NewLog(option *Options) (lager.Logger, error) {
	checkPassLagerDefinition(option)

	localPath := ""
	if !filepath.IsAbs(option.LoggerFile) {
		localPath = os.Getenv("CHASSIS_HOME")
	}
	err := createLogFile(localPath, option.LoggerFile)
	if err != nil {
		return nil, err
	}

	logFilePath = filepath.Join(localPath, option.LoggerFile)

	writers := strings.Split(strings.TrimSpace(option.Writers), ",")

	option.LoggerFile = logFilePath

	seclog.Init(seclog.Config{
		LoggerLevel:   option.LoggerLevel,
		LogFormatText: option.LogFormatText,
		LogColorMode:  option.LogColorMode,
		Writers:       writers,
		LoggerFile:    logFilePath,
		RotateDisable: option.LogRotateDisable,
		MaxSize:       option.LogRotateSize,
		MaxAge:        option.LogRotateAge,
		MaxBackups:    option.LogBackupCount,
		Compress:      option.LogRotateCompress,
	})
	logger := seclog.NewLogger("ut")
	return logger, nil
}

// checkPassLagerDefinition check pass lager definition
func checkPassLagerDefinition(option *Options) {
	if option.LoggerLevel == "" {
		option.LoggerLevel = "DEBUG"
	}

	if option.LoggerFile == "" {
		option.LoggerFile = "log/chassis.log"
	}

	if option.LogRotateAge < 0 || option.LogRotateAge > 10 {
		option.LogRotateAge = LogRotateDate
	}

	if option.LogRotateSize <= 0 || option.LogRotateSize > 500 {
		option.LogRotateSize = LogRotateSize
	}

	if option.LogBackupCount < 0 || option.LogBackupCount > 100 {
		option.LogBackupCount = LogBackupCount
	}
	if option.Writers == "" {
		option.Writers = "file,stdout"
	}
	if option.LogColorMode == "" {
		option.LogColorMode = lager.ColorModeAuto
	}
}

// createLogFile create log file
func createLogFile(localPath, out string) error {
	_, err := os.Stat(strings.Replace(filepath.Dir(filepath.Join(localPath, out)), "\\", "/", -1))
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(strings.Replace(filepath.Dir(filepath.Join(localPath, out)), "\\", "/", -1), os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	f, err := os.OpenFile(strings.Replace(filepath.Join(localPath, out), "\\", "/", -1), os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	return f.Close()
}
