package fileutil

import (
	"os"
	"path/filepath"
)

const (
	//ChassisConfDir is constant of type string
	ChassisConfDir  = "CHASSIS_CONF_DIR"
	ChassisHome     = "CHASSIS_HOME"
	SchemaDirectory = "schema"
)

const (
	//Global is a constant of type string
	Global        = "chassis.yaml"
	LoadBalancing = "load_balancing.yaml"
	RateLimiting  = "rate_limiting.yaml"
	Definition    = "microservice.yaml"
	Hystric       = "circuit_breaker.yaml"
	PaasLager     = "lager.yaml"
	TLS           = "tls.yaml"
	Monitoring    = "monitoring.yaml"
	Auth          = "auth.yaml"
	Tracing       = "tracing.yaml"
)

var configDir string
var homeDir string

//GetWorkDir is a function used to get the working directory
func GetWorkDir() (string, error) {
	wd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return wd, nil
}

func initDir() error {
	if h := os.Getenv(ChassisHome); h != "" {
		homeDir = h
	} else {
		wd, err := GetWorkDir()
		if err != nil {
			return err
		}
		homeDir = wd
	}

	// set conf dir, CHASSIS_CONF_DIR has highest priority
	if confDir := os.Getenv(ChassisConfDir); confDir != "" {
		configDir = confDir
	} else {
		// CHASSIS_HOME has second most high priority
		configDir = filepath.Join(homeDir, "conf")
	}
	return nil
}

//ChassisHomeDir is function used to get the home directory of chassis
func ChassisHomeDir() string {
	return homeDir
}

//ConfDir is a function used to get the configuration directory
func ConfDir() string {
	return configDir
}

//HystrixDefinition is a function used to join .yaml file name with configuration path
func HystrixDefinition() string {
	return filepath.Join(configDir, Hystric)
}

//GetDefinition is a function used to join .yaml file name with configuration path
func GetDefinition() string {
	return filepath.Join(configDir, Definition)
}

//GetLoadBalancing is a function used to join .yaml file name with configuration directory
func GetLoadBalancing() string {
	return filepath.Join(configDir, LoadBalancing)
}

//GetRateLimiting is a function used to join .yaml file name with configuration directory
func GetRateLimiting() string {
	return filepath.Join(configDir, RateLimiting)
}

//GetTLS is a function used to join .yaml file name with configuration directory
func GetTLS() string {
	return filepath.Join(configDir, TLS)
}

//GetMonitoring is a function used to join .yaml file name with configuration directory
func GetMonitoring() string {
	return filepath.Join(configDir, Monitoring)
}

//MicroserviceDefinition is a function used to join .yaml file name with configuration directory
func MicroserviceDefinition(microserviceName string) string {
	return filepath.Join(configDir, microserviceName, Definition)
}

//GetMicroserviceDesc is a function used to join .yaml file name with configuration directory
func GetMicroserviceDesc() string {
	return filepath.Join(configDir, Definition)
}

//GlobalDefinition is a function used to join .yaml file name with configuration directory
func GlobalDefinition() string {
	return filepath.Join(configDir, Global)
}

//PassLagerDefinition is a function used to join .yaml file name with configuration directory
func PassLagerDefinition() string {
	return filepath.Join(configDir, PaasLager)
}

//GetAuth is a function used to join .yaml file name with configuration directory
func GetAuth() string {
	return filepath.Join(configDir, Auth)
}

//GetTracing is a function used to join .yaml file name with configuration directory
func GetTracing() string {
	return filepath.Join(configDir, Tracing)
}

//SchemaDir is a function used to join .yaml file name with configuration path
func SchemaDir(microserviceName string) string {
	return filepath.Join(configDir, microserviceName, SchemaDirectory)
}

//InitConfigDir is a function used to initialize configuration directory
func InitConfigDir() error {
	if err := initDir(); err != nil {
		return err
	}

	return nil
}
