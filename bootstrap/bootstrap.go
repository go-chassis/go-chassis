package bootstrap

import (
	"github.com/go-chassis/go-chassis/core/lager"
)

var bootstrapPlugins = make([]*PluginItem, 0)

//PluginItem include name and plugin implementation
type PluginItem struct {
	Name   string
	Plugin BootstrapPlugin
}

//BootstrapPlugin is a interface which declares Init method
type BootstrapPlugin interface {
	Init() error
}

// The BootstrapFunc type is an adapter to allow the use of ordinary functions as bootstrapPlugin.
type BootstrapFunc func() error

//Init is a method
func (b BootstrapFunc) Init() error {
	return b()
}

//InstallPlugin is a function which installs plugin,
// during initiating of go chassis, plugins will be executed
func InstallPlugin(name string, plugin BootstrapPlugin) {
	bootstrapPlugins = append(bootstrapPlugins, &PluginItem{
		Name:   name,
		Plugin: plugin,
	})
}

//Bootstrap will boot plugins in orders
func Bootstrap() {
	for _, bp := range bootstrapPlugins {
		lager.Logger.Info("Bootstrap " + bp.Name)
		if err := bp.Plugin.Init(); err != nil {
			lager.Logger.Errorf(err, "Failed to init %s.", bp.Name)
		}
	}
}
