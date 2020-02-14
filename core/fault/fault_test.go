package fault_test

import (
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/stretchr/testify/assert"

	"github.com/go-chassis/go-chassis/core/fault"
	"testing"
	"time"
)

func TestApplyFaultInjection(t *testing.T) {
	t.Run("add delay, should not return err", func(t *testing.T) {
		m := map[string]string{
			common.BuildinTagVersion: "0.1",
			common.BuildinTagApp:     "default"}
		inv := &invocation.Invocation{
			MicroServiceName: "service1",
			RouteTags: utiltags.Tags{
				KV:    m,
				Label: utiltags.LabelOfTags(m),
			},
		}
		err := fault.ValidateAndApplyFault(&model.Fault{
			Delay: model.Delay{
				Percent:    50,
				FixedDelay: 1 * time.Second,
			},
		}, inv)
		assert.NoError(t, err)
	})

	t.Run("add abort, should return err", func(t *testing.T) {
		m := map[string]string{
			common.BuildinTagVersion: "0.1",
			common.BuildinTagApp:     "default"}
		inv := &invocation.Invocation{
			MicroServiceName: "service1",
			RouteTags: utiltags.Tags{
				KV:    m,
				Label: utiltags.LabelOfTags(m),
			},
		}
		err := fault.ValidateAndApplyFault(&model.Fault{
			Abort: model.Abort{
				Percent:    100,
				HTTPStatus: 500,
			},
		}, inv)
		assert.Error(t, err)
	})
	t.Run("add delay and abort, should return err", func(t *testing.T) {
		m := map[string]string{
			common.BuildinTagVersion: "0.1",
			common.BuildinTagApp:     "default"}
		inv := &invocation.Invocation{
			MicroServiceName: "service1",
			RouteTags: utiltags.Tags{
				KV:    m,
				Label: utiltags.LabelOfTags(m),
			},
		}
		err := fault.ValidateAndApplyFault(&model.Fault{
			Delay: model.Delay{
				Percent:    50,
				FixedDelay: 1 * time.Second,
			},
			Abort: model.Abort{
				Percent:    100,
				HTTPStatus: 500,
			},
		}, inv)
		assert.Error(t, err)

		err = fault.ValidateAndApplyFault(&model.Fault{
			Delay: model.Delay{
				Percent:    50,
				FixedDelay: 1 * time.Second,
			},
			Abort: model.Abort{
				Percent:    100,
				HTTPStatus: 500,
			},
		}, inv)
		assert.Error(t, err)
	})

	t.Run("add abort and percent=0, should not return err", func(t *testing.T) {
		m := map[string]string{
			common.BuildinTagVersion: "0.1",
			common.BuildinTagApp:     "default"}
		inv := &invocation.Invocation{
			MicroServiceName: "service1",
			RouteTags: utiltags.Tags{
				KV:    m,
				Label: utiltags.LabelOfTags(m),
			},
		}
		err := fault.ValidateAndApplyFault(&model.Fault{
			Delay: model.Delay{
				Percent:    50,
				FixedDelay: 1 * time.Second,
			},
			Abort: model.Abort{
				Percent:    0,
				HTTPStatus: 500,
			},
		}, inv)
		assert.NoError(t, err)
	})

}
