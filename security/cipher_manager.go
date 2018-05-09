package security

import (
	"fmt"
	"os"

	"github.com/ServiceComb/go-chassis/core/goplugin"
	"github.com/ServiceComb/go-chassis/core/lager"
)

const pluginSuffix = ".so"

//CipherPlugins is a map
var cipherPlugins map[string]func() Cipher

//InstallCipherPlugin is a function
func InstallCipherPlugin(name string, f func() Cipher) {
	cipherPlugins[name] = f
}

//GetCipherNewFunc is a function
func GetCipherNewFunc(name string) (func() Cipher, error) {
	if f, ok := cipherPlugins[name]; ok {
		return f, nil
	}
	lager.Logger.Debugf("try to load cipher [%s] from go plugin", name)
	f, err := loadCipherFromPlugin(name)
	if err == nil {
		cipherPlugins[name] = f
		return f, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	return nil, fmt.Errorf("unkown cipher plugin [%s]", name)
}

func loadCipherFromPlugin(name string) (func() Cipher, error) {
	p, err := goplugin.LoadPlugin(name + pluginSuffix)
	if err != nil {
		return nil, err
	}
	c, err := p.Lookup("Cipher")
	if err != nil {
		return nil, err
	}
	customCipher := c.(Cipher)
	f := func() Cipher {
		return customCipher
	}
	return f, nil
}

func init() {
	cipherPlugins = make(map[string]func() Cipher)
}
