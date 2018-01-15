package loadbalance_test

import (
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"testing"
	"time"
)

func TestSessionStorage(t *testing.T) {
	loadbalance.Save("abc", "abc", time.Second)
}
