package cipher

import (
	"errors"
	"fmt"
	"github.com/go-chassis/go-archaius"
	"os"

	"github.com/go-chassis/cari/security"
	"github.com/go-chassis/go-chassis/v2/pkg/goplugin"
	"github.com/go-chassis/openlog"
)

const pluginSuffix = ".so"

// CipherPlugins is a map
var cipherPlugins = make(map[string]func() security.Cipher)

// InstallCipherPlugin is a function
func InstallCipherPlugin(name string, f func() security.Cipher) {
	cipherPlugins[name] = f
}

// NewCipher create and return a cipher
func NewCipher(name string) (security.Cipher, error) {
	f, err := GetCipherNewFunc(name)
	if err != nil {
		return nil, err
	}
	return f(), nil
}

// GetCipherNewFunc return a function which is able to create a cipher
func GetCipherNewFunc(name string) (func() security.Cipher, error) {
	if f, ok := cipherPlugins[name]; ok {
		return f, nil
	}
	openlog.Debug(fmt.Sprintf("try to load cipher [%s] from go plugin", name))
	f, err := loadCipherFromPlugin(name)
	if err == nil {
		cipherPlugins[name] = f
		return f, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	return nil, fmt.Errorf("unknown cipher plugin [%s]", name)
}

func loadCipherFromPlugin(name string) (func() security.Cipher, error) {
	c, err := goplugin.LookUpSymbolFromPlugin(name+pluginSuffix, "Cipher")
	if err != nil {
		return nil, err
	}
	customCipher, ok := c.(security.Cipher)
	if !ok {
		return nil, errors.New("symbol from plugin is not type Cipher")
	}
	f := func() security.Cipher {
		return customCipher
	}
	return f, nil
}

func Init() error {
	pn := archaius.GetString("servicecomb.cipher.plugin", "default")
	c, err := NewCipher(pn)
	if err != nil {
		return fmt.Errorf("init cipher [%s] failed, %w", pn, err)
	}
	defaultCipher = c
	return nil
}
