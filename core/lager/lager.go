package lager

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chassis/openlog"
	"github.com/go-chassis/seclog/third_party/forked/cloudfoundry/lager"
)

// constant values for logrotate parameters
const (
	LogRotateDate     = 1
	LogRotateSize     = 10
	LogBackupCount    = 7
	RollingPolicySize = "size"
)

// log level
const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
	LevelFatal = "FATAL"
)

// output type
const (
	Stdout = "stdout"
	Stderr = "stderr"
	File   = "file"
)

//Logger is the global variable for the object of lager.Logger
//Deprecated. plz use openlog instead
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

	AccessLogFile string `yaml:"access_log_file"`
}

// Init Build constructs a *Lager.Logger with the configured parameters.
func Init(option *Options) {
	var err error
	Logger, err = NewLog(option)
	if err != nil {
		panic(err)
	}
	openlog.SetLogger(Logger)
	openlog.Debug("logger init success")
}

func toLogLevel(option string) (lager.LogLevel, error) {
	logLevel := lager.DEBUG
	switch option {
	case LevelDebug:
	case LevelInfo:
		logLevel = lager.INFO
	case LevelWarn:
		logLevel = lager.WARN
	case LevelError:
		logLevel = lager.ERROR
	case LevelFatal:
		logLevel = lager.FATAL
	default:
		return 0, errors.New("invalid log level, valid: DEBUG, INFO, WARN, ERROR, FATAL")
	}

	return logLevel, nil
}

func toFile(writer string) (*os.File, error) {
	switch writer {
	case Stdout:
		return os.Stdout, nil
	case Stderr:
		return os.Stderr, nil
	case File:
		return os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	}
	return os.Stdout, nil
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

	logger := lager.NewLoggerExt(logFilePath, option.LogFormatText)
	option.LoggerFile = logFilePath

	logLevel, err := toLogLevel(option.LoggerLevel)
	if err != nil {
		return nil, err
	}

	for _, writer := range writers {
		f, err := toFile(writer)
		if err != nil {
			return nil, err
		}
		sink := lager.NewReconfigurableSink(lager.NewWriterSink(writer, f, lager.DEBUG), logLevel)
		logger.RegisterSink(sink)
	}

	Rotators.Rotate(NewRotateConfig(option))
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
func createLogFile(localPath, out string) error {
	_, err := os.Stat(strings.Replace(filepath.Dir(filepath.Join(localPath, out)), "\\", "/", -1))
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(strings.Replace(filepath.Dir(filepath.Join(localPath, out)), "\\", "/", -1), os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	f, err := os.OpenFile(strings.Replace(filepath.Join(localPath, out), "\\", "/", -1), os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	return f.Close()
}
