Features
================================
 - Pluggable registrator and discovery service: Support Service center,istio pilot and file based registry,fit both client side discovery and server side discovery pattern

 - Pluggable Protocol: You can custom your own protocol

 - Circuit breaker: Protect your service in runtime or on-demand

 - Route management: Able to route to different service based on weight and match rule to achieve Canary Release easily

 - Load balancing: Able to custom strategy and filter
 - Rate limiting: Both client side and server side rate limiting
 - Pluggable Cipher: Able to custom your own cipher for AKSK and TLS certs
 - Handler Chain: Able to add your own code during service calling for client and server side
 - Metrics: Able to expose Prometheus metric API automatically and custom metrics reporter
 - Tracing: Use opentracing-go as standard library, easy to integrate tracing impl
 - Logger: You can custom your own writer to sink log, by default support file and stdout
 - Hot-reconfiguraion: A lot of configuration can be reload in runtime, like loadbalancing, circuit breaker, rate limiting
 - Dynamic Configuration framework:   Able to develop a service which has hot-reconfiguration feature easily
 - Fault Injection: In consumer side, you can inject faults to bring chaos testing into your system