package auth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ServiceComb/auth"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/goplugin"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/security"
	yaml "gopkg.in/yaml.v2"
)

const (
	paasAuthPlugin  = "paas_auth.so"
	keyAK           = "cse.credentials.accessKey"
	keySK           = "cse.credentials.secretKey"
	keyProject      = "cse.credentials.project"
	cipherRootEnv   = "CIPHER_ROOT"
	keytoolAkskFile = "certificate.yaml"
	keytoolCipher   = "security"
)

var errAuthConfNotExist = errors.New("Auth config is not exist")

func isAuthConfNotExist(e error) bool {
	return e == errAuthConfNotExist
}

// loadAkskAuth gets the Authentication Mode ak/sk, token and forms required Auth Headers
func loadPaasAuth() error {
	p, err := goplugin.LoadPlugin(paasAuthPlugin)
	if err != nil {
		if os.IsNotExist(err) {
			return errAuthConfNotExist
		}
		return err
	}

	f, err := p.Lookup("GenAuthHeaders")
	if err != nil {
		return err
	}

	genAuthHeaders := f.(func() http.Header)
	authFunc := func(r *http.Request) error {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range genAuthHeaders() {
			r.Header[k] = v
		}
		return nil
	}
	auth.SetAuthFunc(authFunc)
	return nil
}

func getAkskCustomCipher(name string) (security.Cipher, error) {
	f, err := security.GetCipherNewFunc(name)
	if err != nil {
		return nil, err
	}
	cipherPlugin := f()
	if cipherPlugin == nil {
		return nil, fmt.Errorf("Cipher plugin [%s] invalid", name)
	}
	return cipherPlugin, nil
}

func getProjectFromURI(rawurl string) (string, error) {
	errGetProjectFailed := errors.New("Get project from CSE uri failed")
	// rawurl: https://cse.cn-north-1.myhwclouds.com:443
	if rawurl == "" {
		return "", fmt.Errorf("%v, CSE uri empty", errGetProjectFailed)
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("%v, %v", errGetProjectFailed, err)
	}
	parts := strings.Split(u.Host, ".")
	if len(parts) != 4 {
		lager.Logger.Info("CSE uri contains no project")
		return "", nil
	}
	return parts[1], nil
}

func getAkskConfig() (*model.CredentialStruct, error) {
	// 1, if env CIPHER_ROOT exists, read ${CIPHER_ROOT}/certificate.yaml
	// 2, if env CIPHER_ROOT not exists, read chassis config
	var akskFile string
	if v, exist := os.LookupEnv(cipherRootEnv); exist {
		p := filepath.Join(v, keytoolAkskFile)
		if _, err := os.Stat(p); err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
		} else {
			akskFile = p
		}
	}

	c := &model.CredentialStruct{}
	if akskFile == "" {
		c.AccessKey = archaius.GetString(keyAK, "")
		c.SecretKey = archaius.GetString(keySK, "")
		c.Project = archaius.GetString(keyProject, "")
		c.AkskCustomCipher = archaius.GetString(common.AKSKCustomCipher, "")
	} else {
		yamlContent, err := ioutil.ReadFile(akskFile)
		if err != nil {
			return nil, err
		}
		globalConf := &model.GlobalCfg{}
		err = yaml.Unmarshal(yamlContent, globalConf)
		if err != nil {
			return nil, err
		}
		c = &(globalConf.Cse.Credentials)
	}
	if c.AccessKey == "" && c.SecretKey == "" {
		return nil, errAuthConfNotExist
	}
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, errors.New("One of ak and sk is empty")
	}

	// 1, use project in the credential config
	// 2, use project in cse uri contain
	// 3, use project "default"
	if c.Project == "" {
		project, err := getProjectFromURI(config.GlobalDefinition.Cse.Service.Registry.Address)
		if err != nil {
			return nil, err
		}
		if project != "" {
			c.Project = project
		} else {
			c.Project = common.DefaultValue
		}
	}
	return c, nil
}

// loadAkskAuth gets the Authentication Mode ak/sk
func loadAkskAuth() error {
	c, err := getAkskConfig()
	if err != nil {
		return err
	}

	plainSk := c.SecretKey
	cipher := c.AkskCustomCipher
	if cipher != "" {
		if cipher == keytoolCipher {
			lager.Logger.Infof("Use cipher plugin [aes] as plugin [%s]", cipher)
			cipher = "aes"
		}
		cipherPlugin, err := getAkskCustomCipher(cipher)
		if err != nil {
			return err
		}
		res, err := cipherPlugin.Decrypt(c.SecretKey)
		if err != nil {
			return fmt.Errorf("Decrypt sk failed %v", err)
		}
		plainSk = res
	}

	err = auth.UseAKSKAuth(c.AccessKey, plainSk, c.Project)
	if err != nil {
		return err
	}
	return nil
}

// Init initializes auth module
func Init() {
	err := loadAkskAuth()
	if err == nil {
		lager.Logger.Warn("Huawei Cloud auth mode: ak/sk", nil)
		return
	}
	if !isAuthConfNotExist(err) {
		lager.Logger.Error("Load ak/sk failed", err)
		return
	}

	err = loadPaasAuth()
	if err == nil {
		lager.Logger.Warn("Huawei Cloud auth mode: token", nil)
		return
	}
	if !isAuthConfNotExist(err) {
		lager.Logger.Error("Load paas auth failed", err)
		return
	}
	lager.Logger.Debug("No authentication for Huawei Cloud")
}
