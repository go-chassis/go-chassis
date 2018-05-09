package metrics_test

import (
	"github.com/ServiceComb/go-chassis/metrics"
	"github.com/ServiceComb/go-chassis/metrics/prom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstallReporter(t *testing.T) {
	err := metrics.InstallReporter("test", prom.ReportMetricsToPrometheus)
	assert.NoError(t, err)
	err = metrics.InstallReporter("test", prom.ReportMetricsToPrometheus)
	assert.Error(t, err)
}
