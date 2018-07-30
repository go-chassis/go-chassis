package pilot

// Service is the envoy service struct
type Service struct {
	ServiceKey string  `json:"service-key"`
	Hosts      []*Host `json:"hosts"`
}

// Hosts is the struct which contains a list of Host
type Hosts struct {
	Hosts []*Host `json:"hosts"`
}

// Host contains upstream ip address, port and tags
type Host struct {
	Address string `json:"ip_address"`
	Port    int    `json:"port"`
	Tags    *Tags  `json:"tags,omitempty"`
}

// Tags contains az, canary and weight
type Tags struct {
	AZ     string `json:"az,omitempty"`
	Canary bool   `json:"canary,omitempty"`

	// Weight is an integer in the range [1, 100] or empty
	Weight int `json:"load_balancing_weight,omitempty"`
}
