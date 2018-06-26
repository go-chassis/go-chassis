package pilot

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/ServiceComb/go-archaius/core"
	cm "github.com/ServiceComb/go-archaius/core/config-manager"
	"github.com/ServiceComb/go-archaius/core/event-system"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/router"
	wp "github.com/ServiceComb/go-chassis/core/router/weightpool"
	"github.com/ServiceComb/go-chassis/pkg/istio/client"
	"github.com/ServiceComb/go-chassis/pkg/istio/util"
)

const routePilotSourceName = "RoutePilotSource"
const routePilotSourcePriority = 8

// DefaultPilotRefresh is default pilot refresh time
// TODO: use stream instead
var DefaultPilotRefresh = 10 * time.Second

var pilotfetcher core.ConfigMgr
var pilotChan = make(chan string, 10)

func setChanForPilot(k string) bool {
	select {
	case pilotChan <- k:
		return true
	default:
		return false
	}
}

// InitPilotFetcher init the config mgr and add several sources
func InitPilotFetcher(o router.Options) error {
	d := eventsystem.NewDispatcher()

	// register and init pilot fetcher
	d.RegisterListener(&pilotEventListener{}, ".*")
	pilotfetcher = cm.NewConfigurationManager(d)

	return addRoutePilotSource(o)
}

// addRoutePilotSource adds a config source to pilotfetcher
func addRoutePilotSource(o router.Options) error {
	if pilotfetcher == nil {
		return errors.New("pilotfetcher is nil, please init it first")
	}

	s, err := newPilotSource(o)
	if err != nil {
		return err
	}
	lager.Logger.Infof("New [%s] source success", s.GetSourceName())
	return pilotfetcher.AddSource(s, s.GetPriority())
}

// pilotSource keeps the route rule in istio
type pilotSource struct {
	refreshInverval time.Duration
	fetcher         client.PilotClient

	mu             sync.RWMutex
	pmu            sync.RWMutex
	Configurations map[string]interface{}
	PortToService  map[string]string
}

func newPilotSource(o router.Options) (*pilotSource, error) {
	grpcClient, err := client.NewGRPCPilotClient(o.ToPilotOptions())
	if err != nil {
		return nil, fmt.Errorf("connect to pilot failed: %v", err)
	}

	return &pilotSource{
		// TODO: read from config
		refreshInverval: DefaultPilotRefresh,
		Configurations:  map[string]interface{}{},
		PortToService:   map[string]string{},
		fetcher:         grpcClient,
	}, nil
}

func (r *pilotSource) GetSourceName() string { return routePilotSourceName }
func (r *pilotSource) GetPriority() int      { return routePilotSourcePriority }
func (r *pilotSource) Cleanup() error        { return nil }

func (r *pilotSource) AddDimensionInfo(d string) (map[string]string, error)           { return nil, nil }
func (r *pilotSource) GetConfigurationsByDI(d string) (map[string]interface{}, error) { return nil, nil }
func (r *pilotSource) GetConfigurationByKeyAndDimensionInfo(key, d string) (interface{}, error) {
	return nil, nil
}

func (r *pilotSource) GetConfigurations() (map[string]interface{}, error) {
	routerConfigs, err := r.getRouterConfigFromPilot()
	if err != nil {
		lager.Logger.Error("Get router config from pilot failed", err)
		return nil, err
	}
	d := make(map[string]interface{}, 0)
	for k, v := range routerConfigs.Destinations {
		d[k] = v
	}
	r.mu.Lock()
	r.Configurations = d
	r.mu.Unlock()
	return d, nil
}

func (r *pilotSource) GetConfigurationByKey(k string) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if value, ok := r.Configurations[k]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("not found %s", k)
}

// get router config from pilot
func (r *pilotSource) getRouterConfigFromPilot() (*model.RouterConfig, error) {
	routeRules := &model.RouterConfig{
		Destinations: map[string][]*model.RouteRule{},
	}

	sets := GetRouteRule()
	routeConfigs, err := r.fetcher.GetAllRouteConfigurations()
	if err != nil {
		return nil, err
	}
	for _, host := range routeConfigs.VirtualHosts {
		ss, port := util.ServiceAndPort(host.Name)
		if len(ss) == 0 {
			continue
		}

		if _, ok := sets[ss]; ok {
			routeRules.Destinations[ss] = VirtualHostsToRouteRule(&host)
			r.setPortForDestination(ss, port)
		}
	}
	return routeRules, nil
}

func (r *pilotSource) setPortForDestination(service, port string) {
	r.pmu.RLock()
	r.PortToService[port] = service
	r.pmu.RUnlock()
}

func (r *pilotSource) DynamicConfigHandler(callback core.DynamicConfigCallback) error {
	// Periodically refresh configurations
	ticker := time.NewTicker(r.refreshInverval)
	for {
		select {
		case <-pilotChan:
			data, err := r.GetConfigurations()
			if err != nil {
				lager.Logger.Error("pilot pull configuration error", err)
				continue
			}
			for k, d := range data {
				SetRouteRuleByKey(k, d.([]*model.RouteRule))
			}

		case <-ticker.C:
			data, err := r.refreshConfigurations()
			if err != nil {
				lager.Logger.Error("pilot refresh configuration error", err)
				continue
			}
			events, err := r.populateEvents(data)
			if err != nil {
				lager.Logger.Warnf("populate event error", err)
				return err
			}
			//Generate OnEvent Callback based on the events created
			lager.Logger.Debugf("event On receive %+v", events)
			for _, event := range events {
				callback.OnEvent(event)
			}
		}
	}
	return nil
}

func (r *pilotSource) refreshConfigurations() (map[string]interface{}, error) {
	data := make(map[string]interface{}, 0)

	sets := GetRouteRule()
	for port := range r.PortToService {
		routeConfigs, err := r.fetcher.GetRouteConfigurationsByPort(port)
		if err != nil {
			return nil, err
		}
		for _, host := range routeConfigs.VirtualHosts {
			ss, _ := util.ServiceAndPort(host.Name)
			if len(ss) == 0 {
				continue
			}
			if _, ok := sets[ss]; ok {
				data[ss] = VirtualHostsToRouteRule(&host)
			}
		}
	}
	return data, nil
}

func (r *pilotSource) populateEvents(updates map[string]interface{}) ([]*core.Event, error) {
	events := make([]*core.Event, 0)
	new := make(map[string]interface{})

	// generate create and update event
	r.mu.RLock()
	current := r.Configurations
	r.mu.RUnlock()

	for key, value := range updates {
		new[key] = value
		currentValue, ok := current[key]
		if !ok { // if new configuration introduced
			events = append(events, constructEvent(core.Create, key, value))
		} else if !reflect.DeepEqual(currentValue, value) {
			events = append(events, constructEvent(core.Update, key, value))
		}
	}
	// generate delete event
	for key, value := range current {
		_, ok := new[key]
		if !ok { // when old config not present in new config
			events = append(events, constructEvent(core.Delete, key, value))
		}
	}

	// update with latest config
	r.mu.Lock()
	r.Configurations = new
	r.mu.Unlock()
	return events, nil
}

func constructEvent(eventType string, key string, value interface{}) *core.Event {
	return &core.Event{
		EventType:   eventType,
		EventSource: routePilotSourceName,
		Value:       value,
		Key:         key,
	}
}

// pilotEventListener handle event dispatcher
type pilotEventListener struct{}

// update route rule of a service
func (r *pilotEventListener) Event(e *core.Event) {
	if e == nil {
		lager.Logger.Warn("pilot event pointer is nil", nil)
		return
	}

	v := pilotfetcher.GetConfigurationsByKey(e.Key)
	if v == nil {
		DeleteRouteRuleByKey(e.Key)
		lager.Logger.Infof("[%s] route rule of piot is removed", e.Key)
		return
	}
	routeRules, ok := v.([]*model.RouteRule)
	if !ok {
		lager.Logger.Error("value of pilot is not type []*RouteRule", nil)
		return
	}

	if router.ValidateRule(map[string][]*model.RouteRule{e.Key: routeRules}) {
		SetRouteRuleByKey(e.Key, routeRules)
		wp.GetPool().Reset(e.Key)
		lager.Logger.Infof("Update [%s] route rule of pilot success", e.Key)
	}
}
