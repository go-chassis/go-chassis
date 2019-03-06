# CSE
CSE is cloud service engine, it has the core components(registry,distribution config, monitoring etc) for you to run your microservices 

## How to use 
You need to inject authentication in every request to cloud service engine.

1.Set AK SK in auth.yaml file 

```yaml
## Huawei Public Cloud ak/sk
cse:
  credentials:
    accessKey: xxx
    secretKey: xxx
```

2.import auth pkg in main.go

```go
import _ "github.com/huaweicse/auth/adaptor/gochassis"

```
this pkg will inject header into all of request to CSE engine 

After signing the header with authourization looks like this

Authorization: Credential=XXX, SignedHeaders=XXX, Signature=XXX

Complete [example](https://github.com/go-chassis/go-chassis-examples/tree/master/huaweicse)