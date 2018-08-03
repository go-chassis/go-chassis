package control

//LoadBalancingConfig is a standardized model
type LoadBalancingConfig struct {
	Strategy     string
	Filters      []string
	RetryEnabled bool
	RetryOnSame  int
	RetryOnNext  int
	BackOffKind  string
	BackOffMin   int
	BackOffMax   int
}

//RateLimitingConfig is a standardized model
type RateLimitingConfig struct {
	Key     string
	Enabled bool
	Rate    int
}
