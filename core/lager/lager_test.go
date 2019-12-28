package lager_test

import (
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/env"

	//"github.com/go-chassis/go-chassis/core/config"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitialize1(t *testing.T) {
	path := os.Getenv("GOPATH")
	logDir := filepath.Join(path, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server")
	env.ChassisHome = logDir
	os.RemoveAll(logDir)

	t.Log("Initializing lager")
	t.Log("creating log/chassis.log")
	lager.Init(&lager.Options{
		LoggerFile: filepath.Join("log", "chassis.log"),
	})
	defer os.RemoveAll(logDir)

	if _, err := os.Stat(logDir); err != nil {
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
	env.ChassisHome = logDir
	os.RemoveAll(logDir)

	//initializing config for to initialize PassLagerDefinition variable
	t.Log("initializing config for to initialize PassLagerDefinition variable")

	//Initializing lager
	t.Log("Initializing lager")
	lager.Init(&lager.Options{})
	defer os.RemoveAll(logDir)

	if _, err := os.Stat(logDir); err != nil {
		if os.IsNotExist(err) {
			t.Error(err)
		}
	}

	time.Sleep(1 * time.Second)
}
