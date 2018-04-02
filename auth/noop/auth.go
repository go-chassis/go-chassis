package noop

import "github.com/ServiceComb/go-chassis/auth"

// Auth has no operations
type Auth struct {
}

//CheckAuthorization always success
func (a *Auth) CheckAuthorization(check *auth.Check) *auth.CheckResult {
	return &auth.CheckResult{
		Err: nil,
	}
}

//GetAPICertification return no error
func (a *Auth) GetAPICertification(ak, sk, project string) (*auth.Cert, error) {
	return nil, nil
}
func newAuth(service, role string, props map[string]string) auth.Auth {
	return &Auth{}
}
func init() {
	auth.InstallPlugin("noop", newAuth)
}
