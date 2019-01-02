module github.com/go-chassis/go-chassis

require (
	github.com/DataDog/zstd v1.3.4 // indirect
	github.com/Shopify/sarama v1.20.0 // indirect
	github.com/apache/thrift v0.0.0-20180829120307-8de3749235db
	github.com/cenkalti/backoff v2.0.0+incompatible
	github.com/emicklei/go-restful v2.8.0+incompatible
	github.com/go-chassis/go-archaius v0.0.0-20181108111652-ab19b4eae276
	github.com/go-chassis/go-cc-client v0.0.0-20181102101915-dea430061a34
	github.com/go-chassis/go-restful-swagger20 v0.0.0-20181221101811-a33c76fe4a6e
	github.com/go-chassis/go-sc-client v0.0.0-20181229093415-d2797ce547c9
	github.com/go-chassis/paas-lager v0.0.0-20181123014243-005283cca84c
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/go-mesh/openlogging v0.0.0-20181122085847-3daf3ad8ed35
	github.com/gogo/protobuf v1.2.0 // indirect
	github.com/golang/protobuf v1.2.0
	github.com/hashicorp/go-version v1.0.0
	github.com/opentracing/opentracing-go v1.0.2
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.0.0-20180726151020-b85dc675b16b
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/common v0.0.0-20181218105931-67670fe90761 // indirect
	github.com/prometheus/procfs v0.0.0-20181204211112-1dc9a6cbc91a // indirect
	github.com/rcrowley/go-metrics v0.0.0-20180503174638-e2704e165165
	github.com/smartystreets/goconvey v0.0.0-20170602164621-9e8dc3f972df
	github.com/stretchr/testify v1.2.2
	go.uber.org/ratelimit v0.0.0-20180316092928-c15da0234277
	google.golang.org/grpc v1.14.0
	gopkg.in/yaml.v2 v2.2.1
)

replace (
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180824152047-4bcd98cce591 => github.com/golang/net v0.0.0-20180824152047-4bcd98cce591
	golang.org/x/sys v0.0.0-20180824143301-4910a1d54f87 => github.com/golang/sys v0.0.0-20180824143301-4910a1d54f87
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 => github.com/golang/time v0.0.0-20180412165947-fbb02b2291d2
	google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8 => github.com/google/go-genproto v0.0.0-20180817151627-c66870c02cf8
	google.golang.org/grpc v1.14.0 => github.com/grpc/grpc-go v1.14.0
)
