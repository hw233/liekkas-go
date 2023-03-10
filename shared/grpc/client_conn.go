package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	clientv3 "go.etcd.io/etcd/client/v3"

	"shared/grpc/balancer"
	"shared/grpc/module"
	"shared/utility/errors"
	"shared/utility/glog"

	"google.golang.org/grpc"
)

var (
	ErrServerNotMatch = errors.New("server not match")
)

type ClientConn struct {
	*grpc.ClientConn
	*module.Balancer
	*module.Discover
	*module.Recorder
	Service string

	// 用于从grpc的请求参数中获取id
	fetcher func(args interface{}) int64

	// 广播的函数
	broadcastMethods *sync.Map

	// 是否忽略 id!=0 && server == "" 的请求，减少无效请求
	ignoreNullServer bool

	// id指向服务器的缓存，减少对redis的请求次数，redis瓶颈可开启
	// 缓存有两种容错机制
	// 1. 设置缓存超时时间
	// 2. 使用ErrServerNotMatch错误刷新缓存数据重发
	serverCaches       *sync.Map
	serverCachesSwitch bool
	serverCachesExpire time.Duration
}

func NewClientConn(conn *grpc.ClientConn, service string, etcdClient *clientv3.Client, redisClient *redis.Client) *ClientConn {
	return &ClientConn{
		ClientConn: conn,
		Balancer:   module.NewBalancer(service, redisClient),
		Discover:   module.NewDiscover(service, etcdClient),
		Recorder:   module.NewRecorder(service, redisClient),
		Service:    service,
		fetcher:    func(args interface{}) int64 { return 0 },
		// fetcher: func(args interface{}) int64 {
		// 	if f, ok := args.(interface {
		// 		GetID(interface{}) int64
		// 	}); ok {
		// 		return f.GetID(args)
		// 	}
		//
		// 	return 0
		// },
		broadcastMethods:   &sync.Map{},
		ignoreNullServer:   false,
		serverCaches:       &sync.Map{},
		serverCachesSwitch: false,
	}
}

func (cc *ClientConn) RegisterFetcherFunc(fetcher func(interface{}) int64) {
	cc.fetcher = fetcher
}

// method: "Game/Test"
func (cc *ClientConn) RegisterBroadcastMethod(method string) {
	cc.broadcastMethods.Store(method, true)
}

// ignore if id != 0 && server == ""
func (cc *ClientConn) IgnoreNullServer() {
	cc.ignoreNullServer = true
}

func (cc *ClientConn) OpenServerCaches(expire time.Duration) {
	cc.ignoreNullServer = true
}

// balance strategy
// id = 0 round robin
// id != 0 least conn
func (cc *ClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	id := cc.fetcher(args)

	var servers []string
	var err error

	_, ok := cc.broadcastMethods.Load(method)
	if ok {
		// broadcast
		all, err := cc.ListServers(ctx)
		if err != nil {
			return err
		}

		servers = all
	} else {
		// singlecast
		if id == 0 {
			servers = append(servers, "")
		} else {
			server := ""

			if cc.serverCachesSwitch {
				// get server from cache
				cache, ok := cc.serverCaches.Load(id)
				if !ok {
					// get server from recorder
					server, err = cc.GetRecord(ctx, id)
					if err != nil {
						return err
					}

					if server != "" {
						// add cache
						cc.serverCaches.Store(id, newServerCache(server, cc.serverCachesExpire))
					}
				} else {
					serverCache := cache.(*serverCache)
					if serverCache.IsExpired() {
						// cache expired
						cc.serverCaches.Delete(id)

						// get server from recorder
						server, err = cc.GetRecord(ctx, id)
						if err != nil {
							return err
						}

						if server != "" {
							// add cache
							cc.serverCaches.Store(id, newServerCache(server, cc.serverCachesExpire))
						}
					}

					server = serverCache.server
				}
			} else {
				// get server from recorder
				server, err = cc.GetRecord(ctx, id)
				if err != nil {
					return err
				}
			}

			// ignore
			if cc.ignoreNullServer && server == "" {
				glog.Debugf("ignore rpc method[%s], id[%d], service[%s], server[%s], args[%+v]", method, id, cc.Service, server, args)
				return nil
			}

			servers = append(servers, server)
		}
	}

	for _, server := range servers {
		glog.Debugf("call rpc method[%s], id[%d], service[%s], server[%s], args[%+v]", method, id, cc.Service, server, args)
		err = cc.ClientConn.Invoke(balancer.WithCtxMetadata(ctx, server), method, args, reply, opts...)
		if err != nil {
			// cache expired
			if cc.serverCachesSwitch && err == ErrServerNotMatch {
				// refresh cache and retry
				cc.serverCaches.Delete(id)

				// get server from recorder
				server, err = cc.GetRecord(ctx, id)
				if err != nil {
					return err
				}

				if server != "" {
					// add cache
					cc.serverCaches.Store(id, newServerCache(server, cc.serverCachesExpire))
				}

				err = cc.ClientConn.Invoke(balancer.WithCtxMetadata(ctx, server), method, args, reply, opts...)
			} else {
				glog.Warnf("call rpc method[%s], id[%d], service[%s], server[%s], args[%+v], error: %v", method, id, cc.Service, server, args, err)
			}
		}
	}

	return err
}
