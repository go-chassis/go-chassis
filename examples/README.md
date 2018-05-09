# Quik start with examples

1. Launch service center
```sh
cd examples
docker-compose up
```

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