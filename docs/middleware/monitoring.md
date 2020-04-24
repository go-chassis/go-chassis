# Monitoring

## **Introduction**
monitoring handler aim to monitor traffics of server side microservice
see how to export metrics to prometheus [observable](https://go-chassis.readthedocs.io/en/latest/user-guides/metrics.html)

it records 3 different metrics:
- request_count
- request_process_duration
- error_response_count

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
request_count{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0"} 14
# HELP request_process_duration 
# TYPE request_process_duration summary
request_process_duration{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0",quantile="0.5"} 3
request_process_duration{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0",quantile="0.9"} 80
request_process_duration{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0",quantile="0.99"} 80
request_process_duration_sum{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0"} 315
request_process_duration_count{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0"} 14
```




