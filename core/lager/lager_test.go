package lager

import (
	//"github.com/ServiceComb/go-chassis/core/config"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitialize1(t *testing.T) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))

	logdir := path + "/src/github.com/ServiceComb/go-chassis/examples/discovery/server/log"
	os.RemoveAll(logdir)

	t.Log("Initializing lager")
	t.Log("creating log/chassis.log")
	Initialize("", "INFO", "", "size", true, 1, 10, 7)
	defer os.RemoveAll(logdir)

	if _, err := os.Stat(logdir); err != nil {
		if os.IsNotExist(err) {
			t.Error(err)
		}
	}

	t.Log("duplicate initialization")
	Initialize("", "INFO", "", "size", true, 1, 10, 7)
}

func TestInitialize2(t *testing.T) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "ServiceComb",
		"go-chassis", "examples", "discovery", "server"))

	//initializing config for to initialize PassLagerDefinition variable
	t.Log("initializing config for to initialize PassLagerDefinition variable")
	//err := config.Init()
	//if err != nil {
	//	log.Printf("Failed to initialize conf, err=%s\n", err)
	//}

	logdir := path + "/src/github.com/ServiceComb/go-chassis/examples/discovery/server/log"
	os.RemoveAll(logdir)

	//Initializing lager
	t.Log("Initializing lager")
	Initialize("", "INFO", "", "size", true, 1, 10, 7)
	defer os.RemoveAll(logdir)

	if _, err := os.Stat(logdir); err != nil {
		if os.IsNotExist(err) {
			t.Error(err)
		}
	}

	time.Sleep(1 * time.Second)
}
