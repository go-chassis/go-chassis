Discovery
======================

----

Introduction
++++

go-chassis is able to use Istio Pilot as discovery service.

Configuration
++++

edit chassis.yaml.

**registrator.disabled**

 Must disable registrator, because registrator is is used in client side discovery. go-chassis leverage server side discovery which supported by kubernetes

**serviceDiscovery.type**

 specify the plugin type to pilotv2

**serviceDiscovery.address**

 the pilot address


example
++++

::

  cse:
    service:
      Registry:
        registrator:
          disabled: true
        serviceDiscovery:
          type: pilotv2
          address: grpc://istio-pilot.istio-system:15010

