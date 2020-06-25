package config

import (
	stringutil "github.com/go-chassis/go-chassis/pkg/string"
	"gopkg.in/yaml.v2"
)

//OneServiceRule save route rule for one service
type OneServiceRule []*RouteRule

//Len return the length of rule
func (o OneServiceRule) Len() int {
	return len(o)
}

//Value return the rule
func (o OneServiceRule) Value() []*RouteRule {
	return o
}

//NewServiceRule create a rule by raw data
func NewServiceRule(raw string) (*OneServiceRule, error) {
	b := stringutil.Str2bytes(raw)
	r := &OneServiceRule{}
	err := yaml.Unmarshal(b, r)
	return r, err
}

//ServiceComb hold all config items
type ServiceComb struct {
	Prefix Prefix `yaml:"servicecomb"`
}

//Prefix hold all config items
type Prefix struct {
	RouteRule       map[string]string `yaml:"routeRule"`      //service name is key,value is route rule yaml config
	SourceTemplates map[string]string `yaml:"sourceTemplate"` //template name is key, value is template policy
}

// Router define where rule comes from
type Router struct {
	Infra   string `yaml:"infra"`
	Address string `yaml:"address"`
}

// RouteRule is having route rule parameters
type RouteRule struct {
	Precedence int         `json:"precedence" yaml:"precedence"`
	Routes     []*RouteTag `json:"route" yaml:"route"`
	Match      Match       `json:"match" yaml:"match"`
}

// RouteTag gives route tag information
type RouteTag struct {
	Tags   map[string]string `json:"tags" yaml:"tags"`
	Weight int               `json:"weight" yaml:"weight"`
	Label  string
}

// Match is checking source, source tags, and http headers
type Match struct {
	Refer       string                       `json:"refer" yaml:"refer"`
	Source      string                       `json:"source" yaml:"source"`
	SourceTags  map[string]string            `json:"sourceTags" yaml:"sourceTags"`
	HTTPHeaders map[string]map[string]string `json:"httpHeaders" yaml:"httpHeaders"`
	Headers     map[string]map[string]string `json:"headers" yaml:"headers"`
}

//DarkLaunchRule dark launch rule
//Deprecated
type DarkLaunchRule struct {
	Type  string      `json:"policyType"` // RULE/RATE
	Items []*RuleItem `json:"ruleItems"`
}

//RuleItem rule item
//Deprecated
type RuleItem struct {
	GroupName       string   `json:"groupName"`
	GroupCondition  string   `json:"groupCondition"`  // version=0.0.1
	PolicyCondition string   `json:"policyCondition"` // 80/test!=2
	CaseInsensitive bool     `json:"caseInsensitive"`
	Versions        []string `json:"versions"`
}

//MatchPolicy specify a request mach policy
type MatchPolicy struct {
	Headers  map[string]map[string]string `yaml:"headers"`
	APIPaths map[string]string            `yaml:"apiPath"`
	Method   string                       `yaml:"method"`
}

//LimiterConfig is rate limiter policy
type LimiterConfig struct {
	Match string
	QPS   string
}
