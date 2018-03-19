package configcenter

import "crypto/tls"

var configClientPlugins = make(map[string]func(endpoint, serviceName, app, env, version string, tlsConfig *tls.Config) ConfigClient)

//DefaultClient is config server's client
var DefaultClient ConfigClient

const (
	defaultConfigServer = "config_center"
)

//InstallConfigClientPlugin install a config client plugin
func InstallConfigClientPlugin(name string, f func(endpoint, serviceName, app, env, version string, tlsConfig *tls.Config) ConfigClient) {
	configClientPlugins[name] = f
}

//ConfigClient is the interface of config server client, it has basic func to interact with config server
type ConfigClient interface {
	//PullConfigs pull all configs from remote
	PullConfigs(serviceName, version, app, env string)
	//PullConfig pull one config from remote
	PullConfig(serviceName, version, app, env, key, contentType string)
}

//Enable enable config server client
func Enable() {
	//TODO read plugin name from config file,default is config center
	//TODO create DefaultClient,config source should only use default client to interact with config server
}
