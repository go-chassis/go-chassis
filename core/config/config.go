package config

import (
	"errors"
	"os"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/config/schema"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"github.com/go-chassis/go-chassis/pkg/util/iputil"
	"github.com/go-mesh/openlogging"
)

// GlobalDefinition is having the information about region, load balancing, service center, config server,
// protocols, and handlers for the micro service
var GlobalDefinition *model.GlobalCfg
var lbConfig *model.LBWrapper

// MicroserviceDefinition has info about application id, provider info, description of the service,
// and description of the instance
var MicroserviceDefinition *model.ServiceSpec

//MonitorCfgDef has monitor info, including zipkin and apm.
var MonitorCfgDef *model.MonitorCfg

//HystrixConfig is having info about isolation, circuit breaker, fallback properities of the micro service
var HystrixConfig *model.HystrixConfigWrapper

// ErrNoName is used to represent the service name missing error
var ErrNoName = errors.New("micro service name is missing in description file")

//GetConfigServerConf return config server conf
func GetConfigServerConf() model.ConfigClient {
	return GlobalDefinition.ServiceComb.Config.Client
}

//GetTransportConf return transport settings
func GetTransportConf() model.Transport {
	return GlobalDefinition.ServiceComb.Transport
}

//GetDataCenter return data center info
func GetDataCenter() *model.DataCenterInfo {
	return GlobalDefinition.DataCenter
}

//GetAPM return monitor config info
func GetAPM() model.APMStruct {
	return MonitorCfgDef.ServiceComb.APM
}

// readFromArchaius unmarshal configurations to expected pointer
func readFromArchaius() error {
	openlogging.Debug("read from archaius")
	err := ReadGlobalConfigFromArchaius()
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

	populateConfigServerAddress()
	populateServiceRegistryAddress()
	err = ReadMonitorFromArchaius()
	if err != nil {
		return err
	}

	populateServiceEnvironment()
	populateServiceName()
	populateVersion()
	populateApp()

	return nil
}

// populateServiceRegistryAddress populate service registry address
func populateServiceRegistryAddress() {
	//Registry Address , higher priority for environment variable
	registryAddrFromEnv := readEndpoint(common.EnvSCEndpoint)
	if registryAddrFromEnv != "" {
		openlogging.Debug("detect env", openlogging.WithTags(
			openlogging.Tags{
				"ep": registryAddrFromEnv,
			}))
		GlobalDefinition.ServiceComb.Registry.Address = registryAddrFromEnv
	}
}

// populateConfigServerAddress populate config server address
func populateConfigServerAddress() {
	//config server Address , higher priority for environment variable
	configServerAddrFromEnv := readEndpoint(common.EnvCCEndpoint)
	if configServerAddrFromEnv != "" {
		GlobalDefinition.ServiceComb.Config.Client.ServerURI = configServerAddrFromEnv
	}
}

// readEndpoint
func readEndpoint(env string) string {
	addrFromEnv := archaius.GetString(env, archaius.GetString(common.EnvCSEEndpoint, ""))
	if addrFromEnv != "" {
		openlogging.Info("read config " + addrFromEnv)
		return addrFromEnv
	}
	return addrFromEnv
}

// populateServiceEnvironment populate service environment
func populateServiceEnvironment() {
	if e := archaius.GetString(common.Env, ""); e != "" {
		MicroserviceDefinition.Environment = e
	}
}

// populateServiceName populate service name
func populateServiceName() {
	if e := archaius.GetString(common.ServiceName, ""); e != "" {
		MicroserviceDefinition.Name = e
	}
}

// populateVersion populate version
func populateVersion() {
	if e := archaius.GetString(common.Version, ""); e != "" {
		MicroserviceDefinition.Version = e
	}
}

func populateApp() {
	if e := archaius.GetString(common.App, ""); e != "" {
		MicroserviceDefinition.Name = e
	}
}

// ReadGlobalConfigFromArchaius for to unmarshal the global config file(chassis.yaml) information
func ReadGlobalConfigFromArchaius() error {
	GlobalDefinition = &model.GlobalCfg{}
	err := archaius.UnmarshalConfig(&GlobalDefinition)
	if err != nil {
		return err
	}
	MicroserviceDefinition = &GlobalDefinition.ServiceComb.ServiceDescription
	return nil
}

// ReadLBFromArchaius for to unmarshal the global config file(chassis.yaml) information
func ReadLBFromArchaius() error {
	lbMutex.Lock()
	defer lbMutex.Unlock()
	lbConfig = &model.LBWrapper{}
	err := archaius.UnmarshalConfig(lbConfig)
	if err != nil {
		return err
	}
	return nil
}

//ReadMonitorFromArchaius read monitor config from archauis pkg
func ReadMonitorFromArchaius() error {
	MonitorCfgDef = &model.MonitorCfg{}
	err := archaius.UnmarshalConfig(&MonitorCfgDef)
	if err != nil {
		openlogging.Error("Config init failed. " + err.Error())
		return err
	}
	return nil
}

// ReadHystrixFromArchaius is unmarshal hystrix configuration file(circuit_breaker.yaml)
func ReadHystrixFromArchaius() error {
	cbMutex.RLock()
	defer cbMutex.RUnlock()
	HystrixConfig = &model.HystrixConfigWrapper{}
	err := archaius.UnmarshalConfig(&HystrixConfig)
	if err != nil {
		return err
	}
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
	if HystrixConfig != nil {
		return HystrixConfig.HystrixConfig
	}
	return nil
}

// Init is initialize the configuration directory, archaius, route rule, and schema
func Init() error {
	err := InitArchaius()
	if err != nil {
		return err
	}

	//Upload schemas using environment variable SCHEMA_ROOT
	schemaPath := archaius.GetString(common.EnvSchemaRoot, "")
	if schemaPath == "" {
		schemaPath = fileutil.GetConfDir()
	}

	schemaError := schema.LoadSchema(schemaPath)
	if schemaError != nil {
		return schemaError
	}

	//set micro service names
	err = schema.SetMicroServiceNames(schemaPath)
	if err != nil {
		return err
	}

	runtime.NodeIP = archaius.GetString(common.EnvNodeIP, "")

	err = readFromArchaius()
	if err != nil {
		return err
	}

	runtime.ServiceName = MicroserviceDefinition.Name
	runtime.Version = MicroserviceDefinition.Version
	runtime.Environment = MicroserviceDefinition.Environment
	runtime.MD = MicroserviceDefinition.Properties
	runtime.App = MicroserviceDefinition.AppID
	if runtime.App == "" {
		runtime.App = common.DefaultApp
	}

	runtime.HostName = MicroserviceDefinition.Hostname
	if runtime.HostName == "" {
		runtime.HostName, err = os.Hostname()
		if err != nil {
			openlogging.Error("Get hostname failed:" + err.Error())
			return err
		}
	} else if runtime.HostName == common.PlaceholderInternalIP {
		runtime.HostName = iputil.GetLocalIP()
	}
	openlogging.Info("Host name is " + runtime.HostName)
	return err
}
