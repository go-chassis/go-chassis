# ServiceStage
ServiceStage 是华为云推出的云服务，通过其中的CSE(cloud service engine)子服务，你可以快速启动ServiceComb service center服务，配置中心或是其他面向云原生开发的服务。每个服务集群被称为engine。

## 如何使用
### 通过配置文件配置AK SK
1.在auth.yaml文件中进行配置

```yaml
## Huawei Public Cloud ak/sk
cse:
  credentials:
    accessKey: xxx
    secretKey: xxx
```

2.必须在main.go中import auth插件

```go
import _ "github.com/huaweicse/auth/adaptor/gochassis"

```
这个插件会在http的request header中注入认证信息

完整的[example](https://github.com/go-chassis/go-chassis-examples/tree/master/huaweicse)
### 通过使用ServiceStage部署，免AKSK配置
使用ServiceStage进行部署的微服务无需进行ak sk手工配置，框架自动发现service center等服务的地址
1.必须在main.go中import auth插件

```go
import _ "github.com/huaweicse/auth/adaptor/gochassis"

```
2.在此行代码之后import 华为的扩展组件用于自动查询服务的endpoint
```go
import _ "github.com/go-chassis/go-chassis-cloud/provider/huawei/engine"
```
3.在chassis.yaml中设置引擎名字
```yaml
servicecomb:
  engine:
    name: test-engine
```
完整的[example](https://github.com/go-chassis/go-chassis-cloud/tree/master/example)