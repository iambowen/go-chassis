module github.com/go-chassis/go-chassis

require (
	github.com/DataDog/datadog-go v0.0.0-20180330214955-e67964b4021a
	github.com/apache/thrift v0.0.0-20180829120307-8de3749235db
	github.com/cactus/go-statsd-client v3.1.1+incompatible
	github.com/cenkalti/backoff v2.0.0+incompatible
	github.com/emicklei/go-restful v2.8.0+incompatible
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170208215640-dcef7f557305
	github.com/envoyproxy/go-control-plane v0.5.0
	github.com/go-chassis/go-archaius v0.0.0-20180831094429-e91e8f9ed49f
	github.com/go-chassis/go-cc-client v0.0.0-20180831085349-c2bb6cef1640
	github.com/go-chassis/go-chassis-plugins v0.0.0-20180731065901-7b05d8d2fbe6
	github.com/go-chassis/go-sc-client v0.0.0-20180925063328-78ad13b4fbef
	github.com/go-chassis/paas-lager v0.0.0-20180905100939-eff93e5e67db
	github.com/go-mesh/openlogging v0.0.0-20180912071658-0fd4707a75ab
	github.com/golang/protobuf v1.2.0
	github.com/hashicorp/go-version v1.0.0
	github.com/opentracing/opentracing-go v1.0.2
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.0.0-20180726151020-b85dc675b16b
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v0.8.0
	github.com/rcrowley/go-metrics v0.0.0-20180503174638-e2704e165165
	github.com/shirou/gopsutil v0.0.0-20180801053943-8048a2e9c577
	github.com/smartystreets/goconvey v0.0.0-20170602164621-9e8dc3f972df
	github.com/stretchr/testify v1.2.2
	go.uber.org/ratelimit v0.0.0-20180316092928-c15da0234277
	golang.org/x/net v0.0.0-20180824152047-4bcd98cce591
	google.golang.org/grpc v1.14.0
	gopkg.in/yaml.v2 v2.2.1
	k8s.io/apimachinery v0.0.0-20180823151430-017bf4f8f588
)

replace (
	github.com/kubernetes/client-go => ../k8s.io/client-go
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180824152047-4bcd98cce591 => github.com/golang/net v0.0.0-20180824152047-4bcd98cce591
	golang.org/x/sys v0.0.0-20180824143301-4910a1d54f87 => github.com/golang/sys v0.0.0-20180824143301-4910a1d54f87
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 => github.com/golang/time v0.0.0-20180412165947-fbb02b2291d2
	google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8 => github.com/google/go-genproto v0.0.0-20180817151627-c66870c02cf8
	google.golang.org/grpc v1.14.0 => github.com/grpc/grpc-go v1.14.0
)
