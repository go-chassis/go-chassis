package hystrix_test

import (
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix/reporter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstallReporter(t *testing.T) {
	err := hystrix.InstallReporter("test", reporter.ReportMetricsToPrometheus)
	assert.NoError(t, err)
	err = hystrix.InstallReporter("test", reporter.ReportMetricsToPrometheus)
	assert.Error(t, err)
}
