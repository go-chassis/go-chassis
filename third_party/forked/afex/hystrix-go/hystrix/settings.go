package hystrix

import (
	"sync"
	"time"
)

var (
	// DefaultTimeout is how long to wait for command to complete, in milliseconds
	DefaultTimeout = 1000

	// DefaultMaxConcurrent is how many commands of the same type can run at the same time
	DefaultMaxConcurrent = 10

	// DefaultVolumeThreshold is the minimum number of requests needed before a circuit can be tripped due to health
	DefaultVolumeThreshold = 20

	// DefaultSleepWindow is how long, in milliseconds, to wait after a circuit opens before testing for recovery
	DefaultSleepWindow = 5000

	// DefaultErrorPercentThreshold causes circuits to open once the rolling measure of errors exceeds this percent of requests
	DefaultErrorPercentThreshold = 50
)

type Settings struct {
	// isolation 属性
	Timeout               time.Duration
	TimeoutEnabled        bool
	MaxConcurrentRequests int

	// circuit break 属性
	CircuitBreakerEnabled  bool
	RequestVolumeThreshold uint64
	SleepWindow            time.Duration
	ErrorPercentThreshold  int

	//动态治理
	ForceFallback bool
	ForceOpen     bool
	ForceClose    bool
}

// CommandConfig is used to tune circuit settings at runtime
type CommandConfig struct {
	TimeoutEnabled         bool
	Timeout                int `json:"timeout"`
	MaxConcurrentRequests  int `json:"max_concurrent_requests"`
	RequestVolumeThreshold int `json:"request_volume_threshold"`
	SleepWindow            int `json:"sleep_window"`
	ErrorPercentThreshold  int `json:"error_percent_threshold"`
	//动态治理
	ForceFallback         bool
	CircuitBreakerEnabled bool
	ForceOpen             bool
	ForceClose            bool
}

var circuitSettings map[string]*Settings
var settingsMutex *sync.RWMutex

func init() {
	// 配置文件中优先级Operation > Schema > Default,实现源码中配置属性
	circuitSettings = make(map[string]*Settings)
	settingsMutex = &sync.RWMutex{}
}

// Configure applies settings for a set of circuits
func Configure(cmds map[string]CommandConfig) {
	for k, v := range cmds {
		ConfigureCommand(k, v)
	}
}

type CommandConfigOption func(*CommandConfig)

// 新建一个CommandConfig返回
func NewCommandConfig(opt ...CommandConfigOption) CommandConfig {
	cmdconfig := CommandConfig{}

	for _, o := range opt {
		o(&cmdconfig)
	}

	return cmdconfig
}

func WithTimeOut(timeout int) CommandConfigOption {
	return func(c *CommandConfig) {
		c.Timeout = timeout
	}
}

func WithMaxRequests(maxrequests int) CommandConfigOption {
	return func(c *CommandConfig) {
		c.MaxConcurrentRequests = maxrequests
	}
}

func WithVolumeThreshold(volumethreshold int) CommandConfigOption {
	return func(c *CommandConfig) {
		c.RequestVolumeThreshold = volumethreshold
	}
}

func WithSleepWindow(sleepwindow int) CommandConfigOption {
	return func(c *CommandConfig) {
		c.SleepWindow = sleepwindow
	}
}

func WithErrorPercent(errorpercent int) CommandConfigOption {
	return func(c *CommandConfig) {
		c.ErrorPercentThreshold = errorpercent
	}
}

// ConfigureCommand applies settings for a circuit
func ConfigureCommand(name string, config CommandConfig) {

	settingsMutex.Lock()
	defer settingsMutex.Unlock()

	timeout := DefaultTimeout
	if config.Timeout != 0 {
		timeout = config.Timeout
	}

	max := DefaultMaxConcurrent
	if config.MaxConcurrentRequests != 0 {
		max = config.MaxConcurrentRequests
	}

	volume := DefaultVolumeThreshold
	if config.RequestVolumeThreshold != 0 {
		volume = config.RequestVolumeThreshold
	}

	sleep := DefaultSleepWindow
	if config.SleepWindow != 0 {
		sleep = config.SleepWindow
	}

	errorPercent := DefaultErrorPercentThreshold
	if config.ErrorPercentThreshold != 0 {
		errorPercent = config.ErrorPercentThreshold
	}
	circuitSettings[name] = &Settings{
		TimeoutEnabled:         config.TimeoutEnabled,
		ForceClose:             config.ForceClose,
		ForceOpen:              config.ForceOpen,
		CircuitBreakerEnabled:  config.CircuitBreakerEnabled,
		Timeout:                time.Duration(timeout) * time.Millisecond,
		MaxConcurrentRequests:  max,
		RequestVolumeThreshold: uint64(volume),
		SleepWindow:            time.Duration(sleep) * time.Millisecond,
		ErrorPercentThreshold:  errorPercent,
		ForceFallback:          config.ForceFallback,
	}
}

func getSettings(name string) *Settings {
	settingsMutex.RLock()
	s, exists := circuitSettings[name]
	settingsMutex.RUnlock()

	if !exists {
		ConfigureCommand(name, CommandConfig{})
		s = getSettings(name)
	}

	return s
}

func GetCircuitSettings() map[string]*Settings {
	copy := make(map[string]*Settings)

	settingsMutex.RLock()
	for key, val := range circuitSettings {
		copy[key] = val
	}
	settingsMutex.RUnlock()

	return copy
}
