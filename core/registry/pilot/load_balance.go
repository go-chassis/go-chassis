package pilot

import (
	"errors"
	"math/rand"
	"sync/atomic"
)

// ErrNoneAvailable create a new error with Message No available
var ErrNoneAvailable = errors.New("No available")

// Next gives the next object in the list
type Next func() (string, error)

var i = int64(rand.Int())

// RoundRobin Gives the next object in sequence
func RoundRobin(eps []string) Next {
	return func() (string, error) {
		if len(eps) == 0 {
			return "", ErrNoneAvailable
		}
		node := eps[int(atomic.LoadInt64(&i))%len(eps)]
		atomic.AddInt64(&i, 1)
		return node, nil
	}
}
