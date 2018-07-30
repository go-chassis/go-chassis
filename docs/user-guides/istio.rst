与Istio集成
======================

服务发现
----

概述
++++

go-chassis对接Istio的pilot组件，实现服务发现能力。

配置
++++

接入pilot进行服务发现的配置在chassis.yaml。

**registrator.disabled**

 设置禁用服务主动注册能力，与Istio集成后，由其对接的平台进行服务注册

**serviceDiscovery.type**

 设置启用接入pilot插件

**serviceDiscovery.address**

 pilot服务地址 允许配置多个以逗号隔开


示例
++++

::

  cse:
    service:
      Registry:
        registrator:
          disabled: true
        serviceDiscovery:
          type: pilot
          address: http://istio-pilot.istio-system:8080

