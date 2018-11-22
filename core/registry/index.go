package registry

import (
	"github.com/go-chassis/go-chassis/core/common"
)

func (m *MicroServiceInstance) appID() string   { return m.Metadata[common.BuildinTagApp] }
func (m *MicroServiceInstance) version() string { return m.Metadata[common.BuildinTagVersion] }

// Has return whether microservice has tags
func (m *MicroServiceInstance) Has(tags map[string]string) bool {
	for k, v := range tags {
		if mt, ok := m.Metadata[k]; !ok || mt != v {
			return false
		}
	}
	return true
}

// WithAppID add app tag for microservice instance
func (m *MicroServiceInstance) WithAppID(v string) *MicroServiceInstance {
	m.Metadata[common.BuildinTagApp] = v
	return m
}
