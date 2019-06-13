module github.com/go-chassis/go-chassis

require (
	github.com/cenkalti/backoff v2.0.0+incompatible
	github.com/emicklei/go-restful v2.8.0+incompatible
	github.com/go-chassis/foundation v0.0.0-20190516083152-b8b2476b6db7
	github.com/go-chassis/go-archaius v0.18.0
	github.com/go-chassis/go-chassis-config v0.7.0
	github.com/go-chassis/go-restful-swagger20 v1.0.1
	github.com/go-chassis/paas-lager v1.0.2-0.20190328010332-cf506050ddb2
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/go-mesh/openlogging v0.0.0-20181205082104-3d418c478b2d
	github.com/golang/protobuf v1.2.0
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/gorilla/websocket v1.4.0
	github.com/hashicorp/go-version v1.0.0
	github.com/opentracing/opentracing-go v1.0.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v0.9.1
	github.com/prometheus/client_model v0.0.0-20190115171406-56726106282f // indirect
	github.com/prometheus/common v0.2.0
	github.com/prometheus/procfs v0.0.0-20190117184657-bf6a532e95b1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20180503174638-e2704e165165
	github.com/smartystreets/assertions v0.0.0-20190116191733-b6c0e53d7304 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a
	github.com/stretchr/testify v1.2.2
	go.uber.org/ratelimit v0.0.0-20180316092928-c15da0234277
	gopkg.in/go-playground/assert.v1 v1.2.1
	gopkg.in/yaml.v2 v2.2.1
)

replace (
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180824152047-4bcd98cce591 => github.com/golang/net v0.0.0-20180824152047-4bcd98cce591
	golang.org/x/sys v0.0.0-20180824143301-4910a1d54f87 => github.com/golang/sys v0.0.0-20180824143301-4910a1d54f87
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 => github.com/golang/time v0.0.0-20180412165947-fbb02b2291d2
	google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8 => github.com/google/go-genproto v0.0.0-20180817151627-c66870c02cf8
)
