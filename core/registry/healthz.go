package registry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/healthz/client"
)

const (
	timeoutToPending     = 1 * time.Second
	timeoutToPackage     = 100 * time.Millisecond
	timeoutToHealthCheck = 5 * time.Second
	chanCapacity         = 1000
)

var defaultHealthChecker = &HealthChecker{}

func init() {
	defaultHealthChecker.Run()
}

// WrapInstance is the struct defines an instance object with appID/serviceName/version
type WrapInstance struct {
	AppID       string
	ServiceName string
	Version     string
	Instance    *MicroServiceInstance
}

// String is the method returns the string type current instance's key value
func (i *WrapInstance) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", i.ServiceName, i.Version, i.AppID, i.Instance.InstanceID)
}

// ServiceKey is the method returns the string type current instance's service key value
func (i *WrapInstance) ServiceKey() string {
	return fmt.Sprintf("%s:%s:%s", i.ServiceName, i.Version, i.AppID)
}

// checkResult is the struct defines the result from health check
type checkResult struct {
	Item *WrapInstance
	Err  error
}

// HealthChecker is the struct judges the instance health in the removing cache
type HealthChecker struct {
	pendingCh chan *WrapInstance
	delCh     chan map[string]*WrapInstance
}

// Run is the method initializes and starts the health check process
func (hc *HealthChecker) Run() {
	hc.pendingCh = make(chan *WrapInstance, chanCapacity)
	hc.delCh = make(chan map[string]*WrapInstance, chanCapacity)
	go hc.wait()
	go hc.check()
}

// Add is the method adds a key of the instance cache into pending chan
func (hc *HealthChecker) Add(i *WrapInstance) error {
	select {
	case hc.pendingCh <- i:
	case <-time.After(timeoutToPending):
		return errors.New("Health checker is too busy")
	}
	return nil
}

func (hc *HealthChecker) wait() {
	pack := make(map[string]*WrapInstance)
	for {
		select {
		case i, ok := <-hc.pendingCh:
			if !ok {
				// chan closed
				return
			}
			pack[i.String()] = i
		case <-time.After(timeoutToPackage):
			if len(pack) > 0 {
				hc.delCh <- pack
				pack = make(map[string]*WrapInstance)
			}
		}
	}
}

func (hc *HealthChecker) check() {
	for pack := range hc.delCh {
		var rs []<-chan checkResult
		for _, v := range pack {
			rs = append(rs, hc.doCheck(v))
		}
		for _, r := range rs {
			cr := <-r
			if cr.Err != nil {
				lager.Logger.Debugf("Health check instance %s failed, %s",
					cr.Item.ServiceKey(), cr.Err)
				hc.removeFromCache(cr.Item)
				continue
			}
			lager.Logger.Debugf("Health check instance %s %s is still alive, keep it in cache",
				cr.Item.ServiceKey(), cr.Item.Instance.EndpointsMap)
		}
	}
}

func (hc *HealthChecker) doCheck(i *WrapInstance) <-chan checkResult {
	cr := make(chan checkResult)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutToHealthCheck)
		r := checkResult{Item: i, Err: nil}
		defer func() {
			cancel()
			cr <- r
		}()
		req := client.Reply{
			AppId:       i.AppID,
			ServiceName: i.ServiceName,
			Version:     i.Version,
		}

		for protocol, ep := range i.Instance.EndpointsMap {
			r.Err = client.Test(ctx, protocol, ep, req)
			return
		}
	}()
	return cr
}

func (hc *HealthChecker) removeFromCache(i *WrapInstance) {
	c, ok := MicroserviceInstanceIndex.Get(i.ServiceName, nil)
	if !ok {
		return
	}
	var is []*MicroServiceInstance
	for _, inst := range c.([]*MicroServiceInstance) {
		if inst.InstanceID == i.Instance.InstanceID {
			continue
		}
		is = append(is, inst)
	}
	MicroserviceInstanceIndex.Set(i.ServiceName, is)
	lager.Logger.Debugf("Health check: cached [%d] Instances of service [%s]", len(is), i.ServiceName)
}

// HealthCheck is the function adds the instance to HealthChecker
func HealthCheck(service, version, appID string, instance *MicroServiceInstance) error {
	if !config.GetServiceDiscoveryHealthCheck() {
		return fmt.Errorf("Health check is disabled")
	}

	return defaultHealthChecker.Add(&WrapInstance{
		ServiceName: service,
		Version:     version,
		AppID:       appID,
		Instance:    instance,
	})
}

// RefreshCache is the function to filter changes between new pulling instances and cache
func RefreshCache(service string, store []*MicroServiceInstance) {
	var (
		arrUps []*MicroServiceInstance
		downs  = make(map[string]*MicroServiceInstance)
		ups    = make(map[string]*MicroServiceInstance)
	)

	for _, ins := range store {
		if ins.Status != common.DefaultStatus {
			downs[ins.InstanceID] = ins
			lager.Logger.Debugf("do not cache the instance in '%s' status, instanceId = %s/%s",
				ins.Status, ins.ServiceID, ins.InstanceID)
			continue
		}
		arrUps = append(arrUps, ins)
		ups[ins.InstanceID] = ins
	}

	c, ok := MicroserviceInstanceIndex.Get(service, nil)
	if !ok || c == nil || c.([]*MicroServiceInstance) == nil {
		// if full new instances or at less one instance, then refresh cache immediately
		MicroserviceInstanceIndex.Set(service, arrUps)
		return
	}

	var (
		saves []*MicroServiceInstance
		lefts []*MicroServiceInstance
		exps  = make(map[string]*MicroServiceInstance)
	)
	for _, instance := range c.([]*MicroServiceInstance) {
		exps[instance.InstanceID] = instance
	}

	for _, elder := range exps {
		// case: keep still alive instances
		if _, ok := ups[elder.InstanceID]; ok {
			lefts = append(lefts, elder)
			continue
		}
		// case: remove instances with the non-up status
		if _, ok := downs[elder.InstanceID]; ok {
			continue
		}
		// case: keep instances returned HC ok
		if err := HealthCheck(service, elder.version(), elder.appID(), elder); err == nil {
			lefts = append(lefts, elder)
		}
	}

	for _, up := range ups {
		if _, ok := exps[up.InstanceID]; ok {
			continue
		}
		// case: add new come in instances
		saves = append(saves, up)
	}

	lefts = append(lefts, saves...)
	MicroserviceInstanceIndex.Set(service, lefts)
	lager.Logger.Debugf("Cached [%d] Instances of service [%s]", len(lefts), service)
}
