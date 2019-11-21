Go-Chassis 是一个go语言的微服务开发框架，帮助你快速开发微服务，完成架构转型

### 为什么使用 Go chassis

- go chassis被设计为通信协议中立的框架,你可以开发自定义的协议集成到go chassis， 并应用统一的治理能力，如负载均衡，熔断器，限流，流量控制，这些功能使你的服务变得具有韧性，并面向云原生。完全不需要再去自己集成各种方案，使用go chassis只需要通过配置文件来使用这些功能。

- go chassis 使用open tracing与prometheus使调用链与指标可视化

- go chassis 灵活性高，许多组件可以自己定制，比如注册发现，指标上报，调用链追踪 ，分布式配置管理等

- go chassis是插件化设计，所有的功能都是可插拔的，可以简化到成为一个restful框架。

# 特性
 - **可插拔的注册发现组件**: 当前支持Apache ServiceComb，kubernetes与istio，无论是服务端发现还是客户端注册发现都可以适配。
 - **插件化协议**: 当前支持http，grpc，支持开发者定制私有协议
 - **多端口管理**:  对于同协议，可以开放不同的端口，使用端口来划分API，面向不同调用方
 - **断路器**:  在运行时保护你的分布式系统，免于错误雪崩。
 - **流量管理**:  可以根据访问特征，微服务元数据，权重等规则灵活控制流量，可支持金丝雀发布等场景。
 - **客户端复杂均衡**: 支持定制策略，当前支持roundrobin，随机，会话粘纸与延迟权重。
 - **限流**:  客户端，服务端均可限流
 - **可插拔的加解密组件**:   加解密组件会被应用到mTLS等安全敏感的处理流程中，可自定义算法
 - **Handler Chain**:  可以在处理请求的过程中定制特殊逻辑，比如认证鉴权。
 - **Metrics**:  支持上报prometheus
 - **调用链追踪**: 集成opentracing-go作为标准，当前支持zipkin 
 - **运行时热加载配置**: 集成轻量级配置管理框架go-archaius, 配置可以在运行时热加载，无需重启，比如负载均衡，断路器，流量管理等配置
 - **原生支持动态配置框架**: 集成轻量级配置管理框架 go-archaius, 开发者可以实现拥有运行时配置热加载功能的应用
 - **API gateway与service mesh方案**: 由 [servicecomb-mesher](https://github.com/apache/servicecomb-mesher)提供. 
 - **Open API 2.0支持** go chassis会自动生成 Open API 2.0 文档并把它注册到Apache ServiceComb的service center. 你可以在统一的服务查看微服务文档。

go chassis插件库 [plugins](https://github.com/go-chassis/go-chassis-plugins) 可以查看目前的插件

# Get started 
1.生成 go mod
```bash
go mod init
```
2.增加go chassis 
```shell script
GO111MODULE=on go get github.com/go-chassis/go-chassis
```
如果你面临网络问题
```bash
export GOPROXY=https://goproxy.io
```

3.[开始编写你的第一个微服务](http://docs.go-chassis.com/getstarted/writing-rest.html)


# 文档
这是在线文档 [here](http://docs.go-chassis.com/), 
但他总是最新的版本, 如果你使用其他版本的go chassis
你可以跟随[这里](docs/README.md)来生成本地文档

# 例子
当前有两个例子库，提供很多使用场景
- [这里](examples)
- [这里](https://github.com/go-chassis/go-chassis-examples)

# 开发 go chassis

1. 安装[go 1.12+](https://golang.org/doc/install) 

2. Clone 工程

```sh
git clone git@github.com:go-chassis/go-chassis.git
```

3. 下载 vendors
```shell
cd go-chassis
export GO111MODULE=on 
go mod download
#optional
export GO111MODULE=on 
go mod vendor
```

4. 安装[Apache ServiceComb service-center](http://servicecomb.apache.org/)


# 使用go chassis开发的开源项目
- [apache/servicecomb-kie](https://github.com/apache/servicecomb-kie): 
A distributed configuration management service, go chassis and mesher integrate with it,
so that user can manage service configurations by this service.
- [apache/servicecomb-mesher](https://github.com/apache/servicecomb-mesher): 
A service mesh able to co-work with go chassis, 
it is able to run as a [API gateway](https://mesher.readthedocs.io/en/latest/configurations/edge.html) also.
- [KubeEdge](https://github.com/kubeedge/kubeedge): Kubernetes Native Edge Computing Framework (project under CNCF) https://kubeedge.io
