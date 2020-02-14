package fault_test

import (
	"fmt"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/fault"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_ApplyFaultInjection(t *testing.T) {
	var inv = new(invocation.Invocation)
	var fault1 = new(model.Fault)
	inv.Endpoint = "1.2.3.4"
	inv.MicroServiceName = "Server"
	fault1.Delay.FixedDelay = 2 * time.Second
	//delay must not return error
	v := fault.ApplyFaultInjection(fault1, inv, 100, "delay")
	assert.Equal(t, v, nil)
	//abort must return error
	v = fault.ApplyFaultInjection(fault1, inv, 100, "abort")
	assert.Equal(t, v, fmt.Errorf("injecting abort"))
}
