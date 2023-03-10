module portal

go 1.16

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/openzipkin/zipkin-go v0.4.0
	github.com/panjf2000/ants/v2 v2.4.6
	github.com/panjf2000/gnet v1.5.1
	go.etcd.io/etcd/client/v3 v3.5.2
	go.uber.org/zap v1.21.0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
	shared v0.0.0
)

replace shared => ../shared
