A example to show launch same protocol's server and separate API by ports

server separate API by ports
```go
chassis.RegisterSchema("rest", &schemas.Hello{})
chassis.RegisterSchema("rest-legacy", &schemas.Legacy{})
chassis.RegisterSchema("rest-admin", &schemas.Admin{})
```
use chassis.yaml to launch API servers
```yaml
  protocols:
    rest:
      listenAddress: 127.0.0.1:5001
      advertiseAddress: 127.0.0.1:5001
    rest-legacy:
      listenAddress: 127.0.0.1:5002
      advertiseAddress: 127.0.0.1:5002
    rest-admin:
      listenAddress: 127.0.0.1:5003
      advertiseAddress: 127.0.0.1:5003
```
build and run server
```
cd server
go build main.go
./main
```
client use different adress to access API server
```go
	req, err := rest.NewRequest("GET", "http://RESTServer/hello")

	req, err = rest.NewRequest("GET", "http://RESTServer:legacy/legacy")

```
build and run client
```
cd client
go build main.go
./main
```
```bash
2018-09-19 17:20:36.065 +08:00 INFO client/client_manager.go:86 Create client for rest:RESTServer:127.0.0.1:5001
2018-09-19 17:20:36.066 +08:00 INFO client/main.go:34 REST Server sayhello[GET]: hi from hello
2018-09-19 17:20:36.067 +08:00 INFO client/client_manager.go:86 Create client for rest:RESTServer:127.0.0.1:5002
2018-09-19 17:20:36.067 +08:00 INFO client/main.go:49 REST Server sayhello[GET]: hello from legacy

```


access to url to see different response, you can see API is separate by ports

http://127.0.0.1:5001/hello 

http://127.0.0.1:5002/legacy

http://127.0.0.1:5003/admin
