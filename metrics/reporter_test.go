package metrics_test

import (
	"github.com/go-chassis/go-chassis/metrics"
	"github.com/go-chassis/go-chassis/metrics/prom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstallReporter(t *testing.T) {
	err := metrics.InstallReporter("test", prom.ReportMetricsToPrometheus)
	assert.NoError(t, err)
	err = metrics.InstallReporter("test", prom.ReportMetricsToPrometheus)
	assert.Error(t, err)
}
