Concepts
===================
Registry
注册中心负责微服务的注册和发现

Registrator
自注册组件，go-chassis在启动后会连接注册中心，自注册服务信息

Service Discovery
服务发现组件，负责服务发现并周期性轮询注册中心中的服务到本地缓存。

Protocol server and client
支持开发者自己将协议逻辑插入到go chassis中，接入统一的治理和微服务管理当中