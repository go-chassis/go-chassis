# Router example

In this example we will run two different versions service and one client.We will show you how to 
implement grayscale publishing of the application with route management,version V1 is simulation
you old service , version V2 is simulation  you new service. 
 1. Launch service center
 
 ```sh
 cd examples
 docker-compose up
 ```
2. Build and Run Service V1 and Service V2
 
 
 we use microservice.yaml to set service config,and we must set two service with the 
 same service name,set different version for two service
 
 
Service V1
```yaml
#Private property of microservices
service_description:
  name: ROUTERServer
  version: 1.0
```
 
Service V2
```yaml
#Private property of microservices
service_description:
  name: ROUTERServer
  version: 2.0
```
 
 build and run two service
 
```bash
cd serverV1
go build main.go
./main 
```
```bash
cd serverV2
go build main.go
./main 
```

3. Build and run client

client use router.yaml to management router.Configuration of this file , when request header setting "Chassis:info",
all request access V1 . when request header setting "Chassis:say" , all request access V2 . Not setting anything in 
header ,request 80% access V1 , 20% access V2

set header
```go
req.SetHeader("Chassis", "info")
```
use router.yaml launch client 
```yaml
routeRule:
  ROUTERServer: # this value what set in the microservice.yaml file of service,service name of sc too
    - precedence: 1 # the big num the  precedence
      route: # router rules lists
      - tags:
          version: 1.0 #service version,sc default 0.1
        weight: 80 #weight 80% for here
      - tags:
          version: 2.0 #service version,sc default 0.1
        weight: 20 #weight 20% for here
    - precedence: 2
      match: # match strategy
        headers: 
          Chassis: # if request header setting info,will all access V1 
            regex: info
      route: 
      - tags:
          version: 1.0
        weight: 100 
    - precedence: 2
      match:
        headers:
          Chassis: # if request header setting say,will all access V2 
            regex: say
      route: 
      - tags:
          version: 2.0
        weight: 100 
```

build and run client
```bash
cd client
go build main.go
./main 
```
```bash
ROUTER Server equal num [POST]: version V1 : given num is equal the sum of the slice , num  : 10 ,sum : 10 
ROUTER Server equal num [POST]: version V2 : given num not  equal the product of the slice , num  : 10 ,product : 30
ROUTER Server equal num [POST]: version V1 : given num is equal the sum of the slice , num  : 10 ,sum : 10 
ROUTER Server equal num [POST]: version V1 : given num is equal the sum of the slice , num  : 10 ,sum : 10 
ROUTER Server equal num [POST]: version V1 : given num is equal the sum of the slice , num  : 10 ,sum : 10 
```