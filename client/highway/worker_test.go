package tcp

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJobschedule_schedulerCleanup(t *testing.T) {
	assert := assert.New(t)
	jobScheduleObj := &jobSchedule{
		jobChannel: make(map[string]*JobInfo),
	}
	jobScheduleObj.jobChannel["addr1"] = &JobInfo{
		JobChan:     make(chan *job, 10),
		RefreshTime: time.Date(2017, 8, 17, 20, 34, 58, 651387237, time.UTC),
	}
	jobScheduleObj.jobChannel["addr1"].JobChan <- &job{}
	jobScheduleObj.schedulerCleanup()
	assert.Equal(len(jobScheduleObj.jobChannel), 0)

}
func TestRealeaseWorker(t *testing.T) {
	assert := assert.New(t)
	jobSchdlr = &jobSchedule{
		jobChannel: make(map[string]*JobInfo),
	}
	jobSchdlr.jobChannel["addr1"] = &JobInfo{
		JobChan:     make(chan *job, 10),
		RefreshTime: time.Date(2017, 8, 17, 20, 34, 58, 651387237, time.UTC),
	}
	jobSchdlr.jobChannel["addr1"].JobChan <- &job{}
	releaseWorkers("addr1")
	_, ok := jobSchdlr.jobChannel["addr1"]
	assert.Equal(ok, false)

}
