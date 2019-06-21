# Circuit breaker

## **Introduction**
Circuit breaker help to isolate upstream services during runtime,
all of invocation will be executed by circuit, under its protection, if there is too much error, time out or concurrency,
circuit will open to stop network communication
it also monitor each service call to make service [observable](https://docs.go-chassis.com/user-guides/metrics.html)

## **Configuration**

Circuit breaker scope is controlled by 

**cse.circuitBreaker.scope**
> *(optional, string)* service、instance or api, 
default is api, go chassis create a dedicated circuit for every api, invocation will be isolated based on api. 
if set to service, all of APIs of each service share one circuit, it will isolate the service.
if set to instance, each instance will get a dedicated circuit, it will isolate only one instance.
if set to api, each api will get a dedicated circuit, it will isolate only service api.
if set to instance-api, each instance api will get a dedicated circuit, it will isolate only one instance api.

Configuration Format looks like below：

cse.{namespace}.Consumer.{serviceName}.{property}

Explanation:

{namespace}：it can be isolation\|circuitBreaker\|fallback\|fallbackpolicy. 

{serviceName}: it is optional. it is the service level configuration, it represent the target service name

{property}: configuration items

**cse.isolation.timeoutInMilliseconds**
> *(optional, int)* if delay for a time, this call will be considered as failure, default is *30000*

**cse.isolation.maxConcurrentRequests**
> *(optional, int)* max concurrency, default is 1000

**cse.circuitBreaker.enabled**
> *(optional, bool)* enable circuit breaker or not, default is false

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
> *(optional, string)* fallback policy  [*returnnull*| *throwexception*]，default is returnnull. 
you can also [custom fallback policy](http://docs.go-chassis.com/dev-guides/circuit.html)


## **Examples**
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
    scope: api
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
usually you must set bizkeeper-consumer handler in chain before load balancing and transport
, here is a example 
```yaml
handler:
  chain:
    Consumer:
      default: bizkeeper-consumer, router, loadbalance, ratelimiter-consumer,transport
```

if you want to isolate instance or instance-api,you must set 
bizkeeper-consumer handler in chain after load balancing and before transport
,hear is a example
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
    scope: instance-api
    Consumer:
          enabled: false
          forceOpen: false
#......
```

```yaml
handler:
  chain:
    Consumer:
      default: router, loadbalance, bizkeeper-consumer, ratelimiter-consumer,transport
```
