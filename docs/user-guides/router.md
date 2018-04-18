# Router
## 概述

路由策略可应用于AB测试场景和新版本的灰度升级，主要通过路由规则来根据请求的来源、目标服务、Http Header及权重将服务访问请求分发到不同版本的微服务实例中。

## 配置

路由规则当前只支持在配置文件配置，支持rest和highway协议。

#### Consumer配置

灰度发布的路由规则只在服务的消费端配置使用，用于**将特定的请求，按一定权重，分发至同一服务名的不同分组。**用户可在conf/router.yaml 文件中设置：

```yaml
routeRule:  
  {targetServiceName}: # 服务名
    - precedence: {number} #优先级
      match:        #匹配策略
        source: {sourceServiceName} #匹配某个服务名
        headers:          #header匹配
          {key0}:            
            regex: {regex}
          {key1}         
            exact: {=？}   
      route: #路由规则
        - weight: {percent} #权重值
          tags:
            version: {version1}
            app: {appId}
    - precedence: {number1}
      match:        
        refer: {sourceTemplateName} #参考某个source模板ID
      route:
        - weight: {percent}
          tags:
            version: {version2}
            app: {appId}        
sourceTemplate:  #定义source模板
  {templateName}: # source 模板ID
    source: {sourceServiceName} 
    sourceTags：
      {tag}：{value}
    headers:
      {key0}:
        regex: {regex}
      {key1}
        exact: {=?}
      {key2}:
        noEqu: {!=?}
      {key3}
        greater: {>?}    
      {key4}:
        less: {<?}
      {key5}
        noLess: {>=?}      
      {key6}:
        noGreater: {<=?}
```

路由规则说明：

- 匹配特定请求由match配置，匹配条件是：source（源服务名）、source  tags 及headers，另外也可以使用refer字段来使用source模板进行匹配。
- Match中的Source Tags用于和服务调用请求中的sourceInfo中的tags 进行逐一匹配。
- Header中的字段的匹配支持正则, =, !=, &gt;, &lt;, &gt;=, &lt;=6种匹配方式。
- 如果未定义match，则可匹配任何请求。
- 转发权重定义在routeRule.{targetServiceName}.route下，由weight配置。
- 服务分组定义在routeRule.{targetServiceName}.route下，由tags配置，配置内容有version和app。

## API

##### 设置Router Rules

这个接口会彻底覆盖运行时的路由规则

```
router.SetRouteRule(rr map[string][]*config.RouteRule)
```

##### 获取Router Rules

```
router.GetRouteRule() 返回值 map[string][]*config.RouteRule
```

## 例子

#### 目标服务

每个路由规则的目标服务名称都由routeRule中的Key值指定。例如下表所示，所有以“Carts”服务为目标服务的路由规则均被包含在以“Carts”为Key值的列表中。

```yaml
routeRule:
  Carts:
    - precedence: 1
      route:
        - weight: 100 #percent          
          tags:            
            version: 0.0.1
```

Key值（目标服务名称）应该满足是一个合法的域名称。例如，一个在服务中心中注册的服务名称。

#### 规则优先级

针对某个特定的目标服务可以定义多条路由规则，在路由规则匹配过程中的匹配顺序按照各个规则的“precedence”字段的值来确定。“precedence”字段是可选配置，默认为0，值越大则优先级越高。如果两个规则的“precedence”配置值相同，则它们的实际匹配顺序是不确定的。

一种常用的模式是为指定的目标服务提供一条或多条具有高优先级的匹配请求源/Header的路由规则，并同时提供一条低优先级的只依照版本权重分发请求的路由规则，且这条规则不设置任何匹配条件以处理剩余的其他请求。

以下面的路由规则为例，对所有访问“Carts“服务的请求，如果满足header中包含”Foo：bar“，则将请求分发到服务的”2.0“版本的实例中，剩余的其他请求全部分发到”1.0“版本的实例中。

```yaml
routeRule:
  Carts:
    - precedence: 2
      match:
        headers:
          Foo:
            exact: bar
      route:
        - weight: 100           
          tags:            
            version: 2.0
    - precedence: 1
      route:
        - weight: 100   
          tags:            
            version: 1.0
```

#### 请求匹配规则

```yaml
match:
  refer: {templateName}
  source: {sourceServiceName}
  headers:
    {key}:
      regex: {regex}
    {key1}:
      exact: {exact}
```

请求的匹配规则属性配置如下：

| 属性              | 说明                                       |
| :-------------- | :--------------------------------------- |
| 引用（refer）       | 引用的匹配规则模板名称，用户可选择在sourceTemplate中定义匹配规则模板，并在此处引用。若引用了匹配规则模板，则其他配置项不用配置。 |
| 源服务（source）     | 表示发送请求的服务，和consumer是一个意义。                |
| header（headers） | 匹配headers。如果配置了多条Header 字段校验规则，则需要同时满足所有规则才可完成路由规则匹配。匹配方式有以下几种：精确匹配（exact）：header必须等于配置; 正则（regex）：按正则匹配header内容;  不等于（noEqu）：header不等于配置值;  大于等于（noLess）： header不小于配置值;  小于等于（noGreater）：header不大于配置值;   大于（greater）：header大于配置值;    小于（less）： header小于配置值 |

示例：

```yaml
match:
  source: vmall
  headers:
    cookie:
      regex: "^(.*?;)?(user=jason)(;.*)?$"
```

仅适用于来自vmall，header中的“cookie”字段包含“user=jason"的服务访问请求。

#### 分发规则

每个路由规则中都会定义一个或多个具有权重标识的后端服务，这些后端服务对应于用标签标识的不同版本的目标服务的实例。如果某个标签对应的注册服务实例有多个，则指向该标签版本的服务请求会按照用户配置的负责均衡策略进行分发，默认会采用round-robin策略。

分发规则的属性配置如下：

| 属性         | 说明                                       |
| :--------- | :--------------------------------------- |
| 权重（weight） | 本条规则的分发比重，配置为1-100的整数，表示百分比。             |
| 标签（tags）   | 用于区分相同服务名的不同服务分组，支持的key值为：版本（version），默认值为0.0.1。应用（app），默认值为default。 |

下面的例子表示75%的访问请求会被分流到具有“version：2.0”标签的服务实例中，其余25%的访问请求会被分发到1.0版本的实例中。

```yaml
route:
  - weight: 75
    tags:
      version: 2.0
  - weight: 25
    tags:
      version: 1.0
```

#### 定义匹配模板

我们可以通过预定义源模板（模板中的结构为一个Match结构），并在match部分引用该模板来进行路由规则的匹配。在下面的例子中，“vmall-with-special-header”是一个预定义的源模板的Key值，并在Carts的请求匹配规则中被引用。

```yaml
routeRule:
  Carts:
    - precedence: 2
      match:
        refer: vmall-with-special-header
      route:
        - weight: 100           
          tags:            
            version: 2.0
sourceTemplate:
  vmall-with-special-header:
    source: vmall
    headers:
      cookie:
        regex: "^(.*?;)?(user=jason)(;.*)?$"
```