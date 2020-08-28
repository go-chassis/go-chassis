package bootstrap

import (
	"fmt"
	"github.com/go-chassis/openlog"
)

var bootstrapPlugins = make([]*PluginItem, 0)

//PluginItem include name and plugin implementation
type PluginItem struct {
	Name   string
	Plugin Plugin
}

//Plugin is a interface which declares Init method
type Plugin interface {
	Init() error
}

// Func The Func type is an adapter to allow the use of ordinary functions as bootstrapPlugin.
type Func func() error

//Init is a method
func (b Func) Init() error {
	return b()
}

//InstallPlugin is a function which installs plugin,
// during initiating of go chassis, plugins will be executed
func InstallPlugin(name string, plugin Plugin) {
	bootstrapPlugins = append(bootstrapPlugins, &PluginItem{
		Name:   name,
		Plugin: plugin,
	})
}

//Bootstrap will boot plugins in orders
func Bootstrap() {
	for _, bp := range bootstrapPlugins {
		openlog.Info("Bootstrap " + bp.Name)
		if err := bp.Plugin.Init(); err != nil {
			openlog.Error(fmt.Sprintf("Failed to init %s. error [%s]", bp.Name, err.Error()))
		}
	}
}
