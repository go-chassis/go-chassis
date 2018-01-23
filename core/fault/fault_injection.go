package fault

import (
	"errors"
	"fmt"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"time"
)

var (
	faultKeyCount             map[string]int
	initialKeyCount           map[string]int
	invokedCount              map[string]int
	abortNeeded, delayApplied bool
)

func init() {
	faultKeyCount = make(map[string]int)
	initialKeyCount = make(map[string]int)
	invokedCount = make(map[string]int)
}

// constant for default values values of abort and delay percentages
const (
	DefaultAbortPercentage int = 100
	DefaultDelayPercentage int = 100
)

// ValidateAndApplyFault validate and apply the fault rule
func ValidateAndApplyFault(fault *model.Fault, inv *invocation.Invocation) error {
	if fault.Delay != (model.Delay{}) {
		if err := ValidateFaultDelay(fault); err != nil {
			return err
		}

		if fault.Abort != (model.Abort{}) {
			abortNeeded = true
		}

		if err := ApplyFaultInjection(fault, inv, fault.Delay.Percent, "delay"); err != nil {
			return err
		}

	}

	if fault.Abort != (model.Abort{}) {
		if err := ValidateFaultAbort(fault); err != nil {
			return err
		}

		// In case both delay and abort specified ,then fault injection mechanism will apply delay followed by abort
		// percentage of fault injection will be done based on the specified percentage for delay
		if abortNeeded && delayApplied {
			abortNeeded = false
			delayApplied = false
			return errors.New("injecting abort and delay")
		}

		if !abortNeeded {
			if err := ApplyFaultInjection(fault, inv, fault.Abort.Percent, "abort"); err != nil {
				return err
			}
		}

	}

	return nil
}

// ValidateFaultAbort checks that fault injection abort HTTP status and Percentage is valid
func ValidateFaultAbort(fault *model.Fault) error {
	if fault.Abort.HTTPStatus < 100 || fault.Abort.HTTPStatus > 600 {
		return errors.New("invalid httpfault status")
	}
	if fault.Abort.Percent < 0 || fault.Abort.Percent > 100 {
		return fmt.Errorf("invalid httpfault percentage:must be in range 0..100")
	}

	if fault.Abort.Percent == 0 {
		fault.Abort.Percent = DefaultAbortPercentage
	}

	return nil
}

// ValidateFaultDelay checks that fault injection delay fixed delay and Percentage is valid
func ValidateFaultDelay(fault *model.Fault) error {
	if fault.Delay.Percent < 0.0 || fault.Delay.Percent > 100.0 {
		return errors.New("percentage must be in range 0..100")
	}

	if fault.Delay.Percent == 0 {
		fault.Delay.Percent = DefaultDelayPercentage
	}

	if fault.Delay.FixedDelay < time.Millisecond {
		return errors.New("duration must be greater than 1ms")
	}

	return nil
}

//ApplyFaultInjection abort/delay
func ApplyFaultInjection(fault *model.Fault, inv *invocation.Invocation, configuredPercent int, faultType string) error {
	key := inv.MicroServiceName + "-" + inv.Version + "-" + inv.AppID
	count, exist := invokedCount[key]
	if !exist {
		count = 1
		faultKeyCount[key]++
	}

	if exist && count == 1 {
		faultKeyCount[key]++
	}

	failureCount := faultKeyCount[key]
	initialCount, ok := initialKeyCount[key]

	if !ok {
		initialCount = 0
	}

	percentage := calculatePercentage(count, configuredPercent)

	if isMutationOccrred(failureCount, percentage) {
		resetFaultKeyCount(key, percentage-1)
	}

	if percentage == failureCount && initialCount != 1 {
		initialKeyCount[key]++
		incrementKeyCount(key, count+1)
		err := injectFault(faultType, fault)
		return err

	}

	if percentage != failureCount && percentage > 1 {
		faultKeyCount[key]++
		incrementKeyCount(key, count+1)
		err := injectFault(faultType, fault)
		return err
	}

	incrementKeyCount(key, count+1)
	return nil
}

//incrementKeyCount increment the key count with respect to the instance which is going to serve the request
func incrementKeyCount(key string, count int) {
	invokedCount[key] = count

}

//calculatePercentage calculate percentage as whole number
func calculatePercentage(count, percent int) int {
	return (count * percent / 100)
}

//isMutationOccrred confirm whether the configured percent has mutation occured
func isMutationOccrred(actualCount, expectedCount int) bool {
	return (actualCount-expectedCount) > 0 || (expectedCount-actualCount) > 1
}

//resetFaultKeyCount reset the faultKeyCount to expectedCount
func resetFaultKeyCount(key string, expectedCount int) {
	faultKeyCount[key] = expectedCount
}

//injectFault apply fault based on the type
func injectFault(faultType string, fault *model.Fault) error {
	if faultType == "delay" {
		delayApplied = true
		time.Sleep(fault.Delay.FixedDelay)
	}

	if faultType == "abort" {
		return errors.New("injecting abort")
	}

	return nil
}
