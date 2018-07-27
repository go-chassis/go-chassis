package auth

import (
	"errors"
	"fmt"
	"github.com/go-chassis/auth"
	"github.com/go-chassis/go-chassis/core/archaius"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/goplugin"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/security"
	"github.com/go-chassis/http-client"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

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
	projectFromEnv := os.Getenv(paasProjectNameEnv)
	httpclient.SignRequest = func(r *http.Request) error {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, vs := range genAuthHeaders() {
			for _, v := range vs {
				r.Header.Add(k, v)
			}
		}
		if projectFromEnv != "" {
			r.Header.Set(auth.HeaderServiceProject, projectFromEnv)
		}
		return nil
	}
	return nil
}

func getAkskCustomCipher(name string) (security.Cipher, error) {
	f, err := security.GetCipherNewFunc(name)
	if err != nil {
		return nil, err
	}
	cipherPlugin := f()
	if cipherPlugin == nil {
		return nil, fmt.Errorf("cipher plugin [%s] invalid", name)
	}
	return cipherPlugin, nil
}

func getProjectFromURI(rawurl string) (string, error) {
	errGetProjectFailed := errors.New("get project from CSE uri failed")
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

	// 1, use project of env PAAS_PROJECT_NAME
	// 2, use project in the credential config
	// 3, use project in cse uri contain
	// 4, use project "default"
	if v := os.Getenv(paasProjectNameEnv); v != "" {
		c.Project = v
	}
	if c.Project == "" {
		project, err := getProjectFromURI(config.GetRegistratorAddress())
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
	if c.Project == "" {
		lager.Logger.Debug("Huawei Cloud project is empty")
	} else {
		lager.Logger.Debugf("Huawei Cloud project: %s", c.Project)
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
			return fmt.Errorf("decrypt sk failed %v", err)
		}
		plainSk = res
	}

	httpclient.SignRequest, err = auth.UseAKSKAuth(c.AccessKey, plainSk, c.Project)
	if err != nil {
		return err
	}
	return nil
}

// Init initializes auth module
func Init() error {
	err := loadAkskAuth()
	if err == nil {
		lager.Logger.Warn("Huawei Cloud auth mode: ak/sk", nil)
		return nil
	}
	if !isAuthConfNotExist(err) {
		lager.Logger.Error("Load ak/sk failed", err)
		return err
	}

	err = loadPaasAuth()
	if err == nil {
		lager.Logger.Warn("Huawei Cloud auth mode: token", nil)
		return nil
	}
	if !isAuthConfNotExist(err) {
		lager.Logger.Error("Load paas auth failed", err)
		return err
	}
	lager.Logger.Debug("No authentication for Huawei Cloud")
	return nil
}
