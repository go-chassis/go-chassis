# Fault Tolerance
## Introduction

go-chassis support fault-tolerance, so that you can define if error occurs, how to deal with this error.

## Configuration

the fault-tolerace related configurations is all in load_balancing.yaml, prefix is cse.loadbalance.
set retryEnabled to true to enable it



**retryEnabled**
> *(optional, bool)* Enable fault tolerance, default is *false*

**retryOnSame**
> *(optional, int)* if remote call failed, then retry on same instance, default is *0*

**retryOnNext**
> *(optional, int)* if remote call failed, then call load balancing again to get next instance, default is *0*

**backoff.kind**
> *(optional, string)* backoff policy: [jittered|constant|zero] default is *zero*
- zero:  do not wait for any time。
- constant: after each failed retry, wait for constant time. Use backoff.minMs to set the time。
- jittered: time wil exponential growth after each retry, till this time reach to MaxMs. 
Use backoff.minMs to set the the first wait time

**backoff.MinMs**
> *(optional, int)* minimum wait time between each retry, unit is ms, default is *0*

**backoff.MaxMs**
> *(optional, int)* maximum wait time between each retry, unit is ms, default is *0*

## example

edit load_balancing.yaml.

```yaml
cse:
  loadbalance:
    retryEnabled: true
    retryOnNext: 2
    retryOnSame: 3
    backoff:
      kind: jittered
      MinMs: 200
      MaxMs: 400
```



