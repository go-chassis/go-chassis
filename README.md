# Go-Chassis  
[![Build Status](https://travis-ci.org/go-chassis/go-chassis.svg?branch=master)](https://travis-ci.org/go-chassis/go-chassis)  [![Coverage Status](https://coveralls.io/repos/github/go-chassis/go-chassis/badge.svg)](https://coveralls.io/github/go-chassis/go-chassis) [![Go Report Card](https://goreportcard.com/badge/github.com/go-chassis/go-chassis)](https://goreportcard.com/report/github.com/go-chassis/go-chassis) [![GoDoc](https://godoc.org/github.com/go-chassis/go-chassis?status.svg)](https://godoc.org/github.com/go-chassis/go-chassis) [![HitCount](http://hits.dwyl.io/go-chassis/go-chassis.svg)](http://hits.dwyl.io/go-chassis/go-chassis)  [![Join Slack](https://img.shields.io/badge/Join-Slack-orange.svg)](https://join.slack.com/t/go-chassis/shared_invite/enQtMzk0MzAyMjEzNzEyLTRjOWE3NzNmN2IzOGZhMzZkZDFjODM1MDc5ZWI0YjcxYjM1ODNkY2RkNmIxZDdlOWI3NmQ0MTg3NzBkNGExZGU)      

Go-Chassis is a Software Development Kit(SDK) for rapid development of microservices in Go



# Features
 - **Pluggable registrator and discovery service**: Support Service center,istio pilot and file based registry, 
 fit both client side discovery and server side discovery pattern 
 - **Pluggable Protocol**: You can custom your own protocol, by default support http and highway(RPC)
 - **Circuit breaker**: Protect your service in runtime or on-demand
 - **Route management**: Able to route to different service based on weight and match rule to achieve Canary Release easily
 - **Load balancing**: Able to custom strategy and filter
 - **Rate limiting**: Both client side and server side rate limiting
 - **Pluggable Cipher**: Able to custom your own cipher for AKSK and TLS certs
 - **Handler Chain**: Able to add your own code during service calling for client and server side
 - **Metrics**: Able to expose Prometheus metric API automatically and custom metrics reporter
 - **Tracing**:Use opentracing-go as standard library, easy to integrate tracing impl
 - **Logger**: You can custom your own writer to sink log, by default support file and stdout
 - **Hot-reconfiguraion**: A lot of configuration can be reload in runtime, like loadbalancing, circuit breaker, rate limiting
 - **Dynamic Configuration framework**:   Able to develop a service which has hot-reconfiguration feature easily
 - **Fault Injection**: In consumer side, you can inject faults to bring chaos testing into your system
 
You can check [plugins](https://github.com/go-chassis/go-chassis-plugins) to see more features

# Quick Start
You can see more documentations in [here](http://go-chassis.readthedocs.io/en/latest/)

1. Install [go 1.8+](https://golang.org/doc/install) 

2. Clone the project

```sh
git clone git@github.com:go-chassis/go-chassis.git
```

3. Use use go mod(go 1.11+, experimental but a recommended way)
```shell
cd go-chassis
GO111MODULE=on go mod download
#optional
GO111MODULE=on go mod vendor
```


4. Install [service-center](http://servicecomb.incubator.apache.org/release/)

5. [Write your first http micro service](http://go-chassis.readthedocs.io/en/latest/getstarted/writing-rest.html)



# Examples
You can check examples [here](examples)
# Communication Protocols
Go-Chassis supports two types of communication protocol.
1. Rest - REST is an approach that leverages the HTTP protocol for communication.
2. Highway - This is a RPC communication protocol

## Debug suggestion for dlv:
Add `-tags debug` into go build arguments before debugging, if your go version is go1.10 onward.

example:

```shell
go build -tags debug -o server -gcflags "all=-N -l" server.go
```

Chassis customized `debug` tag to resolve dlv debug issue:

https://github.com/golang/go/issues/23733

https://github.com/derekparker/delve/issues/865
