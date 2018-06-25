package pilot

import (
	"os"
	"strings"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
)

const (
	// PODNAMESPACE means pod's namespace
	PODNAMESPACE = "POD_NAMESPACE"
	// DefaultSuffix means suffix of service key of api v1
	DefaultSuffix = "svc.cluster.local|http"
)

// Close : Close all connection.
func close(r *EnvoyDSClient) error {
	err := r.Close()
	if err != nil {
		lager.Logger.Errorf(err, "Conn close failed.")
		return err
	}
	lager.Logger.Debugf("Conn close success.")
	return nil
}

// filterInstances filter instances
func filterInstances(hs []*Host) []*registry.MicroServiceInstance {
	instances := make([]*registry.MicroServiceInstance, 0)
	for _, h := range hs {
		msi := ToMicroServiceInstance(h, nil)
		instances = append(instances, msi)
	}
	return instances
}

func pilotServiceKey(service string) string {
	ns := os.Getenv(PODNAMESPACE)
	if ns == "" {
		ns = "default"
	}
	return strings.Join([]string{service, ns, DefaultSuffix}, ".")
}

func pilotQueryKey(serviceKey string, tags registry.Tags) string {
	if tags == nil {
		return "/" + serviceKey
	}
	ss := make([]string, 0, len(tags))
	for k, v := range tags {
		if k == common.BuildinTagVersion && v == common.LatestVersion {
			continue
		}
		ss = append(ss, k+"="+v)
	}
	return "/" + serviceKey + "|" + strings.Join(ss, ",")
}

func pilotTags(labels []string, key string) map[string]string {
	ss := strings.Split(key, ":")
	if len(ss) != len(labels) {
		return nil
	}
	ret := make(map[string]string, len(labels))
	for i, t := range labels {
		ret[t] = ss[i]
	}
	return ret
}
