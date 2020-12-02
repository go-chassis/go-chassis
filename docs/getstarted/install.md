Minimize Installation
=====
1.Install [go 1.13+](https://golang.org/doc/install) 

2.Generate go mod
```bash
go mod init
```
3.Add go chassis 
```bash
go get github.com/go-chassis/go-chassis/v2@v2.0.2
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
5.Install [service-center](https://service-center.readthedocs.io/en/latest/get-started/install.html)

6.[Write your first http micro service](https://go-chassis.readthedocs.io/en/latest/getstarted/writing-rest.html)


Use gRPC communication
===================
follow https://developers.google.com/protocol-buffers/docs/gotutorial to install grpc 

[Write your first grpc micro service](https://go-chassis.readthedocs.io//getstarted/writing-rpc.html)
