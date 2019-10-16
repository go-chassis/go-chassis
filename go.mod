module github.com/go-chassis/go-chassis

require (
	github.com/cenkalti/backoff v2.0.0+incompatible
	github.com/emicklei/go-restful v2.8.0+incompatible
	github.com/go-chassis/foundation v0.0.0-20190621030543-c3b63f787f4c
	github.com/go-chassis/go-archaius v0.23.0
	github.com/go-chassis/go-chassis-config v0.12.1-0.20190924031444-785da26fa7df
	github.com/go-chassis/go-restful-swagger20 v1.0.1
	github.com/go-chassis/paas-lager v1.0.2-0.20190328010332-cf506050ddb2
	github.com/go-mesh/openlogging v1.0.1
	github.com/golang/protobuf v1.2.0
	github.com/gorilla/websocket v1.4.0
	github.com/hashicorp/go-version v1.0.0
	github.com/opentracing/opentracing-go v1.0.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.1
	github.com/prometheus/common v0.2.0
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a
	github.com/stretchr/testify v1.3.0
	go.uber.org/ratelimit v0.1.0
	gopkg.in/yaml.v2 v2.2.1
)

replace (
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180824152047-4bcd98cce591 => github.com/golang/net v0.0.0-20180824152047-4bcd98cce591
	golang.org/x/sys v0.0.0-20180824143301-4910a1d54f87 => github.com/golang/sys v0.0.0-20180824143301-4910a1d54f87
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 => github.com/golang/time v0.0.0-20180412165947-fbb02b2291d2
	google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8 => github.com/google/go-genproto v0.0.0-20180817151627-c66870c02cf8
	github.com/apache/servicecomb-kie v0.1.0 => github.com/apache/servicecomb-kie v0.0.0-20190905142319-5ee098c8886f // indirect. TODO: remove this line when servicecomb-kie has a stable release
)
