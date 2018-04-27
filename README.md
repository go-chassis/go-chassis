# Go-Chassis  
[![Build Status](https://travis-ci.org/ServiceComb/go-chassis.svg?branch=master)](https://travis-ci.org/ServiceComb/go-chassis)  [![Coverage Status](https://coveralls.io/repos/github/ServiceComb/go-chassis/badge.svg)](https://coveralls.io/github/ServiceComb/go-chassis) [![Go Report Card](https://goreportcard.com/badge/github.com/ServiceComb/go-chassis)](https://goreportcard.com/report/github.com/ServiceComb/go-chassis) [![GoDoc](https://godoc.org/github.com/ServiceComb/go-chassis?status.svg)](https://godoc.org/github.com/ServiceComb/go-chassis) [![HitCount](http://hits.dwyl.io/ServiceComb/go-chassis.svg)](http://hits.dwyl.io/ServiceComb/go-chassis)      

Go-Chassis is a Software Development Kit(SDK) for rapid development of microservices in Go
 
Go-chassis is based on [Go-Micro](https://github.com/micro/go-micro) A pluggable RPC framework

# Features
 - **Pluggable registrator and discovery service**: Support Service center,istio pilot and file based registry by default
 - **Dynamic Configuration framework**:  you are able to develop a service which has hot-reconfiguration  feature easily
 - **Pluggable Protocol**: You can custom your own protocol,by default support http and highway(RPC)
 - **Circuit breaker**: Protect your service in runtime or on-demand
 - **Routing management**: Able to route to different service based on weight and match rule to achieve Canary Release easily
 - **Load balancing**: Add custom strategy and filter
 - **Rate limiting**: Both client side and server side rate limiting
 - **Pluggable Cipher**: Able to custom your own cipher for AKSK and TLS certs
 - **Handler Chain**: Able to add your own code during service calling for client and server side
 - **Metrics**: Able to expose Prometheus metric API automatically and sink metrics to CSE Dashboard
 - **Tracing**: Integrate with Zipkin and namedpipe to sink tracing data
 - **Logger**: You can custom your own writer to sink log, by default support file and stdout
 - **Hot-reconfiguraion**: A lot of configuration can be reload in runtime, like loadbalancing, circuit breaker, rate limiting
 
# Quick Start
You can see more informations in here http://go-chassis.readthedocs.io/en/latest/

## Installation
1. Install go 1.8+ https://golang.org/doc/install

2. Clone the project

```sh
git clone git@github.com:ServiceComb/go-chassis.git
```

3. Use gvt to download deps

```sh
go get -u github.com/FiloSottile/gvt
cd go-chassis 
gvt restore
```

4. Install ServiceComb service-center https://github.com/ServiceComb/service-center/releases

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
chassis.RegisterSchema("rest", &HelloServer{},server.WithSchemaID("HelloServer"))
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
restInvoker := core.NewRestInvoker()
req, _ := rest.NewRequest("GET", "cse://"+providerName+"/sayhello")
resp1, err := restInvoker.ContextDo(context.TODO(), req)
```

# Examples
You can check examples [here](examples)
# Communication Protocols
Go-Chassis supports two types of communication protocol.
1. Rest - REST is an approach that leverages the HTTP protocol for communication.
2. Highway - This is a high performance communication protocol originally developed by Huawei. 

