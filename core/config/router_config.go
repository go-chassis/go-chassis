package config

import (
	"encoding/json"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/util/fileutil"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// constant for route rule keys
const (
	DarkLaunchPrefix   = "cse.darklaunch.policy."
	DarkLaunchTypeRule = "RULE"
	DarkLaunchTypeRate = "RATE"
	TemplateKey        = "sourceTemplate"
)

// routerConfig variable info about route rule, and source template
var routerConfig *RouterConfig

// GetRouterConfig get router configuurations
func GetRouterConfig() *RouterConfig {
	return routerConfig
}

// InitRouter initialize router
func InitRouter() error {
	routerConfig = &RouterConfig{}
	routerRule := &RouterConfig{}
	routerTemplates := &RouterConfig{}
	routerRule, err := GetRouteRules()
	if err != nil {
		return err
	}

	contents, err := GetConfigContents(TemplateKey)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(contents), routerTemplates); err != nil {
		return err
	}

	routerConfig.Destinations = routerRule.Destinations
	routerConfig.SourceTemplates = routerTemplates.SourceTemplates
	return nil
}

// GetRouteRules get route rules
func GetRouteRules() (*RouterConfig, error) {
	f := fileutil.GetRouter()
	routeRules := &RouterConfig{
		Destinations: map[string][]*RouteRule{},
	}
	var contents string

	configMap := archaius.GetConfigs()
	//filter out key:value pairs which are not route rules
	for k := range configMap {
		if !strings.HasPrefix(k, DarkLaunchPrefix) {
			delete(configMap, k)
		}
	}
	//there is no route rules in archaius, so load rules from config file
	if len(configMap) == 0 {
		lager.Logger.Info("Load route rules from file: " + f)
		b, err := ioutil.ReadFile(f)
		if err != nil {
			lager.Logger.Warn("Can not read file", err)
			contents = ""
			return nil, err
		}
		contents = string(b)
		err = yaml.Unmarshal([]byte(contents), routeRules)
		return routeRules, err
	}
	//put route rules in configMap into routeRules
	rule := &DarkLaunchRule{}
	for k, v := range configMap {
		if err := json.Unmarshal([]byte(v.(string)), rule); err != nil {
			return routeRules, err
		}
		key := strings.Replace(k, DarkLaunchPrefix, "", 1)
		routeRules.Destinations[key] = TranslateRules(rule)
	}
	return routeRules, nil
}

// GetConfigContents get configuration contents
func GetConfigContents(key string) (string, error) {
	var err error
	var contents string
	//route rule yaml file's content is value of a key
	//So read from config center first,if it is empty, Try to set file content into memory key value
	contents = archaius.GetString(key, "")
	if contents == "" {
		contents, err = SetKeyValueByFile(key, fileutil.GetRouter())
		if err != nil {
			return "", err
		}
	}
	return contents, nil
}

// SetKeyValueByFile is for adding configurations of the file to the archaius through external configuration source
func SetKeyValueByFile(key, f string) (string, error) {
	var contents string
	if _, err := os.Stat(f); err != nil {
		lager.Logger.Warn(err.Error(), nil)
		return "", err
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		lager.Logger.Error("Can not read router.yaml", err)
		return "", err
	}
	contents = string(b)
	archaius.AddKeyValue(key, contents)
	return contents, nil
}

// TranslateRules translate rules
func TranslateRules(rule *DarkLaunchRule) []*RouteRule {
	if rule.Type == DarkLaunchTypeRate {
		routes := make([]*RouteTag, len(rule.Items))
		for i, v := range rule.Items {
			weight, _ := strconv.Atoi(v.PolicyCondition)
			version := strings.Replace(v.GroupCondition, "version=", "", 1)
			routes[i] = &RouteTag{
				Weight: weight,
				Tags:   map[string]string{"version": version},
			}
		}
		return []*RouteRule{{
			Routes:     routes,
			Precedence: 1,
		}}
	}
	if rule.Type == DarkLaunchTypeRule {
		rules := make([]*RouteRule, len(rule.Items))
		for i, v := range rule.Items {
			con := v.PolicyCondition
			version := strings.Replace(v.GroupCondition, "version=", "", 1)
			match := Match{
				HTTPHeaders: map[string]map[string]string{},
				Headers:     map[string]map[string]string{},
			}
			if strings.Contains(con, "!=") {
				match.HTTPHeaders[strings.Split(con, "!=")[0]] = map[string]string{"noEqu": strings.Split(con, "!=")[1]}
				match.Headers[strings.Split(con, "!=")[0]] = map[string]string{"noEqu": strings.Split(con, "!=")[1]}
			} else if strings.Contains(con, ">=") {
				match.HTTPHeaders[strings.Split(con, ">=")[0]] = map[string]string{"noLess": strings.Split(con, ">=")[1]}
				match.Headers[strings.Split(con, ">=")[0]] = map[string]string{"noLess": strings.Split(con, ">=")[1]}
			} else if strings.Contains(con, "<=") {
				match.HTTPHeaders[strings.Split(con, "<=")[0]] = map[string]string{"noGreater": strings.Split(con, "<=")[1]}
				match.Headers[strings.Split(con, "<=")[0]] = map[string]string{"noGreater": strings.Split(con, "<=")[1]}
			} else if strings.Contains(con, "=") {
				match.HTTPHeaders[strings.Split(con, "=")[0]] = map[string]string{"exact": strings.Split(con, "=")[1]}
				match.Headers[strings.Split(con, "=")[0]] = map[string]string{"exact": strings.Split(con, "=")[1]}
			} else if strings.Contains(con, ">") {
				match.HTTPHeaders[strings.Split(con, ">")[0]] = map[string]string{"greater": strings.Split(con, ">")[1]}
				match.Headers[strings.Split(con, ">")[0]] = map[string]string{"greater": strings.Split(con, ">")[1]}
			} else if strings.Contains(con, "<") {
				match.HTTPHeaders[strings.Split(con, "<")[0]] = map[string]string{"less": strings.Split(con, "<")[1]}
				match.Headers[strings.Split(con, "<")[0]] = map[string]string{"less": strings.Split(con, "<")[1]}
			} else if strings.Contains(con, "~") {
				match.HTTPHeaders[strings.Split(con, "~")[0]] = map[string]string{"regex": strings.Split(con, "~")[1]}
				match.Headers[strings.Split(con, "~")[0]] = map[string]string{"regex": strings.Split(con, "~")[1]}
			}
			newRule := &RouteRule{
				Routes:     GenerateRouteTags(strings.Split(version, ",")),
				Match:      match,
				Precedence: 1,
			}
			rules[i] = newRule
		}
		return rules
	}
	return nil
}

// GenerateRouteTags generate route tags
func GenerateRouteTags(versions []string) []*RouteTag {
	length := len(versions)
	if length == 1 {
		return []*RouteTag{{
			Weight: 100,
			Tags:   map[string]string{"version": versions[0]},
		}}
	}

	tags := make([]*RouteTag, length)
	for i, v := range versions {
		tags[i] = &RouteTag{
			Weight: 100 / length,
			Tags:   map[string]string{"version": v},
		}
	}
	return tags
}
