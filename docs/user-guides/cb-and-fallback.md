# Circuit breaker

## **Introduction**
Circuit breaker help to prevent network failure between service call, 
it also monitor each service call to make service [observable](https://go-chassis.readthedocs.io/en/latest/user-guides/metrics.html)

## **Configuration**

Configuration Format as below：

cse.{namespace}.Consumer.{serviceName}.{property}
explanation:

{namespace}：it can be isolation\|circuitBreaker\|fallback\|fallbackpolicy. 

{serviceName}: it is optional. it is the service level configuration, it represent the target service name

{property}: configuration items

**cse.isolation.timeoutInMilliseconds**
> *(optional, int)* if delay for a time, this call will be considered as failure, default is *30000*

**cse.isolation.maxConcurrentRequests**
> *(optional, int)* max concurrency, default is 1000

**cse.circuitBreaker.enabled**
> *(optional, bool)* enable circuit breaker or not, default is true

**cse.circuitBreaker.forceOpen**
> *(optional, bool)* if it is true, will forcely open the circuit, default is false

**cse.circuitBreaker.forceClosed**
> *(optional, bool)* ignore all configurations forcely close crcuit all the time, default is false

**cse.circuitBreaker.sleepWindowInMilliseconds**
> *(optional, int)* after a circuit open, how long it should wait for next retry, 
if retry failed, circuit will open again.
>default is 15000

**cse.circuitBreaker.requestVolumeThreshold**
> *(optional, int)* it means in 10 seconds after how many request fails, circuit breaker should open
> default is 20

**cse.circuitBreaker.errorThresholdPercentage**
> *(optional, int)* it means how many err percentage met, circuit breaker should open, default is 50

**cse.fallback.enabled**
> *(optional, bool)* enable fallback or not, default is true

**cse.fallbackpolicy.policy**
> *(optional, string)* fallback policy  [*returnnull*| *throwexception*]，default is returnnull


## **examples**
```yaml
cse:
  isolation:
    Consumer:
      timeoutInMilliseconds: 1
      maxConcurrentRequests: 100
      ServerA: # service level config
        timeoutInMilliseconds: 1000
        maxConcurrentRequests: 1000
  circuitBreaker:
    Consumer:
      enabled: false
      forceOpen: false
      forceClosed: true
      sleepWindowInMilliseconds: 10000
      requestVolumeThreshold: 20
      errorThresholdPercentage: 10
      ServerB: # service level config
        enabled: true
        forceOpen: false
        forceClosed: false
        sleepWindowInMilliseconds: 10000
        requestVolumeThreshold: 20
        errorThresholdPercentage: 5
  fallback:
    Consumer:
      enabled: true
      maxConcurrentRequests: 20
  fallbackpolicy:
    Consumer:
      policy: throwexception
```
you must set bizkeeper-consumer handler in chain before load balancing and transport
here is a example 
```yaml
handler:
  chain:
    Consumer:
      default: bizkeeper-consumer, router, loadbalance, ratelimiter-consumer,transport
```


