package resolver

import (
	"context"
	"sync"

	"shared/grpc/module"
	"shared/utility/glog"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/codes"
	gresolver "google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

type builder struct {
	scheme string
	c      *clientv3.Client
}

func (b builder) Build(target gresolver.Target, cc gresolver.ClientConn, opts gresolver.BuildOptions) (gresolver.Resolver, error) {
	r := &resolver{
		c:  b.c,
		cc: cc,
	}

	r.ctx, r.cancel = context.WithCancel(context.Background())

	watcher, err := newWatcher(r.ctx, b.c, target.URL.Host)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resolver: failed to new watcher: %s", err)
	}

	r.watcher = watcher
	r.wg.Add(1)

	go r.watch()

	return r, nil
}

func (b builder) Scheme() string {
	return b.scheme
}

// NewBuilder creates a resolver builder.
func NewBuilder(client *clientv3.Client, scheme string) gresolver.Builder {
	return builder{c: client, scheme: scheme}
}

type resolver struct {
	watcher *watcher
	c       *clientv3.Client
	cc      gresolver.ClientConn
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func (r *resolver) watch() {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		case resolverMessages, ok := <-r.watcher.receiver:
			if !ok {
				return
			}

			addrs := make([]gresolver.Address, 0, len(resolverMessages))

			for _, resolverMessage := range resolverMessages {
				balancerMessage := module.NewBalancerMessage(resolverMessage.Event)

				addr := gresolver.Address{
					Addr:       resolverMessage.Addr,
					ServerName: resolverMessage.Server,
					Attributes: attributes.New(module.AttrBalancerMessage, balancerMessage),
				}

				addrs = append(addrs, addr)
			}

			// 没有任何操作的服务不会更新
			err := r.cc.UpdateState(gresolver.State{Addresses: addrs})
			if err != nil {
				glog.Errorf("resolver: cc.UpdateState error: %v", err)
			}
		}
	}
}

// ResolveNow is a no-op here.
// It's just a hint, resolver can ignore this if it's not necessary.
func (r *resolver) ResolveNow(gresolver.ResolveNowOptions) {}

func (r *resolver) Close() {
	r.cancel()
	r.wg.Wait()
}
