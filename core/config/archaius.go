package config

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"sync"
)

var initArchaiusOnce = sync.Once{}

// InitArchaius is to initialize the archaius
func InitArchaius() error {
	var err error
	initArchaiusOnce.Do(func() {
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

		err = archaius.Init(archaius.WithRequiredFiles(essentialfiles), archaius.WithOptionalFiles(commonfiles))
	})

	return err
}
