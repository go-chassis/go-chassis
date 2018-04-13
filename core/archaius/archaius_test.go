package archaius_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/config/schema"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalancer"
	"github.com/ServiceComb/go-chassis/util/fileutil"
	"github.com/stretchr/testify/assert"
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
	cbBytes := []byte(`
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
  circuitBreaker:
    Consumer:
      enabled: true
      forceOpen: false
      forceClosed: true
      sleepWindowInMilliseconds: 10000
      requestVolumeThreshold: 20
      errorThresholdPercentage: 50
      Server:
        enabled: true
        forceOpen: false
        forceClosed: true
        sleepWindowInMilliseconds: 10000
        requestVolumeThreshold: 20
        errorThresholdPercentage: 50
    Provider:
      Server:
        enabled: true
        forceOpen: false
        forceClosed: true
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
	lbBytes := []byte(`
---
cse: 
  loadbalance: 
    TargetService: 
      backoff: 
        kind: constant
      retryEnabled: false
      strategy: 
        name: WeightedResponse
    target_Service: 
      backoff:
        maxMs: 500
        minMs: 200
        kind: constant
      retryEnabled: false
      strategy: 
        name: WeightedResponse
    backoff: 
      maxMs: 400
      minMs: 200
      kind: constant
    retryEnabled: false
    retryOnNext: 2
    retryOnSame: 3
    serverListFilters: zoneaware
    strategy: 
      name: WeightedResponse

`)
	root, _ := fileutil.GetWorkDir()
	os.Setenv("CHASSIS_HOME", root)
	t.Log(os.Getenv("CHASSIS_HOME"))
	t.Log("Test archaius.go")
	chassisyamlContent := "APPLICATION_ID: CSE\n  \ncse:\n  service:\n    registry:\n      type: servicecenter\n  protocols:\n       highway:\n         listenAddress: 127.0.0.1:8080\n  \nssl:\n  test.consumer.certFile: test.cer\n  test.consumer.keyFile: test.key\n"
	lageryamlcontent := "logger_level: DEBUG\n \nlogger_file: log.log\n \nlog_format_text: false\n \nenable_rsyslog: true\n \nrsyslog_network:\n \nrsyslog_addr: 127.0.0.1:5151\n"

	confdir := filepath.Join(root, "conf")
	filename1 := filepath.Join(root, "conf", "chassis.yaml")
	circuitBreakerFileName := filepath.Join(root, "conf", "circuit_breaker.yaml")
	filename3 := filepath.Join(root, "conf", "lager.yaml")
	lbFileName := filepath.Join(root, "conf", "load_balancing.yaml")
	filename6 := filepath.Join(root, "conf", "microservice.yaml")

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	os.Remove(filename1)
	os.Remove(circuitBreakerFileName)
	os.Remove(filename3)
	os.Remove(lbFileName)
	os.Remove(filename6)
	os.Remove(confdir)
	err := os.Mkdir(confdir, 0777)
	check(err)
	defer os.Remove(confdir)
	f1, err1 := os.Create(filename1)
	check(err1)
	defer f1.Close()
	defer os.Remove(filename1)
	circuitBreakerFile, err2 := os.Create(circuitBreakerFileName)
	check(err2)
	defer circuitBreakerFile.Close()
	defer os.Remove(circuitBreakerFileName)
	f3, err3 := os.Create(filename3)
	check(err3)
	defer f3.Close()
	defer os.Remove(filename3)
	lbFile, err4 := os.Create(lbFileName)
	check(err4)
	defer lbFile.Close()
	defer os.Remove(lbFileName)
	f6, err6 := os.Create(filename6)
	check(err6)
	defer f6.Close()
	defer os.Remove(filename6)

	_, err1 = io.WriteString(f1, chassisyamlContent)
	_, err1 = io.WriteString(circuitBreakerFile, string(cbBytes))
	_, err1 = io.WriteString(f3, lageryamlcontent)
	_, err1 = io.WriteString(lbFile, string(lbBytes))

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
	config.ReadLBFromArchaius()
	fs := config.GetServerListFilters()
	assert.Contains(t, fs, loadbalancer.ZoneAware)
	assert.Equal(t, 20, archaius.GetInt("cse.circuitBreaker.Consumer.requestVolumeThreshold", 0))
	assert.Equal(t, "throwexception", archaius.GetString("cse.fallbackpolicy.Consumer.policy", ""))
	assert.Equal(t, 50, archaius.GetInt("cse.circuitBreaker.Consumer.Server.errorThresholdPercentage", 0))
	assert.Equal(t, true, archaius.GetBool("cse.isolation.Consumer.timeout.enabled", false))
	t.Log("Unmarshall cb")
	cb := model.HystrixConfigWrapper{}
	archaius.UnmarshalConfig(&cb)
	assert.Equal(t, 20, cb.HystrixConfig.FallbackProperties.Consumer.MaxConcurrentRequests)
	assert.Equal(t, 1000, cb.HystrixConfig.IsolationProperties.Consumer.AnyService["Server"].TimeoutInMilliseconds)
	assert.NotEqual(t, 22, cb.HystrixConfig.IsolationProperties.Consumer.AnyService["Server"].TimeoutInMilliseconds)
	t.Log("Unmarshall lb")
	lbConfig := model.LBWrapper{}
	archaius.UnmarshalConfig(&lbConfig)
	t.Log(lbConfig.Prefix.LBConfig)

	assert.Equal(t, loadbalancer.ZoneAware, lbConfig.Prefix.LBConfig.Filters)
	t.Log(lbConfig.Prefix.LBConfig.AnyService)
	assert.Equal(t, "WeightedResponse", lbConfig.Prefix.LBConfig.AnyService["TargetService"].Strategy["name"])
	assert.Equal(t, "WeightedResponse", lbConfig.Prefix.LBConfig.AnyService["target_Service"].Strategy["name"])
	assert.Equal(t, 500, int(lbConfig.Prefix.LBConfig.AnyService["target_Service"].Backoff.MaxMs))
	err = archaius.AddFile(lbFileName)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	err = archaius.AddKeyValue("memorySourcetestKeyCheck", "testmemsource1")
	if err != nil {
		t.Error("Failed to Add Key and value in Memoryconfig source")
	}

	err = archaius.DeleteKeyValue("memorySourcetestKeyCheck", "testmemsource1")
	if err != nil {
		t.Error("Failed to Delete Key and value in Memoryconfig source")
	}

	err = archaius.AddKeyValue("memorySourcetest", "testmemsource1")
	if err != nil {
		t.Error("Failed to Add Key and value in memoryconfig source")
	}

	configvalue := archaius.Get("memorySourcetest")
	if configvalue != "testmemsource1" {
		t.Error("memoryconfigsource key value is mismatched")
	}

	if archaius.Exist("memorySourcetest") != true || archaius.Exist("notexistingkey") != false {
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
