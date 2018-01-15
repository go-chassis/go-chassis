package iputil_test

import (
	"github.com/ServiceComb/go-chassis/util/iputil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHostName(t *testing.T) {
	hostname := iputil.GetHostName()
	assert.NotNil(t, hostname)
}
func TestGetLocapIp(t *testing.T) {
	localip := iputil.GetLocalIP()
	assert.NotNil(t, localip)
}
