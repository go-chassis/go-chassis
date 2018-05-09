# Quik start with examples

1. Launch service center
```sh
cd examples
docker-compose up
```

2. Run rest server

```sh 
cd server
export CHASSIS_HOME=$PWD
go run main.go

```

3. Run Rest client
```sh 
 cd client
 export CHASSIS_HOME=$PWD
 go run main.go
 
```

4. check zikpin at http://127.0.0.1:9411

5. check metrics at http://127.0.0.1:5001

 
