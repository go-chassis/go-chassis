# Introduction

### What is Go chassis

Go chassis is a micro service framework for Go developer. you can develop distributed system with go chassis rapidly.


### Why use Go chassis

go chassis is designed as a protocol-independent framework, any protocol is able to integrate with go chassis and leverage same function like load balancing,
circuit breaker,rate limiting, those function resilient your service

go chassis makes service observable by bring open tracing and prometheus to it.

go chassis is flexible, many different modules can be replaced by other implementation, 
like registry, metrics, handler chain, config center etc 

With many build-in function like route management, circuit breaker, load balancing, monitoring etc,
your don't need to search and integrate a solution yourself

go chassis supports Istio platform, Although Istio is a great platform with a service mesh in data plane, 
it surely decrease the throughput and increase the latency of your service and cost more CPU usage. 
go chassis can bring better performance to go program, you can use Istio configurations to control go chassis.