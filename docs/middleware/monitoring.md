# Monitoring

## **Introduction**
monitoring handler aim to monitor traffics of server side microservice
see how to export metrics to prometheus [observable](https://docs.go-chassis.com/user-guides/metrics.html)

it records 3 different metrics:
- request_count
- request_latency
- request_errors

## **Usage**


1.Import it in your main file
```go
import _ github.com/go-chassis/go-chassis/middleware/monitoring
```

2. you must set monitoring handler in chain provider chain
   , here is a example 
   ```yaml
   handler:
     chain:
       Provider:
         default: monitoring
   ```




metrics API response looks like below：
```text
# HELP request_count 
# TYPE request_count counter
request_count{app="",env="",instance="",service="",version=""} 4
# HELP request_errors_count 
# TYPE request_errors_count counter
request_errors_count{app="",code="Ǵ",env="",instance="",service="",version=""} 2
# HELP request_latency 
# TYPE request_latency histogram
request_latency_bucket{app="",env="",instance="",service="",version="",le="0.05"} 4
request_latency_bucket{app="",env="",instance="",service="",version="",le="0.25"} 4
request_latency_bucket{app="",env="",instance="",service="",version="",le="0.5"} 4
request_latency_bucket{app="",env="",instance="",service="",version="",le="0.75"} 4
request_latency_bucket{app="",env="",instance="",service="",version="",le="0.9"} 4
request_latency_bucket{app="",env="",instance="",service="",version="",le="0.99"} 4
request_latency_bucket{app="",env="",instance="",service="",version="",le="0.995"} 4
request_latency_bucket{app="",env="",instance="",service="",version="",le="+Inf"} 4
request_latency_sum{app="",env="",instance="",service="",version=""} 0
request_latency_count{app="",env="",instance="",service="",version=""} 4
```




