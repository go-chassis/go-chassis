package egress

import (
	"github.com/ServiceComb/go-chassis/core/config/model"
	"sync"
	"regexp"
	"errors"
)

var lock sync.RWMutex

var plainHosts = make(map[string]*model.EgressRule)
var regexHosts = make(map[string]*model.EgressRule)



//Egress return egress rule, you can also set custom egress rule
type Egress interface {
	Init()error
	SetEgressRule(map[string][]*model.EgressRule)
	FetchEgressRule() map[string][]*model.EgressRule
	FetchEgressRuleByName(string) []*model.EgressRule
}

// ErrNoExist means if there is no egress implementation
var ErrNoExist = errors.New("Egress not exists")
var egressServices = make(map[string]func() (Egress, error))

// DefaultEgress is current egress implementation
var DefaultEgress Egress

// InstallEgressService install router service for developer
func InstallEgressService(name string, f func() (Egress, error)) {
	egressServices[name] = f
}

//BuildEgress create a Egress
func BuildEgress(name string) error {
	f, ok := egressServices[name]
	if !ok {
		return ErrNoExist
	}
	r, err := f()
	if err != nil {
		return err
	}
	DefaultEgress = r
	return nil
}

//Check Egress rule matches
func Match(hostname string) (bool, *model.EgressRule){
	EgressRules := DefaultEgress.FetchEgressRule()
	for _, egressRules := range EgressRules {
		for _, egress := range  egressRules {
			for _, host := range egress.Hosts {
				// Check host length greater than 0 and does not
				// start with *
				if len(host) > 0 &&  string(host[0]) != "*"{
						if host == hostname {
							return true, egress
						}
					} else if string(host[0]) == "*" {
						substring := host[1:]
						match, _ := regexp.MatchString(substring+"$", hostname)
						if match == true {
							return true, egress
						}
					}
				}
			}
		}

	return false, nil
}

//Check Egress rule matches
func SplitEgressRules() (map[string]*model.EgressRule, map[string]*model.EgressRule){
	EgressRules := DefaultEgress.FetchEgressRule()
	for _, egressRules := range EgressRules {
		for _, egress := range  egressRules{

			for _, host := range egress.Hosts{
				if len(host) > 1 && string(host[0]) != "*" {
					plainHosts[host] = egress
				}else if string(host[0]) == "*"{
					substring := host[1:]
					regexHosts[substring] = egress
				}
			}
		}
	}

	return plainHosts, regexHosts
}

func MatchHost(hostname string)(bool, *model.EgressRule){
	if val, ok := plainHosts[hostname]; ok{
		return true, val
	}

	for key, value := range regexHosts {
			match, _ := regexp.MatchString(key+"$", hostname);
			if match == true {
				return true, value

		}
	}
	return false, nil
}