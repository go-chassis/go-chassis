Concepts
===================
Registry
 A registry must support both registration and discovery

Registrator
 A registrator service must support registration

Service Discovery
 A Service Discovery service must support discovery service at least.
 For example,ServiceComb service center support both registration and discovery
 Istio and kubernetes only support discovery

 .. image:: registry.PNG


Protocol server and client
 go chassis allows you to integrate any protocol into standardized model Invocation, so that any protocol can reuse same function like circuit breaker, load balancing, rate limiting, route management

 .. image:: protocol.PNG