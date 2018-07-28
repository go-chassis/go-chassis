# Go-Chassis  
[![Build Status](https://travis-ci.org/go-chassis/go-chassis.svg?branch=master)](https://travis-ci.org/go-chassis/go-chassis)  [![Coverage Status](https://coveralls.io/repos/github/go-chassis/go-chassis/badge.svg)](https://coveralls.io/github/go-chassis/go-chassis) [![Go Report Card](https://goreportcard.com/badge/github.com/go-chassis/go-chassis)](https://goreportcard.com/report/github.com/go-chassis/go-chassis) [![GoDoc](https://godoc.org/github.com/go-chassis/go-chassis?status.svg)](https://godoc.org/github.com/go-chassis/go-chassis) [![HitCount](http://hits.dwyl.io/go-chassis/go-chassis.svg)](http://hits.dwyl.io/go-chassis/go-chassis)  [![Join Slack](https://img.shields.io/badge/Join-Slack-orange.svg)](https://join.slack.com/t/go-chassis/shared_invite/enQtMzk0MzAyMjEzNzEyLTRjOWE3NzNmN2IzOGZhMzZkZDFjODM1MDc5ZWI0YjcxYjM1ODNkY2RkNmIxZDdlOWI3NmQ0MTg3NzBkNGExZGU)      

Go-Chassis is a Software Development Kit(SDK) for rapid development of microservices in Go
 
Go-chassis is based on [Go-Micro](https://github.com/micro/go-micro) A pluggable RPC framework



# Features
 - **Pluggable registrator and discovery service**: Support Service center,istio pilot and file based registry by default
 - **Pluggable Protocol**: You can custom your own protocol, by default support http and highway(RPC)
 - **Circuit breaker**: Protect your service in runtime or on-demand
 - **Routing management**: Able to route to different service based on weight and match rule to achieve Canary Release easily
 - **Load balancing**: Able to custom strategy and filter
 - **Rate limiting**: Both client side and server side rate limiting
 - **Pluggable Cipher**: Able to custom your own cipher for AKSK and TLS certs
 - **Handler Chain**: Able to add your own code during service calling for client and server side
 - **Metrics**: Able to expose Prometheus metric API automatically and sink metrics to CSE Dashboard
 - **Tracing**: Integrate with Zipkin and namedpipe to sink tracing data
 - **Logger**: You can custom your own writer to sink log, by default support file and stdout
 - **Hot-reconfiguraion**: A lot of configuration can be reload in runtime, like loadbalancing, circuit breaker, rate limiting
 - **Dynamic Configuration framework**:  you are able to develop a service which has hot-reconfiguration feature easily
 - **Fault Injection**: In consumer side, chassis support you to inject faults 
 
You can check [plugins](https://github.com/go-chassis/go-chassis-plugins) to see more features

# Quick Start
You can see more documentations in [here](http://go-chassis.readthedocs.io/en/latest/)

1. Install [go 1.8+](https://golang.org/doc/install)

2. Clone the project

```sh
git clone git@github.com:go-chassis/go-chassis.git
```

3. Use [glide](https://github.com/Masterminds/glide) to download dependencies

```sh
cd go-chassis 
glide intall
```

4. Install [service-center](https://github.com/go-chassis/service-center/releases)

5. [Write your first http micro service](http://go-chassis.readthedocs.io/en/latest/getstarted/writing-rest.html)



# Examples
You can check examples [here](examples)
# Communication Protocols
Go-Chassis supports two types of communication protocol.
1. Rest - REST is an approach that leverages the HTTP protocol for communication.
2. Highway - This is a high performance communication protocol originally developed by Huawei. 

