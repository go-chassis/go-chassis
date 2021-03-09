package cipher_test

import (
	"github.com/go-chassis/cari/security"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/security/cipher"
	_ "github.com/go-chassis/go-chassis/v2/security/cipher/plugins/aes"
	"github.com/stretchr/testify/assert"
	"testing"
)

//DefaultCipher is a struct
type DefaultCipher struct {
}

func init() {
	cipher.InstallCipherPlugin("default", new)
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
	cipher.InstallCipherPlugin("test", new)
	c, err := cipher.NewCipher("test")
	assert.NoError(t, err)
	r, _ := c.Encrypt("test")
	assert.Equal(t, "test", r)
	_, err = cipher.GetCipherNewFunc("asd")
	assert.Error(t, err)

	_, err = cipher.GetCipherNewFunc("aes")
	assert.NoError(t, err)
	t.Run("Init", func(t *testing.T) {
		archaius.Init(archaius.WithMemorySource())
		archaius.Set("servicecomb.cipher.plugin", "default")
		err := cipher.Init()
		assert.NoError(t, err)
		s, err := cipher.Decrypt("text")
		assert.NoError(t, err)
		assert.Equal(t, "text", s)
		s, err = cipher.Encrypt("text")
		assert.NoError(t, err)
		assert.Equal(t, "text", s)
	})
}
