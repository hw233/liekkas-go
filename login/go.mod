module login

go 1.16

replace shared => ../shared

require (
	github.com/go-redis/redis/v8 v8.11.4
	go.etcd.io/etcd/client/v3 v3.5.2
	google.golang.org/grpc v1.44.0
	shared v0.0.0-00010101000000-000000000000
)
