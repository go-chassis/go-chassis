module github.com/go-chassis/go-chassis

require (
	github.com/DataDog/zstd v1.3.5 // indirect
	github.com/apache/thrift v0.0.0-20180829120307-8de3749235db
	github.com/cenkalti/backoff v2.0.0+incompatible
	github.com/eapache/go-resiliency v1.1.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/emicklei/go-restful v2.8.0+incompatible
	github.com/go-chassis/go-archaius v0.0.0-20181108111652-ab19b4eae276
	github.com/go-chassis/go-cc-client v0.0.0-20181102101915-dea430061a34
	github.com/go-chassis/go-restful-swagger20 v0.0.0-20181221101811-a33c76fe4a6e
	github.com/go-chassis/go-sc-client v0.0.0-20190110124355-7e78d2170dca
	github.com/go-chassis/paas-lager v0.0.0-20181123014243-005283cca84c
	github.com/go-mesh/openlogging v0.0.0-20181122085847-3daf3ad8ed35
	github.com/golang/protobuf v1.2.0
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/hashicorp/go-version v1.0.0
	github.com/jtolds/gls v4.2.1+incompatible // indirect
	github.com/opentracing-contrib/go-observer v0.0.0-20170622124052-a52f23424492 // indirect
	github.com/opentracing/opentracing-go v1.0.2
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.0.0-20180726151020-b85dc675b16b
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v0.9.1
	github.com/prometheus/common v0.0.0-20190107103113-2998b132700a // indirect
	github.com/prometheus/procfs v0.0.0-20190104112138-b1a0a9a36d74 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20180503174638-e2704e165165
	github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d // indirect
	github.com/smartystreets/goconvey v0.0.0-20170602164621-9e8dc3f972df
	github.com/stretchr/testify v1.2.2
	go.uber.org/ratelimit v0.0.0-20180316092928-c15da0234277
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a
	google.golang.org/genproto v0.0.0-20181221175505-bd9b4fb69e2f // indirect
	google.golang.org/grpc v1.16.0
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
