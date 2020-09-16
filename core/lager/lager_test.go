package lager_test

import (
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitialize1(t *testing.T) {
	lager.Init(&lager.Options{
		LoggerFile: filepath.Join("./log", "chassis.log"),
		Writers:    "file",
	})

	if _, err := os.Stat("log"); err != nil {
		if os.IsNotExist(err) {
			t.Error(err)
		}
	}

	t.Log("duplicate initialization")
	lager.Init(&lager.Options{})
}

func TestInitialize2(t *testing.T) {
	path := os.Getenv("GOPATH")
	logDir := filepath.Join(path, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server")
	os.Setenv("CHASSIS_HOME", logDir)

	//initializing config for to initialize PassLagerDefinition variable
	t.Log("initializing config for to initialize PassLagerDefinition variable")

	//Initializing lager
	lager.Init(&lager.Options{})

	if _, err := os.Stat(logDir); err != nil {
		if os.IsNotExist(err) {
			t.Error(err)
		}
	}

	time.Sleep(1 * time.Second)
}
