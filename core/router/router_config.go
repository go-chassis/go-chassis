package router

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/router/model"
	"github.com/ServiceComb/go-chassis/util/fileutil"
	"gopkg.in/yaml.v2"
)

// Init initialize router config
func Init() error {
	// init dests and templates
	routerConfigFromFile, err := getRouterConfigFromFile()
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		lager.Logger.Debugf("%s not exist", fileutil.Router)
	} else {
		if routerConfigFromFile.Destinations != nil {
			dests = routerConfigFromFile.Destinations
		}
		if routerConfigFromFile.SourceTemplates != nil {
			templates = routerConfigFromFile.SourceTemplates
		}
	}

	// the manager use dests to init, so must init after dests
	if err = initRouterManager(); err != nil {
		return err
	}

	if err = refresh(); err != nil {
		return err
	}
	lager.Logger.Info("Router init success")
	return nil
}

// refresh all the router config
func refresh() error {
	configs := routeRuleMgr.GetConfigurations()
	d := make(map[string][]*model.RouteRule)
	for k, v := range configs {
		rules, ok := v.([]*model.RouteRule)
		if !ok {
			err := fmt.Errorf("route rule type assertion fail, key: %s", k)
			return err
		}
		d[k] = rules
	}
	dests = d
	return nil
}

// get router config from router config file
func getRouterConfigFromFile() (*model.RouterConfig, error) {
	f := fileutil.GetRouter()
	if _, err := os.Stat(f); err != nil {
		return nil, err
	}
	routerConfig := &model.RouterConfig{}
	contents, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal([]byte(contents), routerConfig); err != nil {
		return nil, err
	}
	return routerConfig, nil
}
