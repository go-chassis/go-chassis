package registry

import (
	"errors"
	"fmt"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"strings"
	"time"
)

const (
	timeoutToPending = 1 * time.Second
	timeoutToPackage = 100 * time.Millisecond
	chanCapacity     = 1000
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
			if cr := <-r; cr.Err != nil {
				hc.removeFromCache(cr.Item)
				continue
			}
		}
	}
}

func (hc *HealthChecker) doCheck(i *WrapInstance) <-chan checkResult {
	cr := make(chan checkResult)
	go func() {
		// TODO call provider healthz api
		cr <- checkResult{Item: i, Err: errors.New("DELETE")}
		lager.Logger.Debugf("Health check %s failed", i.String())
	}()
	return cr
}

func (hc *HealthChecker) removeFromCache(i *WrapInstance) {
	key := i.ServiceKey()
	c, ok := MicroserviceInstanceCache.Get(key)
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
	MicroserviceInstanceCache.Set(key, is, 0)
	lager.Logger.Debugf("Health check: cached [%d] Instances of service [%s]", len(is), key)
}

// HealthCheck is the function adds the instance to HealthChecker
func HealthCheck(serviceKey string, instance *MicroServiceInstance) error {
	if !config.GetServiceDiscoveryHealthCheck() {
		return fmt.Errorf("Health check is disabled")
	}

	arr := strings.Split(serviceKey, ":")
	return defaultHealthChecker.Add(&WrapInstance{
		ServiceName: arr[0],
		Version:     arr[1],
		AppID:       arr[2],
		Instance:    instance,
	})
}

func RefreshCache(store map[string][]*MicroServiceInstance) {
	for key, v := range store {
		c, ok := MicroserviceInstanceCache.Get(key)
		if !ok {
			MicroserviceInstanceCache.Set(key, v, 0)
			continue
		}
		var (
			news   []*MicroServiceInstance
			lefts  []*MicroServiceInstance
			elders = make(map[string]*MicroServiceInstance)
			newers = make(map[string]*MicroServiceInstance)
		)

		for _, instance := range c.([]*MicroServiceInstance) {
			elders[instance.InstanceID] = instance
		}

		for _, instance := range v {
			newers[instance.InstanceID] = instance
		}

		for _, elder := range elders {
			if _, ok := newers[elder.InstanceID]; ok {
				lefts = append(lefts, elder)
				continue
			}
			if err := HealthCheck(key, elder); err == nil {
				lefts = append(lefts, elder)
			} // else remove the cache immediately if HC failed
		}

		for _, newer := range v {
			if _, ok := elders[newer.InstanceID]; ok {
				continue
			}
			news = append(news, newer)
		}

		lefts = append(lefts, news...)
		MicroserviceInstanceCache.Set(key, lefts, 0)
		lager.Logger.Debugf("Cached [%d] Instances of service [%s]", len(lefts), key)
	}
}
