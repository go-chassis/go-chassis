package hystrix

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestConfigureConcurrency(t *testing.T) {
	convey.Convey("given a command configured for 100 concurrent requests", t, func() {
		ConfigureCommand("", CommandConfig{MaxConcurrentRequests: 100})

		convey.Convey("reading the concurrency should be the same", func() {
			convey.So(getSettings("").MaxConcurrentRequests, convey.ShouldEqual, 100)
		})
	})
}

func TestConfigureTimeout(t *testing.T) {
	convey.Convey("given a command configured for a 10000 milliseconds", t, func() {
		ConfigureCommand("", CommandConfig{Timeout: 10000})

		convey.Convey("reading the timeout should be the same", func() {
			convey.So(getSettings("").Timeout, convey.ShouldEqual, time.Duration(10*time.Second))
		})
	})
}

func TestConfigureRVT(t *testing.T) {
	convey.Convey("given a command configured to need 30 requests before tripping the circuit", t, func() {
		ConfigureCommand("", CommandConfig{RequestVolumeThreshold: 30})

		convey.Convey("reading the threshold should be the same", func() {
			convey.So(getSettings("").RequestVolumeThreshold, convey.ShouldEqual, uint64(30))
		})
	})
}

func TestSleepWindowDefault(t *testing.T) {
	convey.Convey("given default settings", t, func() {
		ConfigureCommand("", CommandConfig{})

		convey.Convey("the sleep window should be 5 seconds", func() {
			convey.So(getSettings("").SleepWindow, convey.ShouldEqual, time.Duration(5*time.Second))
		})
	})
}

func TestGetCircuitSettings(t *testing.T) {
	convey.Convey("when calling GetCircuitSettings", t, func() {
		ConfigureCommand("test", CommandConfig{Timeout: 30000})

		convey.Convey("should read the same setting just added", func() {
			convey.So(GetCircuitSettings()["test"], convey.ShouldEqual, getSettings("test"))
			convey.So(GetCircuitSettings()["test"].Timeout, convey.ShouldEqual, time.Duration(30*time.Second))
		})
	})
}
