package config_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	lbBytes := []byte(`
servicecomb: 
  loadbalance: 
    TargetService: 
      backoff: 
        maxMs: 400
        minMs: 200
        kind: constant
      retryEnabled: false
      retryOnNext: 2
      retryOnSame: 3
      serverListFilters: zoneaware
      strategy: 
        name: WeightedResponse
    backoff: 
      maxMs: 400
      minMs: 200
      kind: constant
    retryEnabled: false
    retryOnNext: 2
    retryOnSame: 3
    serverListFilters: zoneaware
    strategy: 
      name: WeightedResponse

`)
	d, _ := os.Getwd()
	filename3 := filepath.Join(d, "load_balancing.yaml")
	f3, _ := os.OpenFile(filename3, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	_, _ = f3.Write(lbBytes)
	m.Run()
}
