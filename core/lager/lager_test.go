package lager_test

import (
	"github.com/go-chassis/go-chassis/core/lager"
	//"github.com/go-chassis/go-chassis/core/config"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitialize1(t *testing.T) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))

	logdir := path + "/src/github.com/go-chassis/go-chassis/examples/discovery/server/log"
	os.RemoveAll(logdir)

	t.Log("Initializing lager")
	t.Log("creating log/chassis.log")
	lager.Init(&lager.Options{})
	defer os.RemoveAll(logdir)

	if _, err := os.Stat(logdir); err != nil {
		if os.IsNotExist(err) {
			t.Error(err)
		}
	}

	t.Log("duplicate initialization")
	lager.Init(&lager.Options{})
}

func TestInitialize2(t *testing.T) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "go-chassis",
		"go-chassis", "examples", "discovery", "server"))

	//initializing config for to initialize PassLagerDefinition variable
	t.Log("initializing config for to initialize PassLagerDefinition variable")

	logdir := path + "/src/github.com/go-chassis/go-chassis/examples/discovery/server/log"
	os.RemoveAll(logdir)

	//Initializing lager
	t.Log("Initializing lager")
	lager.Init(&lager.Options{})
	defer os.RemoveAll(logdir)

	if _, err := os.Stat(logdir); err != nil {
		if os.IsNotExist(err) {
			t.Error(err)
		}
	}

	time.Sleep(1 * time.Second)
}
