Go-Chassis 是一个go语言的微服务开发框架，专注于帮你实现云原生应用

### 为什么使用 Go chassis
- 强大的中间件 "handler chain":不止拥有 "filter" or "interceptor"的能力. chain中每个handler都可以拿到后面的handler的执行结果，包括业务代码的执行结果。这在很多场景下都很实用，比如:

1.跟踪业务指标，并导出他们让prometheus收集。

2.跟踪关键的业务执行结果，审计这些信息。

3.分布式调用链追踪，end span可以等到业务执行后在handler中去完善，无需写在业务代码中。

4. 客户端调用远程服务时，也需要进行中间处理，比如客户端负载均衡，请求重试

以上场景的共性是帮助你解耦通用功能与业务逻辑，解放业务开发者。否则业务逻辑将于这些通用功能耦合，你可以将这些工作交给基础设施团队，开发出来的中间件，可以供任何业务团队使用。而业务团队只需要专注于业务逻辑开发。

- go chassis沉淀了许多可信软件所需的特性：开箱即用的限流，校验请求，金丝雀发布，高可用，通信保护等功能

- go chassis被设计为通信协议中立的框架,你可以开发自定义的协议集成到go chassis， 并应用统一的治理能力，如负载均衡，熔断器，限流，流量控制，这些功能使你的服务变得具有韧性，并面向云原生。完全不需要再去自己集成各种方案，使用go chassis只需要通过配置文件来使用这些功能。

- go chassis 使用open tracing与prometheus使调用链与指标可视化

- go chassis 灵活性高，许多组件可以自己定制，比如注册发现，指标上报，调用链追踪，分布式配置管理等

- go chassis是插件化设计，所有的功能都是可插拔的，且功能可按需引入，不引入就不会编译到二进制执行文件中。开源三方件依赖很少

# 特性

Less dependencies: checkout the go.mod file, it has less dependency on open source project by default, to import more features checkout plugins to see more features
 - **可插拔的注册发现组件**: 当前支持Apache ServiceComb，kubernetes与istio，无论是服务端发现还是客户端注册发现都可以适配。
 - **插件化协议**: 当前支持http，grpc，支持开发者定制私有协议
 - **多端口管理**:  对于同协议，可以开放不同的端口，使用端口来划分API，面向不同调用方
 - **丰富的中间件**:  利用handler chaiun，提供了多种通用中间件，比如认证鉴权，限流，重试，流量标记等
 - **流量标记**:  定义流量特征并为他标记为一个独有的字符，便于后续根据特征进行流量管理
 - **流量管理**:  可以根据访问特征，微服务元数据，权重等规则灵活控制流量，可支持金丝雀发布，限流等场景。
 - **安全**: cipher插件化设计，可以对接不同加解密算法
 - **客户端复杂均衡**: 
 - **韧性**:  支持重试，限流，客户端负载均衡，断路器，在业务受损，请求失败时保证关键任务运行正常。
 - **遥测**:  提供metrics抽象API，并且默认收集请求数，延迟等通用指标，支持prometheus，集成opentracing-go作为标准，当前支持zipkin 。
 - **后端服务**: 将后端服务视为插件使用，比如配额管理，认证鉴权服务。可以便于测试，并保证组件的可替换性
 - **运行时热加载配置**: 集成轻量级配置管理框架go-archaius, 配置可以在运行时热加载，无需重启，比如负载均衡，断路器，流量管理等配置
 - **原生支持动态配置框架**: 集成轻量级配置管理框架 go-archaius, 开发者可以实现拥有运行时配置热加载功能的应用
  - **API first** go chassis会自动生成 Open API 2.0 文档并把它注册到Apache ServiceComb的service center. 你可以在统一的服务查看微服务文档。
 - **spring cloud与service mesh统一治理**: [servicecomb-mesher](https://github.com/apache/servicecomb-mesher)， [spring cloud](https://github.com/huaweicloud/spring-cloud-huawei)提供。
 -  **极少的开源依赖** 查看go.mod文件，开源库依赖已经做到最少依赖，更多的功能可以查看[插件库](https://github.com/go-chassis/go-chassis-extension)


go chassis插件库 [plugins](https://github.com/go-chassis/go-chassis-extension) 可以查看目前的插件

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

3.[开始编写你的第一个微服务](https://go-chassis.readthedocs.io/en/latest/getstarted/writing-rest.html)


# 文档
这是在线文档 [here](https://go-chassis.readthedocs.io/), 
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

# Known Users

> [欢迎在此登录自己的信息](https://github.com/go-chassis/go-chassis/issues/592)

![huawei](https://raw.githubusercontent.com/go-chassis/go-chassis.github.io/master/known_users/huawei.PNG) 
![qutoutiao](https://raw.githubusercontent.com/go-chassis/go-chassis.github.io/master/known_users/qutoutiao.PNG)
![Shopee](https://raw.githubusercontent.com/go-chassis/go-chassis.github.io/master/known_users/Shopee.png)

