package registry

import (
	"fmt"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	client "github.com/ServiceComb/go-sc-client"
	"github.com/ServiceComb/go-sc-client/model"
	pb "github.com/ServiceComb/go-sc-client/model"
	"os"
	"sync"
	"time"
)

// DefaultRetryTime default retry time
const DefaultRetryTime = 10 * time.Second

// HeartbeatTask heart beat task struct
type HeartbeatTask struct {
	ServiceID  string
	InstanceID string
	Time       time.Time
	Running    bool
}

// HeartbeatService heartbeat service
type HeartbeatService struct {
	Client    *client.RegistryClient
	instances map[string]*HeartbeatTask
	shutdown  bool
	mux       sync.Mutex
}

// Start start the heartbeat system
func (s *HeartbeatService) Start() {
	s.shutdown = false
	defer s.Stop()

	s.run()
}

// Stop stop the heartbeat system
func (s *HeartbeatService) Stop() {
	s.shutdown = true
}

// AddTask add new micro-service instance to the heartbeat system
func (s *HeartbeatService) AddTask(microServiceID, microServiceInstanceID string) {
	key := fmt.Sprintf("%s/%s", microServiceID, microServiceInstanceID)
	lager.Logger.Infof("Add HB task, task:%s", key)
	s.mux.Lock()
	if _, ok := s.instances[key]; !ok {
		s.instances[key] = &HeartbeatTask{
			ServiceID:  microServiceID,
			InstanceID: microServiceInstanceID,
			Time:       time.Now(),
		}
	}
	s.mux.Unlock()
}

// RemoveTask remove micro-service instance from the heartbeat system
func (s *HeartbeatService) RemoveTask(microServiceID, microServiceInstanceID string) {
	key := fmt.Sprintf("%s/%s", microServiceID, microServiceInstanceID)
	s.mux.Lock()
	delete(s.instances, key)
	s.mux.Unlock()
}

// RefreshTask refresh heartbeat for micro-service instance
func (s *HeartbeatService) RefreshTask(microServiceID, microServiceInstanceID string) {
	key := fmt.Sprintf("%s/%s", microServiceID, microServiceInstanceID)
	s.mux.Lock()
	if _, ok := s.instances[key]; ok {
		s.instances[key].Time = time.Now()
	}
	s.mux.Unlock()
}

// toggleTask toggle task
func (s *HeartbeatService) toggleTask(microServiceID, microServiceInstanceID string, running bool) {
	key := fmt.Sprintf("%s/%s", microServiceID, microServiceInstanceID)
	s.mux.Lock()
	if _, ok := s.instances[key]; ok {
		s.instances[key].Running = running
	}
	s.mux.Unlock()
}

// DoHeartBeat do heartbeat for each instance
func (s *HeartbeatService) DoHeartBeat(microServiceID, microServiceInstanceID string) {
	s.toggleTask(microServiceID, microServiceInstanceID, true)
	_, err := RegistryService.Heartbeat(microServiceID, microServiceInstanceID)
	if err != nil {
		lager.Logger.Errorf(err, "Run Heartbeat fail")
		s.RemoveTask(microServiceID, microServiceInstanceID)
		s.RetryRegister(microServiceID)
	}
	s.RefreshTask(microServiceID, microServiceInstanceID)
	s.toggleTask(microServiceID, microServiceInstanceID, false)
}

// run runs the heartbeat system
func (s *HeartbeatService) run() {
	for !s.shutdown {
		s.mux.Lock()
		endTime := time.Now()
		for _, v := range s.instances {
			if v.Running {
				continue
			}
			if endTime.Sub(v.Time) >= pb.DefaultLeaseRenewalInterval*time.Second {
				go s.DoHeartBeat(v.ServiceID, v.InstanceID)
			}
		}
		s.mux.Unlock()
		time.Sleep(time.Second)
	}
}

// RetryRegister retrying to register micro-service, and instance
func (s *HeartbeatService) RetryRegister(sid string) error {
	for {
		time.Sleep(DefaultRetryTime)
		lager.Logger.Infof("Try to re-register self")
		_, err := RegistryService.GetAllMicroServices()
		if err != nil {
			continue
		}
		if _, e := RegistryService.GetMicroService(sid); e != nil {
			err = s.ReRegisterSelfMSandMSI()
		} else {
			err = reRegisterSelfMSI(sid)
		}
		if err == nil {
			break
		}
	}
	lager.Logger.Warn("Re-register self success", nil)
	return nil
}

// ReRegisterSelfMSandMSI 重新注册微服务和实例
func (s *HeartbeatService) ReRegisterSelfMSandMSI() error {
	err := RegisterMicroservice()
	if err != nil {
		lager.Logger.Errorf(err, "The reRegisterSelfMSandMSI() startMicroservice failed.")
		return err
	}

	err = RegisterMicroserviceInstances()
	if err != nil {
		lager.Logger.Errorf(err, "The reRegisterSelfMSandMSI() startInstances failed.")
		return err
	}
	return nil
}

// reRegisterSelfMSI 只重新注册实例
func reRegisterSelfMSI(sid string) error {
	hostname, err := os.Hostname()
	if err != nil {
		lager.Logger.Errorf(err, "Get HostName failed, hostname:%s", hostname)
		return err
	}
	stage := config.Stage
	eps := MakeEndpointMap(config.GlobalDefinition.Cse.Protocols)
	if InstanceEndpoints != nil {
		eps = InstanceEndpoints
	}
	microServiceInstance := &MicroServiceInstance{
		EndpointsMap: eps,
		HostName:     hostname,
		Status:       model.MSInstanceUP,
		Environment:  stage,
	}
	instanceID, err := RegistryService.RegisterServiceInstance(sid, microServiceInstance)
	if err != nil {
		lager.Logger.Errorf(err, "RegisterInstance failed.")
		return err
	}

	value, ok := SelfInstancesCache.Get(microServiceInstance.ServiceID)
	if !ok {
		lager.Logger.Warnf(nil, "RegisterMicroServiceInstance get SelfInstancesCache failed, Mid/Sid: %s/%s", microServiceInstance.ServiceID, instanceID)
	}
	instanceIDs, ok := value.([]string)
	if !ok {
		lager.Logger.Warnf(nil, "RegisterMicroServiceInstance type asserts failed,  Mid/Sid: %s/%s", microServiceInstance.ServiceID, instanceID)
	}
	var isRepeat bool
	for _, va := range instanceIDs {
		if va == instanceID {
			isRepeat = true
		}
	}
	if !isRepeat {
		instanceIDs = append(instanceIDs, instanceID)
	}
	SelfInstancesCache.Set(microServiceInstance.ServiceID, instanceIDs, 0)
	lager.Logger.Warnf(nil, "RegisterMicroServiceInstance success, microServiceID/instanceID: %s/%s.", microServiceInstance.ServiceID, instanceID)

	return nil
}
