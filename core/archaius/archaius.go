package archaius

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
)

// Init is to initialize the archaius
func Init() error {
	essentialfiles := []string{
		fileutil.GlobalDefinition(),
		fileutil.GetMicroserviceDesc(),
	}
	commonfiles := []string{
		fileutil.HystrixDefinition(),
		fileutil.GetLoadBalancing(),
		fileutil.GetRateLimiting(),
		fileutil.GetTLS(),
		fileutil.GetMonitoring(),
		fileutil.GetAuth(),
		fileutil.GetTracing(),
	}

	err := archaius.Init(archaius.WithRequiredFiles(essentialfiles), archaius.WithOptionalFiles(commonfiles))
	return err
}
