# Archaius Config Source Plugin

## Config Source Plugin
Config Source Plugin let's you write your own the Config-Center client implementation for the different types of Config Source.

## Instructions
Go-Chassis can support pulling the configuration from different types of config centers, currently there are 2 
implementation available for Config-Client Plugin (Go-Archaius Config-center, Ctrip Apollo Config-center). If you want to 
implement any new client for another config-center then you have to implement the following ConfigClient Interface.

```go
//ConfigClient is the interface of config server client, it has basic func to interact with config server
type ConfigClient interface {
	//Init the Configuration for the Server
	Init()
	//PullConfigs pull all configs from remote
	PullConfigs(serviceName, version, app, env string) (map[string]interface{}, error)
	//PullConfig pull one config from remote
	PullConfig(serviceName, version, app, env, key, contentType string) (interface{}, error)
	//PullConfigsByDI pulls the configurations with customized DimensionInfo/Project
	PullConfigsByDI(dimensionInfo , diInfo string)(map[string]map[string]interface{}, error)
}
```

Once you implement the above interface then you need to define the type of your configClient

```go
config.GlobalDefinition.Cse.Config.Client.Type
```

```go
cse:
  config:
    client:
      type: your_client_name   #config_center/apollo/your_client_name
```
Based on this type you need to load the plugin in your init()

```go
 func init(){
 	client.InstallConfigClientPlugin("NameOfYourPLugin", InitConfigYourPlugin)
 }

```
You need to import the package path in your application to Enable your plugin. Once the plugin is enabled then you can pull configuration using the ConfigClient

```go
client.DefaultClient.PullConfigs(serviceName, versionName, appName, env)
```
