package registry

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/openlog"
	"github.com/gorilla/websocket"
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
	HeartbeatMode string
	shutdown      bool
	Interval      time.Duration
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
func (s *HeartbeatService) DoHeartBeat(microServiceID, microServiceInstanceID string, instanceHeartbeatMode string) error {
	callback := func() {
		openlog.Error("heartbeat fail,try to recover")
		// register instance
		s.RetryRegister(microServiceID, microServiceInstanceID)
	}
	err := DefaultRegistrator.Heartbeat(microServiceID, microServiceInstanceID, instanceHeartbeatMode, callback)
	return err
}

// run runs the heartbeat system
func (s *HeartbeatService) run() {
	if s.HeartbeatMode == PersistenceHeartBeat {
		err := s.DoHeartBeat(runtime.ServiceID, runtime.InstanceID, s.HeartbeatMode)
		if err != nil {
			openlog.Error("send persistence heartbeat failed: " + err.Error())
			if errors.Is(err, websocket.ErrBadHandshake) {
				openlog.Info("send non-persistence heartbeat")
				s.HeartbeatMode = NonPersistenceHeartBeat
				s.sendNonPersistenceHeartBeat()
			}
		}
		return
	}
	s.sendNonPersistenceHeartBeat()
}

// sendNonPersistenceHeartBeat use http to send heartbeat
func (s *HeartbeatService) sendNonPersistenceHeartBeat() {
	for !s.shutdown {
		// first, wait for the successful registration, and then send the heartbeat to slow down the pressure of SC
		time.Sleep(s.Interval)
		err := s.DoHeartBeat(runtime.ServiceID, runtime.InstanceID, s.HeartbeatMode)
		if err != nil {
			openlog.Error("send heartbeat failed: " + err.Error())
		}
	}
}

//RetryRegister try to register micro-service, and instance
func (s *HeartbeatService) RetryRegister(sid, iid string) {
	for !s.shutdown {
		openlog.Info("try to re-register")
		var err error
		if _, err = DefaultServiceDiscoveryService.GetMicroService(sid); err != nil {
			err = s.ReRegisterSelfMSandMSI()
			if err != nil {
				openlog.Error("recover failed:" + err.Error())
			} else {
				openlog.Warn("recovered service")
			}
		}
		err = reRegisterSelfMSI(sid, iid)
		if err != nil {
			openlog.Error("recover failed:" + err.Error())
		} else {
			openlog.Warn("recovered instance")
			break
		}
		time.Sleep(DefaultRetryTime)
	}
}

// ReRegisterSelfMSandMSI 重新注册微服务和实例
func (s *HeartbeatService) ReRegisterSelfMSandMSI() error {
	err := RegisterService()
	if err != nil {
		openlog.Error(fmt.Sprintf("The reRegisterSelfMSandMSI() startMicroservice failed: %s", err))
		return err
	}

	err = RegisterServiceInstances()
	if err != nil {
		openlog.Error(fmt.Sprintf("The reRegisterSelfMSandMSI() startInstances failed: %s", err))
		return err
	}
	return nil
}

// reRegisterSelfMSI 只重新注册实例
func reRegisterSelfMSI(sid, iid string) error {
	eps, err := MakeEndpointMap(config.GlobalDefinition.ServiceComb.Protocols)
	if err != nil {
		return err
	}
	if len(InstanceEndpoints) != 0 {
		eps = make(map[string]*Endpoint, len(InstanceEndpoints))
		var epObj *Endpoint
		for m, ep := range InstanceEndpoints {

			epObj, err = NewEndPoint(ep)
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
		openlog.Error(fmt.Sprintf("register instance failed: %s", err))
		return err
	}
	openlog.Info(fmt.Sprintf("register instance success, microServiceID/instanceID: %s/%s.", sid, instanceID))
	return nil
}
