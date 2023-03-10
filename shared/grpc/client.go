package grpc

import (
	"fmt"
	"sync"
	"time"

	"shared/grpc/balancer"
	"shared/grpc/balancer/picker/least_conn"
	"shared/grpc/resolver"
	"shared/utility/errors"

	"github.com/go-redis/redis/v8"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"github.com/openzipkin/zipkin-go/reporter"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	scheme = "overlord"

	// load balance
	// LeastConn  LoadBalance = "least_conn"
	// RoundRobin LoadBalance = "round_robin"
)

type LoadBalance string

type Client struct {
	etcdClient  *clientv3.Client
	redisClient *redis.Client
	clientConns *sync.Map
	// serviceCaches *sync.Map
}

func NewClient(etcdClient *clientv3.Client, redisClient *redis.Client) *Client {
	return &Client{
		etcdClient:  etcdClient,
		redisClient: redisClient,
		clientConns: &sync.Map{},
		// serviceCaches: &sync.Map{},
	}
}

func (c *Client) ClientConn(service string) (*ClientConn, bool) {
	v, ok := c.clientConns.Load(service)
	if !ok {
		return nil, false
	}

	return v.(*ClientConn), true
}

// // method: /Game/Connect
// // service: game
// func (c *Client) fetchService(method string) string {
// 	service, ok := c.serviceCaches.Load(method)
// 	if !ok {
// 		ss := strings.Split(method, "/")
// 		if len(ss) != 3 {
// 			return ""
// 		}
//
// 		service = strings.ToLower(ss[1])
//
// 		c.serviceCaches.Store(method, service)
// 	}
//
// 	return service.(string)
// }

// func (c *Client) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
// 	service := c.fetchService(method)
// 	if service == "" {
// 		return errors.New("method invalid")
// 	}
//
// 	clientConn, ok := c.ClientConn(service)
// 	if !ok {
// 		return errors.New("not found clientConn")
// 	}
//
// 	return clientConn.Invoke(ctx, method, args, reply, opts...)
// }
//
// func (c *Client) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
// 	service := c.fetchService(method)
// 	if service == "" {
// 		return nil, errors.New("method invalid")
// 	}
//
// 	clientConn, ok := c.ClientConn(service)
// 	if !ok {
// 		return nil, errors.New("not found clientConn")
// 	}
//
// 	return clientConn.NewStream(ctx, desc, method, opts...)
// }

func (c *Client) Dial(service string, serviceName string, opts ...grpc.DialOption) (*ClientConn, reporter.Reporter, error) {
	balancer.RegisterPickerBuilder(service, least_conn.NewLCPickerBuilder(c.redisClient, service))
	tracer, r, err := NewZipkinTracer(serviceName, "localhost:0")
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}
	cc, err := grpc.Dial(scheme+"://"+service, append(opts,
		grpc.WithStatsHandler(zipkingrpc.NewClientHandler(tracer)),
		grpc.WithResolvers(resolver.NewBuilder(c.etcdClient, scheme)),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{ "loadBalancingConfig": [{"%s": {}}] }`, service)),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             60 * time.Second,
			PermitWithoutStream: false,
		}))...)
	if err != nil {
		return nil, nil, err
	}

	clientConn := NewClientConn(cc, service, c.etcdClient, c.redisClient)

	c.clientConns.Store(service, clientConn)

	return clientConn, r, nil
}
