# Go-chassis
[![Build Status](https://travis-ci.org/ServiceComb/go-chassis.svg?branch=master)](https://travis-ci.org/ServiceComb/go-chassis)  [![Coverage Status](https://coveralls.io/repos/github/ServiceComb/go-chassis/badge.svg)](https://coveralls.io/github/ServiceComb/go-chassis)
Go-Chassis is a Software Development Kit(SDK) for rapid development of microservices in GoLang,
 providing service-discovery,  fault-tolerance, circuit breaker, load balancing, monitoring, hot-reconfiguration features 

Inspired by [Go-Micro](https://github.com/micro/go-micro)
A pluggable RPC framework for distributed systems development.

Go-chassis is not a enhancement based on it, but  did something more
# Features
 - **Pluggable registry**: Support Service center and file based registry by default
 - **Dynamic Configuration framework**:  you are able to develop a service which has hot-reconfiguration  feature easily
 - **Pluggable Protocol**: You can custom your own protocol,by default support http and highway(RPC)
 - **Circuit breaker**: Protect your service in runtime or on-demand
 - **Load balancing**: You can custom strategy and filter
 - **Rate limiting**: Both client side and server side rate limiting
 - **Pluggable Cipher**: Able to custom your own cipher for AKSK and TLS certs
 - **Handler Chain**: Able to add your own code during service calling for client and server side
 - **Metrics**: Able to expose Prometheus metric API automatically and sink metrics to CSE Dashboard
 - **Tracing**: Integrate with Zipkin and namedpipe to sink tracing data
 - **Logger**: You can custom your own writer to sink log, by default support file and stdout
 - **Hot-reconfiguraion**: A lot of configuration can be reload in runtime, like loadbalancing, circuit breaker, rate limiting
 
# Quick Start
You can see more informations in gitbook https://go.huaweicse.cn/

## Write a http service provider

<b>Step 1:</b>
Define your Schema and your business logic.

```go
//API
func (s *HelloServer) SayHello(b *restful.Context) {
	b.Write([]byte("Hello : Welcome to Go-Chassis."))
}
//Specify URL pattern
func (s *HelloServer) URLPatterns() []restful.Route {
	return []restful.Route{
		{http.MethodGet, "/sayhello", "SayHello"},
	}
}
```

<b>Step 2:</b>
Register your Schema to go-chassis
```go
chassis.RegisterSchema("rest", &HelloServer{},
		server.WithSchemaID("HelloServer"))
```

<b>Step 3:</b>
Start the Chassis as a Server
```go
	chassis.Init()
	chassis.Run()
```

## Write a Http Consumer

<b>Step 1:</b>
Initialize your Chassis
```go
chassis.Init()

```
<b>Step 2:</b>
Use Rest Invoker to call the provider
```go
restinvoker := core.NewRestInvoker()
	req, _ := rest.NewRequest("GET", "cse://"+providerName+"/sayhello")
	resp1, err := restinvoker.ContextDo(context.TODO(), req)
```

# Examples
You can check examples [here](examples)
# Communication Protocols
Go-Chassis supports two types of communication protocol.
1. Rest - REST is an approach that leverages the HTTP protocol for communication.
2. Highway - This is a high performance communication protocol originally developed by Huawei. 

