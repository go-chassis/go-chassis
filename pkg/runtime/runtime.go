package runtime

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"os"
)

//HostName is the host name of service host
var HostName string

// Init runtime information
func Init() error {
	var err error
	service := config.MicroserviceDefinition
	HostName = service.ServiceDescription.Hostname
	if HostName == "" {
		HostName, err = os.Hostname()
		if err != nil {
			lager.Logger.Error("Get hostname failed.", err)
			return err
		}
	}
	lager.Logger.Info("Host name is " + HostName)
	return nil
}
