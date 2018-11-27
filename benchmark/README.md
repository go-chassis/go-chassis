# Performance Benchmarking Test

This tool helps to evaluate the performance of your service developed by [go-chassis](https://github.com/go-chassis/go-chassis).


### Quick Start

1. setup service center
2. set service center url in chassis.yaml
3. start your service 
3. benchmark
```bash
cd benchmark
go build github.com/go-chassis/go-chassis/benchmark
./benchmark -c 10 -d 30s -u http://service/hello
```

TPS:  11481.3
Total request:  114813
Err request:  0
latency mean:  0.43261272373540854
latency p05:  0.14776679999999998
latency p25:  0.2251135
latency p50:  0.312795
latency p75:  0.447414
latency p90:  0.6835228000000001
latency p99:  2.7049685300000124