Minimize Installation
=====
1. Install [go 1.8+](https://golang.org/doc/install)

2. Clone the project

```sh
git clone git@github.com:go-chassis/go-chassis.git
```

3. Use [dep](https://github.com/golang/dep) to download dependencies

```sh
cd go-chassis
# behind a proxy, you need setup a http proxy server
# export https_proxy=xxx
dep ensure
```

4. Install go-chassis [service-center](http://servicecomb.incubator.apache.org/release/)


Use RPC communication
===================
Install protobuff 3.2.0 https://github.com/google/protobuf

Install https://github.com/golang/protobuf
