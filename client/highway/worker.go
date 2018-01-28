package tcp

import (
	"errors"
	"fmt"
	"github.com/ServiceComb/go-chassis/client/highway/pb"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/util/metadata"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"strconv"
	"sync"
	"time"
)

//TODO 1. configurable
//TODO 2. could be defined for each service
//TODO 3. dynamic configuration,in other word manually scale
//TODO 4. Auto scale ,Worker should scale in and out according network traffic to a service

//ElapseTime is a constant of type float64
const ElapseTime float64 = 10

type job struct {
	req  *client.Request
	ctx  context.Context
	resp interface{}
	err  chan error
}

type worker struct {
	name string
	pm   proto.Message
	//each worker bind to one conn
	client          transport.Client
	jobs            chan *job
	loggedIn        bool
	reqHeader       *highway.RequestHeader
	respHeader      *highway.ResponseHeader
	loginBody       *highway.LoginRequest
	reqHeaderBytes  []byte
	respHeaderBytes []byte
	respBodyBytes   []byte
}

//each addr has a channel buffer, worker receive jobs from channel buffer
var jobSchdlr *jobSchedule

type jobSchedule struct {
	jobChannel map[string]*JobInfo
	sync.RWMutex
}

//JobInfo is a struct
type JobInfo struct {
	JobChan     chan *job
	RefreshTime time.Time
}

func init() {
	jobSchdlr = new(jobSchedule)
	jobSchdlr.jobChannel = make(map[string]*JobInfo)

	//cleanup the workers if not used for long time(10min)
	go jobSchdlr.cleanUp()
}

func (js *jobSchedule) scheduleJob(addr string, jobInfo *job) error {

	jobDetails, isExist := js.isWorkerExist(addr)
	if !isExist {
		return errors.New("worker does not exist")
	}

	//update refresh time
	jobDetails.RefreshTime = time.Now()

	jobDetails.JobChan <- jobInfo
	return nil
}

func (js *jobSchedule) isWorkerExist(addr string) (*JobInfo, bool) {
	js.RLock()
	job, ok := js.jobChannel[addr]
	if !ok {
		js.RUnlock()
		return nil, false
	}

	js.RUnlock()
	return job, true
}

func (js *jobSchedule) createWorkerSchedulers(addr string, workerNum int, c *highwayClient, host string) error {
	var ch chan *job
	_, ok := js.isWorkerExist(addr)
	if !ok {
		//create buffered channel which will receive jobs
		ch = make(chan *job, 1000)
		lager.Logger.Infof("Starting %d workers for %s", workerNum, addr)
		js.Lock()
		js.jobChannel[addr] = &JobInfo{
			JobChan:     ch,
			RefreshTime: time.Now(),
		}
		js.Unlock()
		for i := 0; i < workerNum; i++ {
			trClient, err := c.opts.Transport.Dial(host)
			if err != nil {
				return err
			}

			worker := &worker{
				jobs:   ch,
				client: trClient,
				name:   addr + "#" + strconv.Itoa(i),
			}

			go worker.run(addr)
		}
	}
	//workers already exist
	return nil
}

func (js *jobSchedule) cleanUp() {
	//check for the refresh time for each five minutes and cleanup the workers
	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ticker.C:
			js.schedulerCleanup()
		}
	}
}

func (js *jobSchedule) schedulerCleanup() {
	//range over the address and if exist, check the last served request time is more than ten minutes
	//if so the worker is inactive remove address and close the  worker channel
	for addr := range js.jobChannel {
		jobDetails, isExist := js.isWorkerExist(addr)
		if isExist {
			refreshedTime := jobDetails.RefreshTime
			if time.Now().Sub(refreshedTime).Minutes() > ElapseTime {
				//close the channel
				close(jobDetails.JobChan)
				//delete the addr and close the channel
				js.Lock()
				delete(js.jobChannel, addr)
				js.Unlock()
			}
		}
	}

}

func releaseWorkers(addr string) {
	//if connection refused error then close the channel and remove the addr
	jobDetails, isExist := jobSchdlr.isWorkerExist(addr)
	if isExist {
		//close the channel
		close(jobDetails.JobChan)
		//delete the addr and close the channel
		jobSchdlr.Lock()
		delete(jobSchdlr.jobChannel, addr)
		jobSchdlr.Unlock()
	}

}

func (w *worker) run(addr string) {
	w.reqHeader = &highway.RequestHeader{}
	w.respHeader = &highway.ResponseHeader{}
	w.loginBody = &highway.LoginRequest{}

	for {
		select {
		case job, ok := <-w.jobs:
			if !ok {
				// channel closed
				lager.Logger.Info("channel is closed")
				return
			}

			if w.loggedIn == false {
				w.login(job, addr)
			}

			w.sendRequest(job, addr)
		}
	}
}

func (w *worker) login(j *job, addr string) {

	//TODO send(Login)
	loginRequestHeader, loginRequestBody, _, err := w.marshalLogin()
	if err != nil {
		j.err <- err
		return
	}
	codeChan := make(chan bool, 1)
	go func() {
		err = w.client.Send(loginRequestHeader, loginRequestBody, nil, j.req.ID)
		if err != nil {
			j.err <- err
			releaseWorkers(addr)
			return
		}
		loginResponseHeader, loginResponseBody, _, _, err := w.client.Recv()

		if err != nil {
			j.err <- err
			releaseWorkers(addr)
			return
		}
		loginRespHeader := &highway.ResponseHeader{}
		loginRespBody := &highway.LoginResponse{}
		//解码响应头
		err = proto.Unmarshal(loginResponseHeader, loginRespHeader)
		if err != nil {
			lager.Logger.Errorf(err, "unmarshal login response header failed")
			j.err <- err
			return
		}
		//解码响应Body
		err = proto.Unmarshal(loginResponseBody, loginRespBody)
		if err != nil {
			lager.Logger.Errorf(err, "unmarshal login response body  failed")
			j.err <- err
			return
		}

		//TODO 获取远端编码方式
		codeChan <- loginRespBody.UseProtobufMapCodec
	}()

	remoteCodeType := <-codeChan
	if remoteCodeType == true {
		close(codeChan)
		w.loggedIn = true
		return
	}

}

func (w *worker) sendRequest(j *job, addr string) {

	requestHeader, requestBody, _, err := w.marshalRequest(j.ctx, j.req)
	if err != nil {
		j.err <- err
		return
	}
	done := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				j.err <- errors.New("Highway request error")
				done <- false
			}
		}()

		//TODO medata 在highway中暂时不用
		err = w.client.Send(requestHeader, requestBody, nil, j.req.ID)

		if err != nil {
			j.err <- err
			done <- false
			releaseWorkers(addr)
			return
		}

		w.respHeaderBytes, w.respBodyBytes, _, _, err = w.client.Recv()
		if err != nil {
			j.err <- err
			done <- false
			releaseWorkers(addr)
			return
		}
		// success
		done <- true
	}()

	d := <-done
	if !d {
		close(done)
		return
	}
	close(done)

	err = proto.Unmarshal(w.respHeaderBytes, w.respHeader)
	if err != nil {
		lager.Logger.Errorf(err, "unmarshal response header failed.")
		j.err <- err
		return
	}

	if w.respHeader.StatusCode != 200 {
		j.err <- fmt.Errorf("highway get an error: %s", w.respHeader.Reason)
		return
	}

	err = proto.Unmarshal(w.respBodyBytes, j.resp.(proto.Message))
	if err != nil {
		lager.Logger.Errorf(err, "unmarshal response body failed.")
		j.err <- err
		return
	}

	j.err <- nil
}

func (w *worker) marshalRequest(ctx context.Context, req *client.Request) ([]byte, []byte, map[string]string, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	body, err := proto.Marshal(req.Arg.(proto.Message))
	if err != nil {
		return nil, body, nil, fmt.Errorf("client marshal request args failed %s", err)
	}

	//flags:标示是否被压缩 暂时不用
	//reuse mem
	*w.reqHeader = highway.RequestHeader{
		MsgType:          highway.MsgTypeRequest,
		Flags:            int32(0),
		DestMicroservice: req.MicroServiceName,
		OperationName:    req.Method,
		SchemaID:         req.Struct,
		Context:          md,
	}

	header, err := proto.Marshal(w.reqHeader)
	if err != nil {
		lager.Logger.Errorf(err, "client marshal highway request header failed.")
		return header, body, nil, err
	}
	return header, body, nil, nil

}

func (w *worker) marshalLogin() ([]byte, []byte, map[string]string, error) {
	//reuse mem
	*w.reqHeader = highway.RequestHeader{
		MsgType:          highway.MsgTypeLogin,
		Flags:            int32(0),
		DestMicroservice: "",
		OperationName:    "",
		SchemaID:         "",
		Context:          nil,
	}
	header, err := proto.Marshal(w.reqHeader)
	if err != nil {
		lager.Logger.Errorf(err, "client marshal highway login header failed")
		return header, nil, nil, nil
	}
	*w.loginBody = highway.LoginRequest{
		Protocol:            "highway",
		ZipName:             "z",
		UseProtobufMapCodec: localSupportLogin,
	}
	body, err := proto.Marshal(w.loginBody)
	if err != nil {
		lager.Logger.Errorf(err, "client marshal highway login body failed")
		return header, body, nil, nil
	}
	return header, body, nil, nil
}
func (w *worker) freeRequestHeader() {
	*w.reqHeader = highway.RequestHeader{}
}
