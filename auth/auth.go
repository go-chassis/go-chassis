package auth

var authPlugin = make(map[string]func(role, service string, props map[string]string) Auth)

//InstallPlugin install auth plugin
func InstallPlugin(name string, f func(role, service string, props map[string]string) Auth) {
	authPlugin[name] = f
}

//GetPlugin return plugin
func GetPlugin(name string) func(role, service string, props map[string]string) Auth {
	return authPlugin[name]
}

// Check includes information to be checked by auth service
type Check struct {
	TargetService           string
	TargetMethod            string
	TargetServiceProperties map[string]string
}

//CheckResult is returned by auth service
type CheckResult struct {
	Message string
	Err     error
}

// Cert is certification
type Cert struct {
	Project    string
	AK         string
	ShaAKAndSK string
}

// Auth is for 2 cased
// 1. check source service is authorized or not when try to request micro service
// 2. Get API service certification from any system
type Auth interface {
	CheckAuthorization(check *Check) *CheckResult
	GetAPICertification(ak, sk, project string) (*Cert, error)
}
