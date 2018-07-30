# Performance Benchmarking Test

This repo helps to evaluate the performance of [go-chassis](https://github.com/go-chassis/go-chassis).

This repo consist of a client and a server which is written is Go using go-chassis sdk. The client calls the api of the server and records the resource usage in csv file which can be used to analyse the aggregate performance of the go-chassis.

### Quick Start

To run this performance benchmark please follow the below steps:

##### Step-1
Clone Go-Chassis and download the dependency.
```go

git clone https://github.com/go-chassis/go-chassis $GOPATH/src/github.com/go-chassis/go-chassis

cd $GOPATH/src/github.com/go-chassis/go-chassis

glide install
```

##### Step-2
Start the service-center locally using this [guide](https://github.com/apache/incubator-servicecomb-service-center#quick-start).

##### Step-3
Build and run server
```go
cd benchmark/server

go build

./server
```

##### Step-4
Build and run client
```go
cd benchmark/client

go build

./client
```

##### Step-5 Verify the performance report

In client a csv file will be generated naming with the timestamp, the files contains the aggregate value for the performance of client and server.

