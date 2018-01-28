// Package loadbalance is a way to load balance service nodes
package selector

import (
	"github.com/ServiceComb/go-chassis/core/registry"
)

var (
	// ErrNoneAvailable is to represent load balance error
	ErrNoneAvailable = LBError{Message: "No available"}
)

// LBError load balance error
type LBError struct {
	Message string
}

// Error for to return load balance error message
func (e LBError) Error() string {
	return "lb: " + e.Message
}

// Selector builds on the registry as a mechanism to pick nodes
// and mark their status. This allows host pools and other things
// to be built using various algorithms.
type Selector interface {
	Init(opts ...Option) error
	Options() Options
	// Select returns a function which should return the next node
	Select(microserviceName, version string, opts ...SelectOption) (Next, error)
	// Name of the selector
	String() string
}

// Next is a function that returns the next node
// based on the selector's strategy
type Next func() (*registry.MicroServiceInstance, error)

// Filter is used to filter a service during the selection process
type Filter func([]*registry.MicroServiceInstance) []*registry.MicroServiceInstance

// Strategy is a selection strategy e.g random, round robin
type Strategy func([]*registry.MicroServiceInstance, interface{}) Next
