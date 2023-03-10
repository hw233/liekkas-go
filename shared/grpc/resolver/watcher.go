package resolver

import (
	"context"

	"shared/grpc/module"
	"shared/utility/glog"
	"shared/utility/key"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// watch status of server
type watcher struct {
	client   *clientv3.Client
	receiver chan []*module.ResolverMessage
}

func newWatcher(ctx context.Context, client *clientv3.Client, service string) (*watcher, error) {
	// init server
	resp, err := client.Get(ctx, service, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}

	messages := make([]*module.ResolverMessage, 0, len(resp.Kvs))

	for _, v := range resp.Kvs {
		messages = append(
			messages,
			module.NewResolverMessage(key.SubEtcdKey(string(v.Key)), string(v.Value)),
		)
	}

	c := make(chan []*module.ResolverMessage, 1)
	if len(messages) > 0 {
		c <- messages
	}

	watcher := &watcher{
		client:   client,
		receiver: c,
	}

	go watcher.watch(ctx, service, resp.Header.Revision+1)

	return watcher, nil
}

func (w *watcher) watch(ctx context.Context, service string, rev int64) {
	defer close(w.receiver)

	watch := w.client.Watch(ctx, service, clientv3.WithRev(rev), clientv3.WithPrefix())

	glog.Infof("start watching service [%s]", service)

	for {
		select {
		case <-ctx.Done():
			return
		case resp, ok := <-watch:
			if !ok {
				return
			}
			if resp.Err() != nil {
				return
			}

			messages := make([]*module.ResolverMessage, 0, len(resp.Events))
			for _, v := range resp.Events {
				server := ""
				ret := key.SplitEtcdKey(string(v.Kv.Key))
				if len(ret) >= 2 {
					server = ret[1]
				}

				message := module.NewResolverMessage(server, string(v.Kv.Value))

				switch v.Type {
				case clientv3.EventTypePut:
					message.Event = module.Register
				case clientv3.EventTypeDelete:
					message.Event = module.Unregister
				default:
					continue
				}

				messages = append(messages, message)
			}

			if len(messages) > 0 {
				w.receiver <- messages
			}
		}
	}
}
