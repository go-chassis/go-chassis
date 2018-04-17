package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/config/schema"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/util/fileutil"

	"gopkg.in/yaml.v2"
)

// GlobalDefinition is having the information about region, load balancing, service center, config center,
// protocols, and handlers for the micro service
var GlobalDefinition *model.GlobalCfg
var lbConfig *model.LBWrapper

// MicroserviceDefinition is having the info about application id, provider info, description of the service,
// and description of the instance
var MicroserviceDefinition *model.MicroserviceCfg

// PaasLagerDefinition is having the information about loging
var PaasLagerDefinition *model.PassLagerCfg

// RouterDefinition is route rule config
var RouterDefinition *model.RouterConfig

//HystrixConfig is having info about isolation, circuit breaker, fallback properities of the micro service
var HystrixConfig *model.HystrixConfigWrapper

// NodeIP gives the information of node ip
var NodeIP string

// SelfServiceID 单进程多微服务根本没法记录依赖关系，因为一个进程里有多个微服务，你在调用别的微服务时到底该怎么添加依赖关系？
//只能随意赋值个id
var SelfServiceID string

// SelfServiceName is self micro service name
var SelfServiceName string

// SelfMetadata is gives meta data of the self micro service
var SelfMetadata map[string]string

// SelfVersion gives version of the self micro service
var SelfVersion string

// ErrNoName is used to represent the service name missing error
var ErrNoName = errors.New("Microservice name is missing in description file")

// parse unmarshal configurations on respective structure
func parse() error {
	err := readGlobalConfigFile()
	if err != nil {
		return err
	}
	err = ReadLBFromArchaius()
	if err != nil {
		return err
	}

	err = ReadHystrixFromArchaius()
	if err != nil {
		return err
	}

	err = readMicroserviceConfigFiles()
	if err != nil {
		return err
	}

	populateConfigCenterAddress()
	populateServiceRegistryAddress()
	populateMonitorServerAddress()
	populateServiceEnvironment()
	populateServiceName()
	populateVersion()
	populateTenant()

	return nil
}

// populateServiceRegistryAddress populate service registry address
func populateServiceRegistryAddress() {
	//Registry Address , higher priority for environment variable
	registryAddrFromEnv := archaius.GetString(common.CseRegistryAddress, "")
	if registryAddrFromEnv != "" {
		GlobalDefinition.Cse.Service.Registry.Registrator.Address = registryAddrFromEnv
		GlobalDefinition.Cse.Service.Registry.ServiceDiscovery.Address = registryAddrFromEnv
		GlobalDefinition.Cse.Service.Registry.ContractDiscovery.Address = registryAddrFromEnv
		GlobalDefinition.Cse.Service.Registry.Address = registryAddrFromEnv
	}
}

// populateConfigCenterAddress populate config center address
func populateConfigCenterAddress() {
	//Config Center Address , higher priority for environment variable
	configCenterAddrFromEnv := archaius.GetString(common.CseConfigCenterAddress, "")
	if configCenterAddrFromEnv != "" {
		GlobalDefinition.Cse.Config.Client.ServerURI = configCenterAddrFromEnv
	}
}

// populateMonitorServerAddress populate monitor server address
func populateMonitorServerAddress() {
	//Monitor Center Address , higher priority for environment variable
	monitorServerAddrFromEnv := archaius.GetString(common.CseMonitorServer, "")
	if monitorServerAddrFromEnv != "" {
		GlobalDefinition.Cse.Monitor.Client.ServerURI = monitorServerAddrFromEnv
	}
}

// populateServiceEnvironment populate service environment
func populateServiceEnvironment() {
	if e := archaius.GetString(common.Env, ""); e != "" {
		MicroserviceDefinition.ServiceDescription.Environment = e
	}
}

// populateServiceName populate service name
func populateServiceName() {
	if e := archaius.GetString(common.ServiceName, ""); e != "" {
		MicroserviceDefinition.ServiceDescription.Name = e
	}
}

// populateVersion populate version
func populateVersion() {
	if e := archaius.GetString(common.Version, ""); e != "" {
		MicroserviceDefinition.ServiceDescription.Version = e
	}
}

// populateTenant populate tenant
func populateTenant() {
	if GlobalDefinition.Cse.Service.Registry.Tenant == "" {
		GlobalDefinition.Cse.Service.Registry.Tenant = common.DefaultApp
	}
}

// readGlobalConfigFile for to unmarshal the global config file(chassis.yaml) information
func readGlobalConfigFile() error {
	globalDef := model.GlobalCfg{}
	err := archaius.UnmarshalConfig(&globalDef)
	if err != nil {
		return err
	}
	GlobalDefinition = &globalDef

	return nil
}

// ReadLBFromArchaius for to unmarshal the global config file(chassis.yaml) information
func ReadLBFromArchaius() error {
	lbMutex.Lock()
	defer lbMutex.Unlock()
	lbDef := model.LBWrapper{}
	err := archaius.UnmarshalConfig(&lbDef)
	if err != nil {
		return err
	}
	lbConfig = &lbDef

	return nil
}

type pathError struct {
	Path string
	Err  error
}

func (e *pathError) Error() string { return e.Path + ": " + e.Err.Error() }

// parsePaasLagerConfig unmarshals the paas lager configuration file(lager.yaml)
func parsePaasLagerConfig(file string) error {
	PaasLagerDefinition = &model.PassLagerCfg{}
	err := unmarshalYamlFile(file, PaasLagerDefinition)
	if err != nil && !os.IsNotExist(err) {
		return &pathError{Path: file, Err: err}
	}
	return err
}

// parseRouterConfig is unmarshal the paas lager configuration file(lager.yaml)
func parseRouterConfig(file string) error {
	RouterDefinition = &model.RouterConfig{}
	err := unmarshalYamlFile(file, RouterDefinition)
	if err != nil && !os.IsNotExist(err) {
		return &pathError{Path: file, Err: err}
	}
	return err
}

func unmarshalYamlFile(file string, target interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, target)
}

// ReadHystrixFromArchaius is unmarshal hystrix configuration file(circuit_breaker.yaml)
func ReadHystrixFromArchaius() error {
	cbMutex.RLock()
	defer cbMutex.RUnlock()
	hystrixCnf := model.HystrixConfigWrapper{}
	err := archaius.UnmarshalConfig(&hystrixCnf)
	if err != nil {
		return err
	}
	HystrixConfig = &hystrixCnf
	return nil
}

// readMicroserviceConfigFiles read micro service configuration file
func readMicroserviceConfigFiles() error {
	MicroserviceDefinition = &model.MicroserviceCfg{}
	//find only one microservice yaml
	microserviceNames := schema.GetMicroserviceNames()
	defPath := fileutil.GetMicroserviceDesc()
	data, err := ioutil.ReadFile(defPath)
	if err != nil {
		lager.Logger.Errorf(err, fmt.Sprintf("WARN: Missing microservice description file: %s", err.Error()))
		if len(microserviceNames) == 0 {
			return errors.New("Missing microservice description file")
		}
		msName := microserviceNames[0]
		msDefPath := fileutil.MicroserviceDefinition(msName)
		lager.Logger.Warnf(fmt.Sprintf("Try to find microservice description file in [%s]", msDefPath))
		data, err := ioutil.ReadFile(msDefPath)
		if err != nil {
			return fmt.Errorf("Missing microservice description file: %s", err.Error())
		}
		ReadMicroserviceConfigFromBytes(data)
		return nil
	}
	return ReadMicroserviceConfigFromBytes(data)
}

// ReadMicroserviceConfigFromBytes read micro service configurations from bytes
func ReadMicroserviceConfigFromBytes(data []byte) error {
	microserviceDef := model.MicroserviceCfg{}
	err := yaml.Unmarshal([]byte(data), &microserviceDef)
	if err != nil {
		return err
	}
	if microserviceDef.ServiceDescription.Name == "" {
		return ErrNoName
	}
	if microserviceDef.ServiceDescription.Version == "" {
		microserviceDef.ServiceDescription.Version = common.DefaultVersion
	}

	MicroserviceDefinition = &microserviceDef
	return nil
}

//GetLoadBalancing return lb config
func GetLoadBalancing() *model.LoadBalancing {
	if lbConfig != nil {
		return lbConfig.Prefix.LBConfig
	}
	return nil
}

//GetHystrixConfig return cb config
func GetHystrixConfig() *model.HystrixConfig {
	return HystrixConfig.HystrixConfig
}

// Init is initialize the configuration directory, lager, archaius, route rule, and schema
func Init() error {
	err := parsePaasLagerConfig(fileutil.PaasLagerDefinition())
	//initialize log in any case
	if err != nil {
		lager.Initialize("", "", "", "", true, 1, 10, 7)
		if os.IsNotExist(err) {
			lager.Logger.Infof("[%s] not exist", fileutil.PaasLagerDefinition())
		} else {
			return err
		}
	} else {
		lager.Initialize(PaasLagerDefinition.Writers, PaasLagerDefinition.LoggerLevel,
			PaasLagerDefinition.LoggerFile, PaasLagerDefinition.RollingPolicy,
			PaasLagerDefinition.LogFormatText, PaasLagerDefinition.LogRotateDate,
			PaasLagerDefinition.LogRotateSize, PaasLagerDefinition.LogBackupCount)
	}

	if err = parseRouterConfig(fileutil.RouterDefinition()); err != nil {
		if os.IsNotExist(err) {
			lager.Logger.Infof("[%s] not exist", fileutil.RouterDefinition())
		} else {
			return err
		}
	}
	err = archaius.Init()
	if err != nil {
		return err
	}
	lager.Logger.Infof("archaius init success")

	var schemaError error

	//Upload schemas using environment variable SCHEMA_ROOT
	schemaEnv := archaius.GetString(common.EnvSchemaRoot, "")
	if schemaEnv != "" {
		schemaError = schema.LoadSchema(schemaEnv, true)
	} else {
		schemaError = schema.LoadSchema(fileutil.GetConfDir(), false)
	}

	if schemaError != nil {
		return schemaError
	}

	//set microservice names
	msError := schema.SetMicroServiceNames(fileutil.GetConfDir())
	if msError != nil {
		return msError
	}

	NodeIP = archaius.GetString(common.EnvNodeIP, "")
	err = parse()
	if err != nil {
		return err
	}

	SelfServiceName = MicroserviceDefinition.ServiceDescription.Name
	SelfVersion = MicroserviceDefinition.ServiceDescription.Version
	SelfMetadata = MicroserviceDefinition.ServiceDescription.Properties
	if GlobalDefinition.AppID == "" {
		GlobalDefinition.AppID = common.DefaultApp
	}

	return err
}
