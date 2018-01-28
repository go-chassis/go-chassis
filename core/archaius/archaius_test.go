package archaius_test

import (
	"testing"

	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/config/schema"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
	"github.com/ServiceComb/go-chassis/util/fileutil"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"time"
)

type EListener struct{}

type ConfigStruct struct {
	Yamltest1 int `yaml:"yamltest1"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestInit2(t *testing.T) {
	file := []byte(`
cse:
  isolation:
    Consumer:
      timeout:
        enabled: true
      timeoutInMilliseconds: 10
      maxConcurrentRequests: 100
      Server:
        timeoutInMilliseconds: 1000
        maxConcurrentRequests: 100
    Provider:
      Server:
        timeoutInMilliseconds: 10
        maxConcurrentRequests: 100
  loadbalance:
    serverListFilters: zoneaware
  circuitBreaker:
    Consumer:
      enabled: true
      forceOpen: false
      forceClose: true
      sleepWindowInMilliseconds: 10000
      requestVolumeThreshold: 20
      errorThresholdPercentage: 50
      Server:
        enabled: true
        forceOpen: false
        forceClose: true
        sleepWindowInMilliseconds: 10000
        requestVolumeThreshold: 20
        errorThresholdPercentage: 50
    Provider:
      Server:
        enabled: true
        forceOpen: false
        forceClose: true
        sleepWindowInMilliseconds: 10000
        requestVolumeThreshold: 20
        errorThresholdPercentage: 50
  #容错处理函数，目前暂时按照开源的方式来不进行区分处理，统一调用fallback函数
  fallback:
    Consumer:
      enabled: false
      maxConcurrentRequests: 20
  fallbackpolicy:
    Consumer:
      policy: throwexception
`)
	root, _ := fileutil.GetWorkDir()
	os.Setenv("CHASSIS_HOME", root)
	t.Log(os.Getenv("CHASSIS_HOME"))
	t.Log("Test archaius.go")
	chassisyamlContent := "APPLICATION_ID: CSE\n  \ncse:\n  service:\n    registry:\n      type: servicecenter\n  protocols:\n       highway:\n         listenAddress: 127.0.0.1:8080\n  \nssl:\n  test.consumer.certFile: test.cer\n  test.consumer.keyFile: test.key\n"
	lageryamlcontent := "logger_level: DEBUG\n \nlogger_file: log.log\n \nlog_format_text: false\n \nenable_rsyslog: true\n \nrsyslog_network:\n \nrsyslog_addr: 127.0.0.1:5151\n"
	testfilecontent := "NAME1: test1\n \nNAME2: test2"

	confdir := filepath.Join(root, "conf")
	filename1 := filepath.Join(root, "conf", "chassis.yaml")
	filename2 := filepath.Join(root, "conf", "circuit_breaker.yaml")
	filename3 := filepath.Join(root, "conf", "lager.yaml")
	filename4 := filepath.Join(root, "conf", "test_addfile.yaml")
	filename5 := filepath.Join(root, "conf", "chassis.yaml")
	filename6 := filepath.Join(root, "conf", "microservice.yaml")

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	os.Remove(filename1)
	os.Remove(filename2)
	os.Remove(filename3)
	os.Remove(filename4)
	os.Remove(filename5)
	os.Remove(filename6)
	os.Remove(confdir)
	err := os.Mkdir(confdir, 0777)
	check(err)
	defer os.Remove(confdir)
	f1, err1 := os.Create(filename1)
	check(err1)
	defer f1.Close()
	defer os.Remove(filename1)
	f2, err2 := os.Create(filename2)
	check(err2)
	defer f2.Close()
	defer os.Remove(filename2)
	f3, err3 := os.Create(filename3)
	check(err3)
	defer f3.Close()
	defer os.Remove(filename3)
	f4, err4 := os.Create(filename4)
	check(err4)
	defer f4.Close()
	defer os.Remove(filename4)
	f5, err5 := os.Create(filename5)
	check(err5)
	defer f5.Close()
	defer os.Remove(filename5)
	f6, err6 := os.Create(filename6)
	check(err6)
	defer f6.Close()
	defer os.Remove(filename6)

	_, err1 = io.WriteString(f1, chassisyamlContent)
	_, err1 = io.WriteString(f2, string(file))
	_, err1 = io.WriteString(f3, lageryamlcontent)
	_, err1 = io.WriteString(f4, testfilecontent)

	t.Log(os.Getenv("CHASSIS_HOME"))

	err = schema.LoadSchema(fileutil.GetConfDir(), false)
	if err != nil {
		t.Error(err)
	}
	archaius.Init()

	time.Sleep(10 * time.Millisecond)
	eventHandler := EListener{}
	err = archaius.RegisterListener(eventHandler, "a*")
	if err != nil {
		t.Error(err)
	}
	defer archaius.UnRegisterListener(eventHandler, "a*")

	time.Sleep(10 * time.Millisecond)
	data := archaius.GetStringByDI("darklaunch@default#0.0.1", "hi", "")
	assert.Equal(t, data, "")

	dataForDI, _ := archaius.AddDI("darklaunch@default#0.0.1")
	assert.NotEqual(t, dataForDI, nil)

	configsForDI := archaius.GetConfigsByDI("darklaunch@default#0.0.1")
	assert.NotEqual(t, configsForDI, nil)

	Stage := archaius.GetString(common.Env, "test")
	if Stage != "test" {
		t.Error("set stage is failed")
	}
	chassishome := archaius.Get("CHASSIS_HOME")
	if chassishome != root {
		t.Error("Get config by key is failed")
	}
	keyexist := archaius.Exist("CHASSIS_HOME")
	if keyexist != true {
		t.Error("Getting exist key status is failed")
	}
	ciper := archaius.GetString(common.SslCipherPluginKey, "cipertest")
	if ciper != "cipertest" {
		t.Error("Getting the string of  string cipherPlugin is failed;")
	}
	fs := archaius.GetServerListFilters()
	assert.Contains(t, fs, selector.ZoneAware)
	assert.Equal(t, 20, archaius.GetInt("cse.circuitBreaker.Consumer.requestVolumeThreshold", 0))
	assert.Equal(t, "throwexception", archaius.GetString("cse.fallbackpolicy.Consumer.policy", ""))
	assert.Equal(t, 50, archaius.GetInt("cse.circuitBreaker.Consumer.Server.errorThresholdPercentage", 0))
	assert.Equal(t, true, archaius.GetBool("cse.isolation.Consumer.timeout.enabled", false))

	yamltestkey := model.HystrixConfigWrapper{}
	archaius.UnmarshalConfig(&yamltestkey)
	assert.Equal(t, 20, yamltestkey.HystrixConfig.FallbackProperties.Consumer.MaxConcurrentRequests)

	err = archaius.AddFile(filename4)
	if err != nil {
		t.Error("Failed to add new file into the archaius")
	}

	time.Sleep(1 * time.Second)
	name1 := archaius.Get("NAME1")
	name2 := archaius.Get("NAME2")

	if name1 != "test1" && name2 != "test2" {
		t.Error("Failed to get the added file configuration key values")
	}

	err = archaius.AddKeyValue("externalSourcetest", "testextsource1")
	if err != nil {
		t.Error("Failed to Add Key and value in Externalconfig source")
	}

	configvalue := archaius.Get("externalSourcetest")
	if configvalue != "testextsource1" {
		t.Error("externalconfigsource key value is mismatched")
	}

	if archaius.Exist("externalSourcetest") != true || archaius.Exist("notexistingkey") != false {
		t.Error("Failed to get the exist status of the keys")
	}

	archaius.AddKeyValue("boolkey", "true")
	time.Sleep(10 * time.Millisecond)
	configvalue2 := archaius.GetBool("boolkey", true)
	if configvalue2 != true {
		t.Error("failed to get the value in bool")
	}
	configvalue2 = archaius.GetBool("boolkey", false)
	if configvalue2 != true {
		t.Error("failed to get the value in bool")
	}
	configvalue2 = archaius.GetBool("notexistingkey", false)
	if configvalue2 != false {
		t.Error("failed to get the value in bool")
	}

	archaius.AddKeyValue("intkey", 12)
	time.Sleep(10 * time.Millisecond)
	configvalue3 := archaius.GetInt("intkey", 12)
	if configvalue3 != 12 {
		t.Error("failed to get the value in int")
	}
	archaius.AddKeyValue("intkey", "12")
	time.Sleep(10 * time.Millisecond)
	configvalue3 = archaius.GetInt("intkey", 0)
	if configvalue3 != 12 {
		t.Error("failed to get the value in int")
	}
	configvalue3 = archaius.GetInt("notexistingkey", 0)
	if configvalue3 != 0 {
		t.Error("failed to get the value in int")
	}

	archaius.AddKeyValue("stringkey", "hello")
	time.Sleep(10 * time.Millisecond)
	configvalue = archaius.GetString("stringkey", "hello")
	if configvalue != "hello" {
		t.Error("failed to get the value in string")
	}
	configvalue = archaius.GetString("stringkey", "")
	if configvalue != "hello" {
		t.Error("failed to get the value in string")
	}
	configvalue = archaius.GetString("notexistingkey", "")
	if configvalue != "" {
		t.Error("failed to get the value in string")
	}

	archaius.AddKeyValue("floatkey", 10.12)
	time.Sleep(10 * time.Millisecond)
	configvalue4 := archaius.GetFloat64("floatkey", 0.0)
	if configvalue4 != 10.12 {
		t.Error("failed to get the value in float64")
	}
	configvalue4 = archaius.GetFloat64("notexistingkey", 0.0)
	if configvalue4 != 0.0 {
		t.Error("failed to get the value in float64")
	}

}

func (e EListener) Event(event *core.Event) {
	lager.Logger.Infof("config value after change ", event.Key, " | ", event.Value)
}
