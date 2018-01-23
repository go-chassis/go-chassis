package archaius

import (
	"github.com/ServiceComb/go-archaius"
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-archaius/sources/external-source"
	"github.com/ServiceComb/go-archaius/sources/file-source"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/util/fileutil"
)

// Config is the struct of configuration files, and configuration factory
type Config struct {
	ConfigFiles   []string
	ConfigFactory goarchaius.ConfigurationFactory
}

// DefaultConf is the variable having information about configuration files, and configuration factory
var DefaultConf *Config

// NewConfig is gives the object of Config(it is having configuration files, and configuration factory)
func NewConfig(essentialfiles, commonfiles []string) (*Config, error) {
	// created config factory object
	factory, err := goarchaius.NewConfigFactory()
	if err != nil {
		return nil, err
	}
	factory.DeInit()
	factory.Init()

	files := make([]string, 0)
	// created file source object
	fileSource := filesource.NewYamlConfigurationSource()
	// adding all files with file source
	for _, v := range essentialfiles {
		if err := fileSource.AddFileSource(v, filesource.DefaultFilePriority); err != nil {
			lager.Logger.Errorf(err, "add file source error.")
			return nil, err
		}
		files = append(files, v)
	}
	for _, v := range commonfiles {
		if err := fileSource.AddFileSource(v, filesource.DefaultFilePriority); err != nil {
			lager.Logger.Infof("%v", err)
		}
		files = append(files, v)
	}

	err = factory.AddSource(fileSource)
	if err != nil {
		return nil, err
	}
	conf := &Config{
		ConfigFiles:   files,
		ConfigFactory: factory,
	}

	eventHandler := EventListener{
		Name:    "EventHandler",
		Factory: factory,
	}

	factory.RegisterListener(eventHandler, "a*")

	return conf, nil
}

// EventListener is a struct having information about registering key and object
type EventListener struct {
	Name    string
	Factory goarchaius.ConfigurationFactory
}

// Event is invoked while generating events at run time
func (e EventListener) Event(event *core.Event) {
	value := e.Factory.GetConfigurationByKey(event.Key)
	lager.Logger.Infof("config value after change %s | %s", event.Key, value)
}

// Init is to initialize the archaius
func Init() error {
	essentialfiles := []string{
		fileutil.GlobalDefinition(),
		fileutil.GetMicroserviceDesc(),
	}
	commonfiles := []string{
		fileutil.HystrixDefinition(),
		fileutil.GetLoadBalancing(),
		fileutil.GetRateLimiting(),
		fileutil.GetTLS(),
		fileutil.GetMonitoring(),
		fileutil.GetAuth(),
		fileutil.GetTracing(),
	}

	lager.Logger.Infof("Essential Configuration Path: %v, Configuration Paths %v", essentialfiles, commonfiles)
	dConf, err := NewConfig(essentialfiles, commonfiles)
	DefaultConf = dConf
	return err
}

// Get is for to get the value of configuration key
func Get(key string) interface{} {
	return DefaultConf.ConfigFactory.GetConfigurationByKey(key)
}

// Exist is check the configuration key existence
func Exist(key string) bool {
	return DefaultConf.ConfigFactory.IsKeyExist(key)
}

// UnmarshalConfig is for unmarshalling the configuraions of receiving object
func UnmarshalConfig(obj interface{}) error {
	return DefaultConf.ConfigFactory.Unmarshal(obj)
}

// GetBool is gives the key value in the form of bool
func GetBool(key string, defaultValue bool) bool {
	b, err := DefaultConf.ConfigFactory.GetValue(key).ToBool()
	if err != nil {
		return defaultValue
	}
	return b
}

// GetFloat64 gives the key value in the form of float64
func GetFloat64(key string, defaultValue float64) float64 {
	result, err := DefaultConf.ConfigFactory.GetValue(key).ToFloat64()
	if err != nil {
		return defaultValue
	}
	return result
}

// GetInt gives the key value in the form of GetInt
func GetInt(key string, defaultValue int) int {
	result, err := DefaultConf.ConfigFactory.GetValue(key).ToInt()
	if err != nil {
		return defaultValue
	}
	return result
}

// GetString gives the key value in the form of GetString
func GetString(key string, defaultValue string) string {
	result, err := DefaultConf.ConfigFactory.GetValue(key).ToString()
	if err != nil {
		return defaultValue
	}
	return result
}

// GetConfigs gives the information about all configurations
func GetConfigs() map[string]interface{} {
	return DefaultConf.ConfigFactory.GetConfigurations()
}

// GetStringByDI is for to get the value of configuration key based on dimension info
func GetStringByDI(dimensionInfo, key string, defaultValue string) string {
	result, err := DefaultConf.ConfigFactory.GetValueByDI(dimensionInfo, key).ToString()
	if err != nil {
		return defaultValue
	}
	return result
}

// GetConfigsByDI is for to get the all configurations received dimensionInfo
func GetConfigsByDI(dimensionInfo string) map[string]interface{} {
	return DefaultConf.ConfigFactory.GetConfigurationsByDimensionInfo(dimensionInfo)
}

// AddDI adds a NewDimensionInfo of which configurations needs to be taken
func AddDI(dimensionInfo string) (map[string]string, error) {
	config, err := DefaultConf.ConfigFactory.AddByDimensionInfo(dimensionInfo)
	return config, err
}

//RegisterListener to Register all listener for different key changes, each key could be a regular expression
func RegisterListener(listenerObj core.EventListener, key ...string) error {
	return DefaultConf.ConfigFactory.RegisterListener(listenerObj, key...)
}

// UnRegisterListener is to remove the listener
func UnRegisterListener(listenerObj core.EventListener, key ...string) error {
	return DefaultConf.ConfigFactory.UnRegisterListener(listenerObj, key...)
}

// AddFile is for to add the configuration files into the configfactory at run time
func AddFile(file string) error {
	return filesource.NewYamlConfigurationSource().AddFileSource(file, filesource.DefaultFilePriority)
}

// AddKeyValue is for to add the configuration key, value pairs into the configfactory at run time
func AddKeyValue(key string, value interface{}) error {
	return externalconfigsource.NewExternalConfigurationSource().AddKeyValue(key, value)
}
