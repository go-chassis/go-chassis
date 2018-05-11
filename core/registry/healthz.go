package registry

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	timeoutToPending = 1 * time.Second
	timeoutToPackage = 100 * time.Millisecond
	chanCapacity     = 1000
)

var defaultHealthChecker = &HealthChecker{}

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
	for {
		pack := make(map[string]*WrapInstance)
		select {
		case i, ok := <-hc.pendingCh:
			if !ok {
				// chan closed
				return
			}
			pack[i.String()] = i
		case <-time.After(timeoutToPackage):
			hc.delCh <- pack
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
}

// HealthCheck is the function adds the instance to HealthChecker
func HealthCheck(serviceKey string, instance *MicroServiceInstance) error {
	arr := strings.Split(serviceKey, ":")
	return defaultHealthChecker.Add(&WrapInstance{
		ServiceName: arr[0],
		Version:     arr[1],
		AppID:       arr[2],
		Instance:    instance,
	})
}
