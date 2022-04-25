# Monitoring

## **Introduction**
monitoring handler aim to monitor traffics of server side microservice
see how to export metrics to prometheus [observable](https://go-chassis.readthedocs.io/en/latest/user-guides/metrics.html)

it records 3 different metrics:
- request_count
- request_process_duration
- error_response_count

all the metrics name starts with "scb_", it stands for servicecomb system
## **Usage**


1.Import it in your main file
```go
import _ github.com/go-chassis/go-chassis/v2/middleware/monitoring
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
# HELP scb_request_count 
# TYPE scb_request_count counter
scb_request_count{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0"} 14
# HELP scb_request_process_duration 
# TYPE scb_request_process_duration summary
scb_request_process_duration{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0",quantile="0.5"} 3
scb_request_process_duration{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0",quantile="0.9"} 80
scb_request_process_duration{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0",quantile="0.99"} 80
scb_request_process_duration_sum{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0"} 315
scb_request_process_duration_count{app="default",env="",instance="",service="servicecomb-kie",version="0.1.0"} 14
```




