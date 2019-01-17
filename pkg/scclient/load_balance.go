package client

import (
	"errors"
	"math/rand"
	"sync"
)

// ErrNoneAvailable create a new error with Message No available
var ErrNoneAvailable = errors.New("No available")

// Next gives the next object in the list
type Next func() (string, error)

var i = rand.Int()

// RoundRobin Gives the next object in sequence
func RoundRobin(eps []string) Next {
	var mtx sync.Mutex
	return func() (string, error) {
		if len(eps) == 0 {
			return "", ErrNoneAvailable
		}
		mtx.Lock()
		node := eps[i%len(eps)]
		i++
		mtx.Unlock()
		return node, nil
	}
}
