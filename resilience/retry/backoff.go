package retry

import (
	"github.com/cenkalti/backoff"
	"time"
)

//retry kind
const (
	KindExponential    = "exponential"
	KindConstant       = "constant"
	KindZero           = "zero"
	DefaultBackOffKind = KindExponential
)

//GetBackOff return the the back off policy
//min and max unit is million second
func GetBackOff(kind string, min, max int) backoff.BackOff {
	switch kind {
	case KindExponential:
		return &backoff.ExponentialBackOff{
			InitialInterval:     time.Duration(min) * time.Millisecond,
			RandomizationFactor: backoff.DefaultRandomizationFactor,
			Multiplier:          backoff.DefaultMultiplier,
			MaxInterval:         time.Duration(max) * time.Millisecond,
			MaxElapsedTime:      0,
			Clock:               backoff.SystemClock,
		}
	case KindConstant:
		return backoff.NewConstantBackOff(time.Duration(min) * time.Millisecond)
	case KindZero:
		return &backoff.ZeroBackOff{}
	default:
		return &backoff.ExponentialBackOff{}
	}

}
