# Transport

## Introduction
you can define what can be considered as failure and make it count in circuit breaker and fault-tolerance module

## Configurations

**transport.failure.{protocol_name}**
> *(required, string)* the name of the protocol client, now only support rest failure. a string list connect with comma,each string 
starts with http_, combine with http status code.


## Example
The cases of http_500,http_502 are considered as unsuccessful attempts
```
cse:
  transport:
    failure:
      rest: http_500,http_502
```
