package pilot

// envoy service struct
type service struct {
	ServiceKey string  `json:"service-key"`
	Hosts      []*host `json:"hosts"`
}

// a list of hosts that make up the service
type hosts struct {
	Hosts []*host `json:"hosts"`
}

// host contains upstream ip address, port and tags
type host struct {
	Address string `json:"ip_address"`
	Port    int    `json:"port"`
	Tags    *tags  `json:"tags,omitempty"`
}

// optional tags per host
type tags struct {
	AZ     string `json:"az,omitempty"`
	Canary bool   `json:"canary,omitempty"`

	// Weight is an integer in the range [1, 100] or empty
	Weight int `json:"load_balancing_weight,omitempty"`
}
