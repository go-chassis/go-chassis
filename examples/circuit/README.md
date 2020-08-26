# A circuit breaker example to show api level isolation
here is the config
```yaml
---
servicecomb:
  isolation:
    Consumer:
      timeoutInMilliseconds: 100
      maxConcurrentRequests: 1000
  circuitBreaker:
    scope: api # service|api
    Consumer:
      enabled: true
      forceOpen: false
      forceClosed: false
      sleepWindowInMilliseconds: 10000
      requestVolumeThreshold: 10
      errorThresholdPercentage: 10
  fallback:
    Consumer:
      enabled: true
  fallbackpolicy:
    Consumer:
      policy: throwexception
```
time out is 100ms

client calls a API which has a dead lock, circuit breaker will isolate this API

client can still call other API of this same service

if you change scope to "service"

it will isolate the service 

the runtime metrics is exported by prometheus exporter
check 127.0.0.1:5000/metrics and 127.0.0.1:5001/metrics to observe services