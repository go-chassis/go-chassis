package config

// RouterConfig is the struct having info about route rule destinations, source templates
type RouterConfig struct {
	Destinations    map[string][]*RouteRule `yaml:"routeRule"`
	SourceTemplates map[string]*Match       `yaml:"sourceTemplate"`
}

// RouteRule is having route rule parameters
type RouteRule struct {
	Precedence int         `yaml:"precedence"`
	Routes     []*RouteTag `yaml:"route"`
	Match      Match       `yaml:"match"`
}

// RouteTag gives route tag information
type RouteTag struct {
	Tags   map[string]string `yaml:"tags"`
	Weight int               `yaml:"weight"`
}

// Match is checking source, source tags, and http headers
type Match struct {
	Refer       string                       `yaml:"refer"`
	Source      string                       `yaml:"source"`
	SourceTags  map[string]string            `yaml:"sourceTags"`
	HTTPHeaders map[string]map[string]string `yaml:"httpHeaders"`
	Headers     map[string]map[string]string `yaml:"headers"`
}

// DarkLaunchRule dark launch rule
type DarkLaunchRule struct {
	Type  string      `json:"policyType"` // RULE/RATE
	Items []*RuleItem `json:"ruleItems"`
}

// RuleItem rule item
type RuleItem struct {
	GroupName       string `json:"groupName"`
	GroupCondition  string `json:"groupCondition"`  // version=0.0.1
	PolicyCondition string `json:"policyCondition"` // 80/test!=2
}
