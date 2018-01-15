package loadbalance

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
)

type defaultSelector struct {
	opts Options
}

func init() {
	rand.Seed(time.Now().Unix())
}

func (r *defaultSelector) Init(opts ...Option) error {
	for _, o := range opts {
		o(&(r.opts))
	}
	return nil
}

func (r *defaultSelector) Options() Options {
	return r.opts
}

func (r *defaultSelector) Select(serviceName, version string, opts ...SelectOption) (Next, error) {
	sopts := SelectOptions{}
	for _, opt := range opts {
		opt(&sopts)
	}

	if sopts.Strategy == nil {
		sopts.Strategy = r.opts.Strategy
	}

	var isFilterExist = true
	for _, filter := range sopts.Filters {
		if filter == nil {
			isFilterExist = false
		}

	}

	// get the service
	if sopts.AppID == "" {
		sopts.AppID = config.GlobalDefinition.AppID
	}

	instances, err := r.opts.Registry.FindMicroServiceInstances(sopts.ConsumerID, sopts.AppID, serviceName, version)
	if err != nil {
		lbErr := LBError{err.Error()}
		lager.Logger.Errorf(lbErr, "Lb err")
		return nil, lbErr
	}

	// apply the filters
	if isFilterExist {
		for _, filter := range sopts.Filters {
			instances = filter(instances)
		}

	}

	// if there's nothing left, return
	if len(instances) == 0 {
		lbErr := LBError{fmt.Sprintf("No available instance, key: %s:%s:%s", sopts.AppID, serviceName, version)}
		lager.Logger.Error(lbErr.Error(), nil)
		return nil, lbErr
	}

	return sopts.Strategy(instances, sopts.Metadata), nil
}

func (r *defaultSelector) String() string {
	return "default"
}

func newDefaultSelector(opts ...Option) Selector {
	sopts := Options{
		Strategy: RoundRobin, //default
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	if sopts.Registry == nil {
		sopts.Registry = registry.RegistryService
	}

	lager.Logger.Debugf("Set default selector's registry: %s.", sopts.Registry)
	return &defaultSelector{
		opts: sopts,
	}
}
