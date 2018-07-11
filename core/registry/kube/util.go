package kuberegistry

import (
	"os"
	"strings"

	"github.com/ServiceComb/go-chassis/core/common"
)

func splitServiceKey(key string) (name, namespace string) {
	sets := strings.Split(key, ".")
	if len(sets) >= 2 {
		return sets[0], sets[1]
	}

	ns := os.Getenv("POD_NAMESPACE")
	if ns == "" {
		ns = common.DefaultValue
	}
	if len(sets) == 1 {
		return sets[0], ns
	}
	return key, ns
}
