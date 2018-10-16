Minimize Installation
=====
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


Use gRPC communication
===================
follow https://developers.google.com/protocol-buffers/docs/gotutorial to install grpc 

[Write your first grpc micro service](http://go-chassis.readthedocs.io/en/latest/getstarted/writing-rpc.html)