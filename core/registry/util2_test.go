package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_fillUnspecifiedIp(t *testing.T) {
	host := "0.0.0.0"
	ipaddr, err := fillUnspecifiedIP(host)
	assert.NoError(t, err)
	assert.NotEmpty(t, ipaddr)
	assert.NotEqual(t, ipaddr, host)

	//host = "::"	// IPv6 address needed for local NIC
	//ipaddr, err = fillUnspecifiedIP(host)
	//assert.NoError(t, err)
	//assert.NotEmpty(t, ipaddr)
	//assert.NotEqual(t, ipaddr, host)

	host = "114.116.58.51"
	ipaddr, err = fillUnspecifiedIP(host)
	assert.NoError(t, err)
	assert.Equal(t, host, ipaddr)

	host = "fe80::c706:e006:d53e:f9fb"
	ipaddr, err = fillUnspecifiedIP(host)
	assert.NoError(t, err)
	assert.Equal(t, host, ipaddr)

	host = "abc"
	ipaddr, err = fillUnspecifiedIP(host)
	assert.Equal(t, "", ipaddr)
}
