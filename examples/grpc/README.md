# Quik start with examples

1. Launch service center
```sh
cd examples
docker-compose up
```

2. Run rest server

```sh 
cd examples/grpc/server
export CHASSIS_HOME=$PWD
go run main.go

```

3. Run Rest client
```sh 
 cd examples/grpc/client
 export CHASSIS_HOME=$PWD
 go run main.go
 
```