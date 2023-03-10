module gamesvr

go 1.16

require (
	github.com/antlabs/timer v0.0.5
	github.com/go-redis/redis/v8 v8.11.4
	github.com/openzipkin/zipkin-go v0.4.0
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9
	go.etcd.io/etcd/client/v3 v3.5.2
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
	shared v0.0.0
)

replace shared => ../shared
