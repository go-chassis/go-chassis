package model


type  EgressConfig struct {
	Destinations map[string][]*EgressRule `yaml:"egressRule"`
}

type EgressRule struct {
	Hosts  []string      `yaml:"hosts"`
	Ports  []*EgressPort `yaml:"ports"`

}


type EgressPort struct {
	Port     int32  `yaml:"port"`
	Protocol string `yaml:"protocol"`
}
