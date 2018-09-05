package aes

import (
	"os"

	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/goplugin"
	"github.com/go-chassis/go-chassis/security"
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
	cipher, err := goplugin.LookUpSymbolFromPlugin(cipherPlugin, "Cipher")
	if err != nil {
		if os.IsNotExist(err) {
			lager.Logger.Errorf("%s not found", cipherPlugin)
		} else {
			lager.Logger.Errorf("Load %s failed, err [%s]", cipherPlugin, err.Error())
		}
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
