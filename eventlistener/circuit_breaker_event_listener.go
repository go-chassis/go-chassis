package eventlistener

import (
	"github.com/go-chassis/go-archaius/event"
	"regexp"
	"strings"

	"github.com/go-chassis/go-chassis/control/servicecomb"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
)

// constants for consumer isolation, circuit breaker, fallback keys
const (
	// ConsumerIsolationKey is a variable of type string
	ConsumerIsolationKey      = "cse.isolation"
	ConsumerCircuitbreakerKey = "cse.circuitBreaker"
	ConsumerFallbackKey       = "cse.fallback"
	ConsumerFallbackPolicyKey = "cse.fallbackpolicy"
	regex4normal              = "cse\\.(isolation|circuitBreaker|fallback|fallbackpolicy)\\.Consumer\\.(.*)\\.(timeout|timeoutInMilliseconds|maxConcurrentRequests|enabled|forceOpen|forceClosed|sleepWindowInMilliseconds|requestVolumeThreshold|errorThresholdPercentage|enabled|maxConcurrentRequests|policy)\\.(.+)"
	regex4mesher              = "cse\\.(isolation|circuitBreaker|fallback|fallbackpolicy)\\.(.+)\\.Consumer\\.(.*)\\.(timeout|timeoutInMilliseconds|maxConcurrentRequests|enabled|forceOpen|forceClosed|sleepWindowInMilliseconds|requestVolumeThreshold|errorThresholdPercentage|enabled|maxConcurrentRequests|policy)\\.(.+)"
)

//CircuitBreakerEventListener is a struct with one string variable
type CircuitBreakerEventListener struct {
	Key string
}

//Event is a method which triggers flush circuit
func (el *CircuitBreakerEventListener) Event(e *event.Event) {
	openlogging.Info("circuit change e: %v", openlogging.WithTags(openlogging.Tags{
		"key": e.Key,
	}))
	if err := config.ReadHystrixFromArchaius(); err != nil {
		openlogging.Error("can not unmarshal new cb config: " + err.Error())
	}
	servicecomb.SaveToCBCache(config.GetHystrixConfig())
	switch e.EventType {
	case common.Update:
		FlushCircuitByKey(e.Key)
	case common.Create:
		FlushCircuitByKey(e.Key)
	case common.Delete:
		FlushCircuitByKey(e.Key)
	}
}

//FlushCircuitByKey is a function used to flush for a particular key
func FlushCircuitByKey(key string) {
	sourceName, serviceName := GetNames(key)
	cmdName := GetCircuitName(sourceName, serviceName)
	if cmdName == common.Consumer {
		openlogging.Info("Global Key changed For circuit: [" + cmdName + "], will flush all circuit")
		hystrix.Flush()
	} else {
		openlogging.Info("Specific Key changed For circuit: [" + cmdName + "], will only flush this circuit")
		hystrix.FlushByName(cmdName)
	}

}

//GetNames is function
func GetNames(key string) (string, string) {
	regNormal := regexp.MustCompile(regex4normal)
	regMesher := regexp.MustCompile(regex4mesher)
	var sourceName string
	var serviceName string
	if regNormal.MatchString(key) {
		s := regNormal.FindStringSubmatch(key)
		openlogging.Debug("Normal Key")
		return "", s[2]

	}
	if regMesher.MatchString(key) {
		s := regMesher.FindStringSubmatch(key)
		openlogging.Debug("Mesher Key")
		return s[2], s[3]
	}
	return sourceName, serviceName
}

//GetCircuitName is a function used to get circuit names
func GetCircuitName(sourceName, serviceName string) string {
	if sourceName != "" {
		return strings.Join([]string{sourceName, "Consumer", serviceName}, ".")
	}
	if sourceName == "" && serviceName != "" {
		return strings.Join([]string{"Consumer", serviceName}, ".")
	}
	if sourceName == "" && serviceName == "" {
		return common.Consumer
	}
	return ""
}
