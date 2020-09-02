# Rate limiting v2
## Introduction
Based on traffic marker, this rate limiter can limit traffic rate

## Usage
it allows you to define a list of rate limit policy.
each policy is defined as a yaml value, the yaml key path is servicecomb.rateLimiting.{policy_name}.
**match**
> *(required, string)* match policy name, if you want to limit un-marked traffic, use "none"

**rate**
> *(required, int)* allowed request in second
>
**burst**
> *(required, int)* if traffic rate reached rate limit. still allow n(burst) request going through 
>
>
Import middleware
```go
import _ github.com/go-chassis/go-chassis/v2/middleware/ratelimiter
```

Add marker and limiter to handler 
```yaml
servicecomb:
  handler:
    chain:
      Provider:
        default: traffic-marker,rate-limiter
```
#### Example
```yaml
servicecomb:
  match:
    traffic-to-some-api-from-jack: |
        headers:
          cookie:
            regex: "^(.*?;)?(user=jack)(;.*)?$"
          os:
            contains: linux
        apiPath:
          exact: "/some/api" 
        method: GET 
        trafficMarkPolicy: once
  rateLimiting:
    limiterPolicy1: |
      match: traffic-to-some-api-from-jack
      rate: 10
      burst: 1
    nonMarked: |
      match: none
      rate: 100
      burst: 1
```
