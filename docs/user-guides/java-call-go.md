# Java 和 Go 互相调用

当你在一个项目中需要同时使用java和go语言时，java语言对应的框架为servicecomb-java-chassis，
go对应的框架为go-chassis，他们之间可以通过servicecomb-service-center进行关联。

## Go as provider

请注意，java调用go时，反序列化需要用到schema中的x-java-interface标示所需要的class，但是，go-chassis并不能帮助你生成，需要你把这个参数添加到schema中，所以，请在chassis.yaml中设定 noRefreshSchema: true，表示不会自动生成schema，否则每次启动都会被覆盖。除此之外，并不需要任何的特殊设置。

代码请参考https://github.com/go-chassis/go-chassis-examples/java-call-go

## Java as provider

当Java作为提供者时，需要注意，必须使用HTTP通信，go-chassis不支持highway协议。
除此之外，无需任何特殊配置。