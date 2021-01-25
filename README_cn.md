Go-Chassis 是一个go语言的微服务开发框架，专注于帮你实现云原生应用。Logo的含义是开发者可以通过引入go chassis重新创造和定制自己的“轮子”（即开发框架），以此来加速云原生应用的交付速度

## 为什么使用 Go chassis
强大的中间件 "handler chain":不止拥有 "filter" or "interceptor"的能力. chain中每个handler都可以拿到后面的handler的执行结果，包括业务代码的执行结果。这在很多场景下都非常实用，比如:

1. 跟踪业务指标，并导出让prometheus收集。

2. 跟踪关键的业务执行结果，审计这些信息。

3. 分布式调用链追踪，end span可以等到业务执行后在handler中去完善，无需写在业务代码中。

4. 客户端调用远程服务时，也会需要中间件进行处理。比如客户端负载均衡，请求重试

以上场景的共性是帮助你解耦通用功能与业务逻辑，解放业务开发者。业务开发者可以将这些工作交给基础设施团队，自己只专注于业务逻辑开发。基础设施团队开发出来的中间件，可作为基础能力供任给何业务团队使用。

- 强大的特性基线：
开箱即用的限流，请求校验，金丝雀发布，客户端负载均衡，通信保护等功能，
这些特性均面向云原生应用的开发场景。

- 面向交付：
在将代码编译发布后，向软件使用者交付最终应用软件包（一个或多个）时，
可以通过插件机制按需引用不同插件以满足不同用户诉求。同时还可以根据业务场景对软件按需裁剪，例如开源与商业的源码是不同的，不同用户对功能诉求不同等。

- 通信协议中立：
你可以将自定义的协议集成到go chassis中。应用统一特性基线，面向云原生。go chassis只需要通过配置文件来使用这些功能，减去了集成各种方案的烦恼。

- 轻量级内核：所有的功能都是可插拔的支持按需引入。不引入的功能就不会编译到二进制执行文件中。此外，开源三方件依赖很少。

## 特性

 - **注册发现**: 当前支持Apache ServiceComb，kubernetes与Istio，无论是服务端发现还是客户端注册发现都可以适配。
 - **客户端负载均衡**: consumer实时缓存依赖服务的网络信息拓扑，并直接进行负载均衡算法选择
 - **流量标记**:  定义流量特征并为他标记为一个独有的字符，便于后续根据特征进行流量管理
 - **流量管理**:  可以根据访问特征，微服务元数据，权重等规则灵活控制流量，可支持金丝雀发布，限流等场景。
 - **韧性**:  支持重试、限流、客户端负载均衡、断路器。在业务受损和请求失败时保证关键任务运行正常。
 - **丰富的中间件**:  利用handler chain，提供了多种通用中间件。比如认证鉴权、限流、重试、流量标记等。
 - **插件化协议**: 当前支持http、gRPC，支持开发者定制私有协议。
 - **多端口管理**:  对于同协议可以开放不同的端口。使用端口来划分API，面向不同调用方。
 - **内置安全**: cipher插件化设计，可以定制不同加解密算法。
 - **遥测**:  提供metrics抽象API，并且默认收集请求数、延迟等通用指标。支持prometheus、zipkin。集成opentracing-go作为标准。
 - **后端服务**: 将后端服务视为插件使用，比如配额管理、认证鉴权服务。这样便于测试并保证组件的可替换性。
 - **原生支持配置热加载**: 集成轻量级配置管理框架 go-archaius, 开发者可以轻松实现配置热加载功能的云应用。
 - **API first**： 自动生成 Open API 2.0 文档，并把它注册到Apache ServiceComb的service center。 可在统一的服务查看微服务文档。
 - **spring cloud与service mesh统一治理**: 由[servicecomb-mesher](https://github.com/apache/servicecomb-mesher)， [spring cloud](https://github.com/huaweicloud/spring-cloud-huawei)提供。
 - **极少的开源依赖** 查看go.mod文件，已做到做少的开源库依赖，更多的扩展和插件功能请查看[插件库](https://github.com/go-chassis/go-chassis-extension)

## Get started 
1. 生成 go mod
```bash
go mod init
```
2. 增加go chassis 
```shell script
GO111MODULE=on go get github.com/go-chassis/go-chassis
```
如果你面临网络问题
```bash
export GOPROXY=https://goproxy.io
```

3. [开始编写你的第一个微服务](https://go-chassis.readthedocs.io/en/latest/getstarted/writing-rest.html)


## 文档
这是在线文档 [here](https://go-chassis.readthedocs.io/), 
你可以跟随[这里](docs/README.md)来生成本地文档

## 例子
当前有两个例子库，提供很多使用场景
- [这里](examples)
- [这里](https://github.com/go-chassis/go-chassis-examples)

## 开发 go chassis

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

## Known Users

> [欢迎在此登录自己的信息](https://github.com/go-chassis/go-chassis/issues/592)

![huawei](https://raw.githubusercontent.com/go-chassis/go-chassis.github.io/master/known_users/huawei.PNG) 
![qutoutiao](https://raw.githubusercontent.com/go-chassis/go-chassis.github.io/master/known_users/qutoutiao.PNG)
![Shopee](https://raw.githubusercontent.com/go-chassis/go-chassis.github.io/master/known_users/Shopee.png)
![ieg](https://raw.githubusercontent.com/go-chassis/go-chassis.github.io/master/known_users/tencent-ieg.png)

