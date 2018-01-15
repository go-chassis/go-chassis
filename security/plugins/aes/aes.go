//+build linux

package aes

import (
	"os"

	"github.com/ServiceComb/go-chassis/core/goplugin"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/security"
)

const cipherPlugin = "cipher_plugin.so"

//Cipher interface declares Init(), Encrypyt(), Decrypyt() methods
type Cipher interface {
	Init()
	Encrypt(src string) (string, error)
	Decrypt(src string) (string, error)
}

// AESCipher is a cipher used in huawei
type AESCipher struct {
	gcryptoEngine Cipher
}

func init() {
	if v, exist := os.LookupEnv("CIPHER_ROOT"); exist {
		os.Setenv("PAAS_CRYPTO_PATH", v)
	}
	security.InstallCipherPlugin("aes", new)
}

func new() security.Cipher {
	p, err := goplugin.LoadPlugin(cipherPlugin)
	if err != nil {
		if os.IsNotExist(err) {
			lager.Logger.Errorf(nil, "%s not found", cipherPlugin)
		} else {
			lager.Logger.Errorf(err, "Load %s failed", cipherPlugin)
		}
		return nil
	}
	cipher, err := p.Lookup("Cipher")
	if err != nil {
		lager.Logger.Errorf(err, "Get init method error!")
		return nil
	}
	cipherInstance, ok := cipher.(Cipher)
	if !ok {
		lager.Logger.Infof("E: Expecting Cipher interface, but got something else.")
		return nil
	}
	cipherInstance.Init()
	return &AESCipher{
		gcryptoEngine: cipherInstance,
	}
}

//Encrypt is method used for encryption
func (ac *AESCipher) Encrypt(src string) (string, error) {
	return ac.gcryptoEngine.Encrypt(src)
}

//Decrypt is method used for decryption
func (ac *AESCipher) Decrypt(src string) (string, error) {
	return ac.gcryptoEngine.Decrypt(src)
}
