Minimize Installation
=====
1.Install [go 1.12+](https://golang.org/doc/install) 

2.Generate go mod
```bash
go mod init
```
3.Add go chassis 
```bash
GO111MODULE=on go get github.com/go-chassis/go-chassis
```

4.Use go mod
    ```bash
    GO111MODULE=on go mod download
    #optional
    GO111MODULE=on go mod vendor
    ```
if you are facing network issue 
```bash
export GOPROXY=https://goproxy.io
```
5.Install [service-center](http://servicecomb.apache.org/release/)

6.[Write your first http micro service](http://docs.go-chassis.com/getstarted/writing-rest.html)


Use gRPC communication
===================
follow https://developers.google.com/protocol-buffers/docs/gotutorial to install grpc 

[Write your first grpc micro service](http://docs.go-chassis.com/getstarted/writing-rpc.html)