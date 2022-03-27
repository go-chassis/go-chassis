package lager_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/openlog"
)

func TestInitialize1(t *testing.T) {
	logDir := t.TempDir()
	lager.Init(&lager.Options{
		LoggerFile: filepath.Join(logDir, "chassis.log"),
		Writers:    "file",
	})

	if _, err := os.Stat(logDir); err != nil {
		if os.IsNotExist(err) {
			t.Error(err)
		}
	}

	t.Log("duplicate initialization")
	lager.Init(&lager.Options{})
}

func TestInitialize2(t *testing.T) {
	homeDir := t.TempDir()
	os.Setenv("CHASSIS_HOME", homeDir)
	logDir := filepath.Join(homeDir, "log")

	//initializing config for to initialize PassLagerDefinition variable
	t.Log("initializing config for to initialize PassLagerDefinition variable")

	//Initializing lager
	lager.Init(&lager.Options{LoggerLevel: "INFO"})
	openlog.Debug("no output")
	openlog.Info("output")
	if _, err := os.Stat(logDir); err != nil {
		if os.IsNotExist(err) {
			t.Error(err)
		}
	}

	time.Sleep(1 * time.Second)
}
