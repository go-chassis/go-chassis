package client

import (
	"github.com/go-mesh/openlogging"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var onceInit sync.Once
var onceMonitor sync.Once
var instance *AddressPool

const (
	available               string = "available"
	unavailable             string = "unavailable"
	defaultCheckSCIInterval        = 25 // default sc instance health check interval in second
)

// AddressPool registry address pool
type AddressPool struct {
	addressMap map[string]string
	status     map[string]string
	mutex      sync.RWMutex
}

// GetInstance Get registry pool instance
func GetInstance() *AddressPool {
	onceInit.Do(func() {
		instance = &AddressPool{
			addressMap: make(map[string]string),
			status:     make(map[string]string),
		}
	})
	return instance
}

// SetAddress set addresses to pool
func (p *AddressPool) SetAddress(addresses []string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	openlogging.Info("set service center endpoints " + strings.Join(addresses, ","))
	p.addressMap = make(map[string]string)
	for _, v := range addresses {
		p.status[v] = available
		p.addressMap[v] = v
	}
}

// GetAvailableAddress Get an available address from pool by roundrobin
func (p *AddressPool) GetAvailableAddress() string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	addrs := make([]string, 0)
	for _, v := range p.addressMap {
		if p.status[v] == available {
			addrs = append(addrs, v)
		}
	}

	next := RoundRobin(addrs)
	addr, err := next()
	if err != nil {
		openlogging.Warn("use default service center address")
		return DefaultAddr
	}
	return addr
}

func (p *AddressPool) checkConnectivity() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	timeOut := time.Duration(1) * time.Second
	for _, v := range p.addressMap {
		conn, err := net.DialTimeout("tcp", v, timeOut)
		if err != nil {
			openlogging.Error("can not connect to sc endpoint", openlogging.WithTags(openlogging.Tags{
				"err": err.Error(),
			}))
			p.status[v] = unavailable
		} else {
			p.status[v] = available
			conn.Close()
		}
	}
}

//Monitor monitor each service center network connectivity
func (p *AddressPool) Monitor() {
	onceMonitor.Do(func() {
		p.checkConnectivity()
		var interval time.Duration
		v, isExist := os.LookupEnv(EnvCheckSCIInterval)
		if !isExist {
			interval = defaultCheckSCIInterval
		} else {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				interval = defaultCheckSCIInterval
			} else {
				interval = time.Duration(i)
			}
		}
		ticker := time.NewTicker(interval * time.Second)
		quit := make(chan struct{})

		go func() {
			for {
				select {
				case <-ticker.C:
					p.checkConnectivity()
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	})
}
