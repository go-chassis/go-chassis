package accesslog

import (
	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/initiator"
	"github.com/go-chassis/go-chassis/pkg/util/iputil"
	"github.com/go-mesh/openlogging"
	"time"
)

type Recorder func(startTime time.Time, i *invocation.Invocation)

var (
	instance = &accessLog{
		recorder: restfulRecorder,
	}

	log openlogging.Logger
)

const handlerNameAccessLog = "access-log"

func init() {
	if initiator.LoggerOptions == nil || len(initiator.LoggerOptions.AccessLogFile) == 0 {
		openlogging.GetLogger().Info("lager.yaml non exist, skip init")
		return
	}

	if initiator.LoggerOptions.AccessLogFile == "stdout" {
		log = openlogging.GetLogger()
	} else {
		var err error
		opts := &lager.Options{
			Writers:        lager.File,
			LoggerLevel:    lager.LevelInfo,
			LoggerFile:     initiator.LoggerOptions.AccessLogFile,
			LogFormatText:  initiator.LoggerOptions.LogFormatText,
			RollingPolicy:  initiator.LoggerOptions.RollingPolicy,
			LogRotateDate:  initiator.LoggerOptions.LogRotateDate,
			LogRotateSize:  initiator.LoggerOptions.LogRotateSize,
			LogBackupCount: initiator.LoggerOptions.LogBackupCount,
		}
		log, err = lager.NewLog(opts)
		if err != nil {
			openlogging.GetLogger().Errorf("new access log failed, %s", err.Error())
			return
		}
	}

	err := handler.RegisterHandler(handlerNameAccessLog, func() handler.Handler {
		return instance
	})
	if err != nil {
		openlogging.GetLogger().Errorf("register access handler failed, %s", err.Error())
	}
}

// CustomizeRecorder support customize recorder
func CustomizeRecorder(recorder Recorder) {
	instance.recorder = recorder
}

type accessLog struct {
	recorder func(time.Time, *invocation.Invocation)
}

// Handle ...
func (a *accessLog) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	now := time.Now()
	chain.Next(i, func(response *invocation.Response) error {
		err := cb(response)
		if err != nil {
			return err
		}
		a.recorder(now, i)
		return nil
	})
}

// Name ...
func (a *accessLog) Name() string {
	return handlerNameAccessLog
}

func restfulRecorder(startTime time.Time, i *invocation.Invocation) {
	req := i.Args.(*restful.Request)
	resp := i.Reply.(*restful.Response)
	log.Infof("%s %s from %s %d %dms", req.Request.Method, req.Request.URL.String(),
		iputil.ClientIP(req.Request), resp.StatusCode(), time.Now().Sub(startTime).Nanoseconds()/1000000)
}
