# Profiling
## Overview

Go-chassis provides the ability to inquire route rules and discovered microservice instance information in the current program cache when the program is running.

It allows users to easily locate issues of route and service discovery.

## Configurations

**cse.profile.enable**
> *(optional, bool)* If it is true, 
a new http API defined in "cse.profile.apiPath" will serve for client.
Default is *false*.

**cse.profile.apiPath**
> *(optional, string)* It's the root path of the profile interface,
default is */profile*.
The specific profile path will be under this root path.


## Example

```yaml
servicecomb:
  profile:
    enable: true
    apiPath: /profile
```

If the rest is listening on 127.0.0.1:8080, after performing the above configuration,
you can get route rules through [http://127.0.0.1:8080/profile/route-rule](http://127.0.0.1:8080/profile/route-rule) and
discovered microservice instance information through [http://127.0.0.1:8080/profile/discovery](http://127.0.0.1:8080/profile/discovery).

Or you can get all profile data through root path [http://127.0.0.1:8080/profile](http://127.0.0.1:8080/profile).
It includes information for all the above sub-paths.

