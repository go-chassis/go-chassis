Simple Rest service

1.Launch service center

follow https://github.com/apache/servicecomb-service-center/tree/master/examples/infrastructures/docker

2.Run rest server

```sh 
cd examples/rest/server
export CHASSIS_HOME=$PWD
go run main.go

```

3.Run Rest client
```sh 
 cd examples/rest/client
 export CHASSIS_HOME=$PWD
 go run main.go

 
```

you can find rest api doc in local, in the meantime it is uploaded automatically to service center.
```shell
vim conf/RESTServer/schema/RESTServer.yaml
``` 
