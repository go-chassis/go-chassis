package fault

import (
	"errors"
	"fmt"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
	"math/rand"
	"time"
)

// constant for default values values of abort and delay percentages
const (
	DefaultAbortPercentage int = 100
	DefaultDelayPercentage int = 100

	MaxPercentage int = 100
	MinPercentage int = 0
)

// ValidateAndApplyFault validate and apply the fault rule
func ValidateAndApplyFault(fault *model.Fault, inv *invocation.Invocation) error {
	if fault.Delay != (model.Delay{}) {
		if err := ValidateFaultDelay(fault); err != nil {
			return err
		}

		if err := ApplyFaultInjection(fault, inv, fault.Delay.Percent, "delay"); err != nil {
			return err
		}
	}

	if fault.Abort != (model.Abort{}) {
		if err := ValidateFaultAbort(fault); err != nil {
			return err
		}

		if err := ApplyFaultInjection(fault, inv, fault.Abort.Percent, "abort"); err != nil {
			return err
		}
	}

	return nil
}

// ValidateFaultAbort checks that fault injection abort HTTP status and Percentage is valid
func ValidateFaultAbort(fault *model.Fault) error {
	if fault.Abort.HTTPStatus < 100 || fault.Abort.HTTPStatus > 600 {
		return errors.New("invalid http fault status")
	}
	if fault.Abort.Percent < MinPercentage || fault.Abort.Percent > MaxPercentage {
		return fmt.Errorf("invalid httpfault percentage:must be in range 0..100")
	}

	return nil
}

// ValidateFaultDelay checks that fault injection delay fixed delay and Percentage is valid
func ValidateFaultDelay(fault *model.Fault) error {
	if fault.Delay.Percent < MinPercentage || fault.Delay.Percent > MaxPercentage {
		return errors.New("percentage must be in range 0..100")
	}

	if fault.Delay.FixedDelay < time.Millisecond {
		return errors.New("duration must be greater than 1ms")
	}

	return nil
}

//ApplyFaultInjection abort/delay
func ApplyFaultInjection(fault *model.Fault, inv *invocation.Invocation, configuredPercent int, faultType string) error {
	if rand.Intn(MaxPercentage)+1 <= configuredPercent {
		return injectFault(faultType, fault)
	}
	return nil
}

//injectFault apply fault based on the type
func injectFault(faultType string, fault *model.Fault) error {
	if faultType == "delay" {
		time.Sleep(fault.Delay.FixedDelay)
	}

	if faultType == "abort" {
		return errors.New("injecting abort")
	}

	return nil
}
