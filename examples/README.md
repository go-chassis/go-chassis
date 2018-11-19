# Quik start with examples

1. Launch service center

follow https://github.com/apache/servicecomb-service-center/tree/master/examples/infrastructures/docker

2. Run rest server

```sh 
cd examples/rest/server
export CHASSIS_HOME=$PWD
go run main.go

```

3. Run Rest client
```sh 
 cd examples/rest/client
 export CHASSIS_HOME=$PWD
 go run main.go
 
```

# Examples

### communication
 
A simple end to end communication without any service registry involved

### discovery

A complicated example with most of go chassis-features

### rest

simple rest communication 

### rpc

simple rpc communication

### file upload

a file upload example by using go-chassis rest communication

### pilot 

this example use  Istio pilot as service registry

### monitoring

show tracing on zipkin and export prometheus metrics format  

### metadata

Demonstrate how to set local scope parameter, header, and use them in your schema

In handler chain you can set a local data or protocol header, then you can read it in next handler in chain and your Restful handler

### router

This example uses the rest protocol to show how to perform route management and achieve  grayscale publishing.
