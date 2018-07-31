package config

import (
	tcm "github.com/go-chassis/go-chassis/benchmark/helpers/config/model"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var TestDefinition *tcm.TestCfg

func GetTestDefinition() (*tcm.TestCfg, error) {
	err := readTestConfigFile()
	if err != nil {
		return nil, err
	}

	return TestDefinition, nil
}

func readTestConfigFile() error {
	defPath := filepath.Join(fileutil.GetConfDir(), "test.yaml")
	dat, err := ioutil.ReadFile(defPath)
	if err != nil {
		TestDefinition = nil
		return err
	}
	testDef := tcm.TestCfg{}
	err = yaml.Unmarshal([]byte(dat), &testDef)
	if err != nil {
		return err
	}
	TestDefinition = &testDef
	return nil
}
