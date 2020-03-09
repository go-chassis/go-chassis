package registry

import (
	"sync"
	"time"

	"github.com/go-chassis/go-chassis/core/config"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-mesh/openlogging"
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
	shutdown bool
	mux      sync.Mutex
}

// Start start the heartbeat system
func (s *HeartbeatService) Start() {
	s.shutdown = false
	s.run()
}

// Stop stop the heartbeat system
func (s *HeartbeatService) Stop() {
	s.shutdown = true
}

// DoHeartBeat do heartbeat for each instance
func (s *HeartbeatService) DoHeartBeat(microServiceID, microServiceInstanceID string) {
	_, err := DefaultRegistrator.Heartbeat(microServiceID, microServiceInstanceID)
	if err != nil {
		openlogging.GetLogger().Errorf("heartbeat fail,try to recover, err: %s", err)
		s.RetryRegister(microServiceID, microServiceInstanceID)
	}
}

// run runs the heartbeat system
func (s *HeartbeatService) run() {
	for !s.shutdown {
		s.DoHeartBeat(runtime.ServiceID, runtime.InstanceID)
		time.Sleep(30 * time.Second)
	}
}

//RetryRegister try to register micro-service, and instance
func (s *HeartbeatService) RetryRegister(sid, iid string) {
	for !s.shutdown {
		openlogging.Info("try to re-register")
		var err error
		if _, err = DefaultServiceDiscoveryService.GetMicroService(sid); err != nil {
			err = s.ReRegisterSelfMSandMSI()
			if err != nil {
				openlogging.Error("recover failed:" + err.Error())
			} else {
				openlogging.Warn("recovered service")
			}
		}
		err = reRegisterSelfMSI(sid, iid)
		if err != nil {
			openlogging.Error("recover failed:" + err.Error())
		} else {
			openlogging.Warn("recovered instance")
			break
		}
		time.Sleep(DefaultRetryTime)
	}
}

// ReRegisterSelfMSandMSI 重新注册微服务和实例
func (s *HeartbeatService) ReRegisterSelfMSandMSI() error {
	err := RegisterService()
	if err != nil {
		openlogging.GetLogger().Errorf("The reRegisterSelfMSandMSI() startMicroservice failed: %s", err)
		return err
	}

	err = RegisterServiceInstances()
	if err != nil {
		openlogging.GetLogger().Errorf("The reRegisterSelfMSandMSI() startInstances failed: %s", err)
		return err
	}
	return nil
}

// reRegisterSelfMSI 只重新注册实例
func reRegisterSelfMSI(sid, iid string) error {
	eps, err := MakeEndpointMap(config.GlobalDefinition.Cse.Protocols)
	if err != nil {
		return err
	}
	if len(InstanceEndpoints) != 0 {
		eps = make(map[string]*Endpoint, len(InstanceEndpoints))
		for m, ep := range InstanceEndpoints {
			epObj, err := NewEndPoint(ep)
			if err != nil {
				continue
			}
			eps[m] = epObj
		}
	}
	microServiceInstance := &MicroServiceInstance{
		InstanceID:   iid,
		EndpointsMap: eps,
		HostName:     runtime.HostName,
		Status:       common.DefaultStatus,
		Metadata:     runtime.InstanceMD,
	}
	instanceID, err := DefaultRegistrator.RegisterServiceInstance(sid, microServiceInstance)
	if err != nil {
		openlogging.GetLogger().Errorf("RegisterInstance failed: %s", err)
		return err
	}
	openlogging.GetLogger().Infof("register instance success, microServiceID/instanceID: %s/%s.", sid, instanceID)
	return nil
}
