package goplugin

import (
	"os"
	"plugin"

	"github.com/ServiceComb/go-chassis/util/fileutil"
)

// LookupPlugin lookup plugin
// Caller needs to determine itself whether the plugin file exists
func LookupPlugin(name string) (string, error) {
	var pluginPath string
	var err error
	// firstly search plugin in {ChassisHome}/lib
	pluginPath = fileutil.ChassisHomeDir() + "/lib/" + name
	if _, err = os.Stat(pluginPath); err == nil {
		return pluginPath, nil
	}
	if !os.IsNotExist(err) {
		return "", err
	}

	// secondly search plugin in /usr/lib
	pluginPath = "/usr/lib/" + name
	if _, err = os.Stat(pluginPath); err == nil {
		return pluginPath, nil
	}
	return "", err
}

// LoadPlugin load plugin
func LoadPlugin(name string) (*plugin.Plugin, error) {
	path, err := LookupPlugin(name)
	if err != nil {
		return nil, err
	}
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	return p, nil
}
