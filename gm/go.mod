module gm

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis/v8 v8.11.4
	github.com/openzipkin/zipkin-go v0.4.0
	go.etcd.io/etcd/client/v3 v3.5.2
	shared v0.0.0
)

replace shared => ../shared
