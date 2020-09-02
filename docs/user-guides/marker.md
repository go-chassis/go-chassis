# Traffic marker(alpha)
## Introduction
Traffic marker module is able to mark requests in both client(consumer) or server(provider) side,
it is the foundation of traffic management.

it is a alpha feature.
## Configurations
it allows you to define a list of match policy.
each policy is defined as a yaml value, the yaml key path is servicecomb.match.{policy_name}.

**apiPath**
> *(optional, map)* relative http API path, for example "/v2/some_component/some_resource".
> you can use operator to specify how to match this API, like "contains,exact,regex"

**headers**
> *(optional, map)* define a condition map, the logic is "and".

**method**
> *(optional, string)* standard http method in upper case, for example GET,POST

**trafficMarkPolicy**
> *(optional, string)* cloud be once or perService. 
> if set to once, then once marked by a service, all services will reuse this mark, 
> if set to perService, each service will ignore mark result which is done by other service and try to mark traffic again.

must add marker to handler 
```yaml
servicecomb:
  handler:
    chain:
      Provider:
        default: traffic-marker,rate-limiter
```
## Supported features
More feature will will support marker, currently support below features:

- router management
- rate limiter

## Example
you can mark traffic and then limit the rate of this kind of traffic
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
 
```

