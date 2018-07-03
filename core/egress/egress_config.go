package egress

import (
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"fmt"
	"errors"
	"regexp"
	"strings"
)


const (
	dns1123LabelMaxLength  int     = 63
	dns1123LabelFmt        string  = "[a-zA-Z0-9]([-a-z-A-Z0-9]*[a-zA-Z0-9])?"
	wildcardPrefix         string  = "(\\*)?" + dns1123LabelFmt
	DefaultRouterType              = "cse"

)
var (
	dns1123LabelRegexp   = regexp.MustCompile("^" + dns1123LabelFmt + "$")
	wildcardPrefixRegexp = regexp.MustCompile("^" + wildcardPrefix + "$")
)


// Init initialize Egress config
func Init() error {
	// init dests
	egressConfigFromFile := config.EgressDefinition
	BuildEgress(DefaultRouterType)

	if egressConfigFromFile != nil {
		if egressConfigFromFile.Destinations != nil {
			DefaultEgress.SetEgressRule(egressConfigFromFile.Destinations)
		}
	}

	DefaultEgress.Init()
	// storing the egress rules based on host in two maps
	// one host having wild card and other without wildcard
	plainHosts, regexHosts = SplitEgressRules()
	lager.Logger.Info("Egress init success")
	return nil
}

// ValidateRule validate the route rules of each service
func ValidateEgressRule(rules map[string][]*model.EgressRule) (bool, error) {
	for _, rule := range rules {
		for _, egressrule := range rule {
			if len(egressrule.Hosts) == 0 {
				return false, errors.New("Egress rule should have atleast one host")
			}
			for _, host :=  range egressrule.Hosts{
				err:= ValidateHostName(host)
				if err != nil {
					return false, err
				}
			}
		}

	}
	return true, nil
}


func ValidateHostName(host string) (error){
	if len(host) > 255 {
		return  fmt.Errorf("host name %q too long (max 255)", host)
	}
	if len(host) == 0 {
		return  fmt.Errorf("empty host name not allowed")
	}

	parts := strings.SplitN(host, ".", 2)
	if !IsWildcardDNS1123Label(parts[0]) {
		return fmt.Errorf("host name %q invalid (label %q invalid)", host, parts[0])
	} else if len(parts) > 1 {
		err := validateDNS1123Labels(parts[1])
			return err
	}

	return nil

}

func IsWildcardDNS1123Label(value string) bool {
	return len(value) <= dns1123LabelMaxLength && wildcardPrefixRegexp.MatchString(value)
}
func validateDNS1123Labels(host string) error {
	for _, label := range strings.Split(host, ".") {
		if !IsDNS1123Label(label) {
			return fmt.Errorf("host name %q invalid (label %q invalid)", host, label)
		}
	}
	return nil
}

func IsDNS1123Label(value string) bool {
	return len(value) <= dns1123LabelMaxLength && dns1123LabelRegexp.MatchString(value)
}