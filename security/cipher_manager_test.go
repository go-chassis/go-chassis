package security_test

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/security"
	_ "github.com/ServiceComb/go-chassis/security/plugins/aes"
	"github.com/stretchr/testify/assert"
	"testing"
)

//DefaultCipher is a struct
type DefaultCipher struct {
}

func init() {
	security.InstallCipherPlugin("default", new)
}
func new() security.Cipher {

	return &DefaultCipher{}
}

//Encrypt is method used for encryption
func (c *DefaultCipher) Encrypt(src string) (string, error) {
	return src, nil
}

//Decrypt is method used for decryption
func (c *DefaultCipher) Decrypt(src string) (string, error) {
	return src, nil
}

func TestInstallCipherPlugin(t *testing.T) {
	lager.Initialize("", "DEBUG", "", "size", true, 1, 10, 7)

	security.InstallCipherPlugin("test", new)
	f, err := security.GetCipherNewFunc("test")
	assert.NoError(t, err)
	c := f()
	r, _ := c.Encrypt("test")
	assert.Equal(t, "test", r)
	_, err = security.GetCipherNewFunc("asd")
	assert.Error(t, err)

	_, err = security.GetCipherNewFunc("aes")
	assert.NoError(t, err)
}
