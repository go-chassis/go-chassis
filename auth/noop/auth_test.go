package noop_test

import (
	"github.com/ServiceComb/go-chassis/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuth_CheckAuthorization(t *testing.T) {
	f := auth.GetPlugin("noop")
	a := f("a", "", nil)
	_, err := a.GetAPICertification("", "", "")
	assert.NoError(t, err)
	r := a.CheckAuthorization(nil)
	assert.NoError(t, r.Err)
}
